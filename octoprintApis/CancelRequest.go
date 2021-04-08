package octoprintApis

import (
	"bytes"
	"encoding/json"
	// "io"

	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// CancelRequest cancels the current print job.
type CancelRequest struct{}

// Do sends an API request and returns an error if any.
func (cmd *CancelRequest) Do(c *Client) error {
	payload := map[string]string{"command": "cancel"}

	buffer := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buffer).Encode(payload); err != nil {
		return err
	}

	_, err := c.doJsonRequest("POST", JobApiUri, buffer, JobToolErrors, true)
	return err
}
