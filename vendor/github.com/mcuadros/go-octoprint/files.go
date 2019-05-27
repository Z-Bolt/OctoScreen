package octoprint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
)

type Location string

const (
	URIFiles = "/api/files"

	Local  Location = "local"
	SDCard Location = "sdcard"
)

var (
	FilesLocationGETErrors = statusMapping{
		404: "Location is neither local nor sdcard",
	}
	FilesLocationPOSTErrors = statusMapping{
		400: "No file or foldername are included in the request, userdata was provided but could not be parsed as JSON or the request is otherwise invalid.",
		404: "Location is neither local nor sdcard or trying to upload to SD card and SD card support is disabled",
		409: "The upload of the file would override the file that is currently being printed or if an upload to SD card was requested and the printer is either not operational or currently busy with a print job.",
		415: "The file is neither a gcode nor an stl file (or it is an stl file but slicing support is disabled)",
		500: "The upload failed internally",
	}
	FilesLocationPathPOSTErrors = statusMapping{
		400: "The command is unknown or the request is otherwise invalid",
		415: "A slice command was issued against something other than an STL file.",
		404: "Location is neither local nor sdcard or the requested file was not found",
		409: "Selected file is supposed to start printing directly but the printer is not operational or if a file to be sliced is supposed to be selected or start printing directly but the printer is not operational or already printing.",
	}
	FilesLocationDeleteErrors = statusMapping{
		404: "Location is neither local nor sdcard",
		409: "The file to be deleted is currently being printed",
	}
)

// FileRequest retrieves the selected file’s or folder’s information.
type FileRequest struct {
	// Location of the file for which to retrieve the information, either
	// `local` or `sdcard`.
	Location Location
	// Filename of the file for which to retrieve the information
	Filename string
	// Recursive if set to true, return all files and folders recursively.
	// Otherwise only return items on same level.
	Recursive bool
}

// Do sends an API request and returns the API response
func (cmd *FileRequest) Do(c *Client) (*FileInformation, error) {
	uri := fmt.Sprintf("%s/%s/%s?recursive=%t", URIFiles,
		cmd.Location, cmd.Filename, cmd.Recursive,
	)

	b, err := c.doJSONRequest("GET", uri, nil, FilesLocationGETErrors)
	if err != nil {
		return nil, err
	}

	r := &FileInformation{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}

// FilesRequest retrieve information regarding all files currently available and
// regarding the disk space still available locally in the system.
type FilesRequest struct {
	// Location is the target location .
	Location Location
	// Recursive if set to true, return all files and folders recursively.
	// Otherwise only return items on same level.
	Recursive bool
}

// Do sends an API request and returns the API response.
func (cmd *FilesRequest) Do(c *Client) (*FilesResponse, error) {
	uri := fmt.Sprintf("%s?recursive=%t", URIFiles, cmd.Recursive)
	if cmd.Location != "" {
		uri = fmt.Sprintf("%s/%s?recursive=%t", URIFiles, cmd.Location, cmd.Recursive)
	}

	b, err := c.doJSONRequest("GET", uri, nil, FilesLocationGETErrors)
	if err != nil {
		return nil, err
	}

	r := &FilesResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	if len(r.Children) > 0 {
		r.Files = r.Children
	}

	return r, err
}

// UploadFileRequest uploads a file to the selected location or create a new
// empty folder on it.
type UploadFileRequest struct {
	// Location is the target location to which to upload the file. Currently
	// only `local` and `sdcard` are supported here, with local referring to
	// OctoPrint’s `uploads` folder and `sdcard` referring to the printer’s
	// SD card. If an upload targets the SD card, it will also be stored
	// locally first.
	Location Location
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
func (req *UploadFileRequest) Do(c *Client) (*UploadFileResponse, error) {
	req.addSelectPrintAndClose()

	uri := fmt.Sprintf("%s/%s", URIFiles, req.Location)
	b, err := c.doRequest("POST", uri, req.w.FormDataContentType(), req.b, FilesLocationPOSTErrors)
	if err != nil {
		return nil, err
	}

	r := &UploadFileResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
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

// DeleteFileRequest delete the selected path on the selected location.
type DeleteFileRequest struct {
	// Location is the target location on which to delete the file, either
	// `local` (for OctoPrint’s uploads folder) or \sdcard\ for the printer’s
	// SD card (if available)
	Location Location
	// Path of the file to delete
	Path string
}

// Do sends an API request and returns error if any.
func (req *DeleteFileRequest) Do(c *Client) error {
	uri := fmt.Sprintf("%s/%s/%s", URIFiles, req.Location, req.Path)
	if _, err := c.doJSONRequest("DELETE", uri, nil, FilesLocationDeleteErrors); err != nil {
		return err
	}

	return nil
}

// SelectFileRequest selects a file for printing.
type SelectFileRequest struct {
	// Location is target location on which to send the command for is located,
	// either local (for OctoPrint’s uploads folder) or sdcard for the
	// printer’s SD card (if available)
	Location Location `json:"-"`
	// Path  of the file for which to issue the command
	Path string `json:"-"`
	// Print, if set to true the file will start printing directly after
	// selection. If the printer is not operational when this parameter is
	// present and set to true, the request will fail with a response of
	// 409 Conflict.
	Print bool `json:"print"`
}

// Do sends an API request and returns an error if any.
func (cmd *SelectFileRequest) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	uri := fmt.Sprintf("%s/%s/%s", URIFiles, cmd.Location, cmd.Path)
	_, err := c.doJSONRequest("POST", uri, b, FilesLocationPathPOSTErrors)
	return err
}

func (cmd *SelectFileRequest) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		SelectFileRequest
	}{
		Command:           "select",
		SelectFileRequest: *cmd,
	})
}
