package octoprintApis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// UploadFileRequest uploads a file to the selected location or create a new
// empty folder on it.
type UploadFileRequest struct {
	// Location is the target location to which to upload the file. Currently
	// only `local` and `sdcard` are supported here, with local referring to
	// OctoPrint’s `uploads` folder and `sdcard` referring to the printer’s
	// SD card. If an upload targets the SD card, it will also be stored
	// locally first.
	Location dataModels.Location

	// Select whether to select the file directly after upload (true) or not
	// (false). Optional, defaults to false. Ignored when creating a folder.
	Select bool

	//Print whether to start printing the file directly after upload (true) or
	// not (false). If set, select is implicitely true as well. Optional,
	// defaults to false. Ignored when creating a folder.
	Print bool
	b     *bytes.Buffer
	w     *multipart.Writer
}

// AddFile adds a new file to be uploaded from a given reader.
func (req *UploadFileRequest) AddFile(filename string, r io.Reader) error {
	w, err := req.writer().CreateFormFile("file", filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	return err

}

func (req *UploadFileRequest) writer() *multipart.Writer {
	if req.w == nil {
		req.b = bytes.NewBuffer(nil)
		req.w = multipart.NewWriter(req.b)
	}

	return req.w
}

// AddFolder adds a new folder to be created.
func (req *UploadFileRequest) AddFolder(folder string) error {
	return req.writer().WriteField("foldername", folder)
}

// Do sends an API request and returns the API response.
func (req *UploadFileRequest) Do(c *Client) (*dataModels.UploadFileResponse, error) {
	req.addSelectPrintAndClose()

	uri := fmt.Sprintf("%s/%s", FilesApiUri, req.Location)
	bytes, err := c.doRequest("POST", uri, req.w.FormDataContentType(), req.b, FilesLocationPOSTErrors, true)
	if err != nil {
		return nil, err
	}

	response := &dataModels.UploadFileResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, err
}

func (req *UploadFileRequest) addSelectPrintAndClose() error {
	err := req.writer().WriteField("select", fmt.Sprintf("%t", req.Select))
	if err != nil {
		return err
	}

	err = req.writer().WriteField("print", fmt.Sprintf("%t", req.Print))
	if err != nil {
		return err
	}

	return req.writer().Close()
}
