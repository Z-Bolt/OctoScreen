package octoprintApis

import (
	"bytes"
	"encoding/json"
	// "io"

	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// RestartRequest restart the print of the currently selected file from the
// beginning. There must be an active print job for this to work and the print
// job must currently be paused
type RestartRequest struct{}

// Do sends an API request and returns an error if any.
func (cmd *RestartRequest) Do(c *Client) error {
	payload := map[string]string{"command": "restart"}

	buffer := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buffer).Encode(payload); err != nil {
		return err
	}

	_, err := c.doJsonRequest("POST", JobApiUri, buffer, JobToolErrors, true)
	return err
}
