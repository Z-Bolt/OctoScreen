package octoprintApis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// SelectFileRequest selects a file for printing.
type SelectFileRequest struct {
	// Location is target location on which to send the command for is located,
	// either local (for OctoPrint’s uploads folder) or sdcard for the
	// printer’s SD card (if available).
	Location dataModels.Location `json:"-"`

	// Path  of the file for which to issue the command.
	Path string `json:"-"`

	// Print, if set to true the file will start printing directly after
	// selection. If the printer is not operational when this parameter is
	// present and set to true, the request will fail with a response of
	// 409 Conflict.
	Print bool `json:"print"`
}

// Do sends an API request and returns an error if any.
func (cmd *SelectFileRequest) Do(c *Client) error {
	buffer := bytes.NewBuffer(nil)
	if err := cmd.encode(buffer); err != nil {
		return err
	}

	uri := fmt.Sprintf("%s/%s/%s", FilesApiUri, cmd.Location, cmd.Path)
	_, err := c.doJsonRequest("POST", uri, buffer, FilesLocationPathPOSTErrors, true)
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
