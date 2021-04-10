package octoprintApis

import (
	"bytes"
	"encoding/json"
	// "io"

	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// StartRequest starts the print of the currently selected file.
type StartRequest struct{}

// Do sends an API request and returns an error if any.
func (cmd *StartRequest) Do(c *Client) error {
	payload := map[string]string{"command": "start"}

	buffer := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buffer).Encode(payload); err != nil {
		return err
	}

	_, err := c.doJsonRequest("POST", JobApiUri, buffer, JobToolErrors, true)
	return err
}
