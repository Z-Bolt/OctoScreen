package octoprintApis

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// PauseRequest pauses/resumes/toggles the current print job.
type PauseRequest struct {
	// Action specifies which action to take.
	// In order to stay backwards compatible to earlier iterations of this API,
	// the default action to take if no action parameter is supplied is to
	// toggle the print job status.
	Action dataModels.PauseAction `json:"action"`
}

// Do sends an API request and returns an error if any.
func (cmd *PauseRequest) Do(c *Client) error {
	buffer := bytes.NewBuffer(nil)
	if err := cmd.encode(buffer); err != nil {
		return err
	}

	_, err := c.doJsonRequest("POST", JobApiUri, buffer, JobToolErrors, true)
	return err
}

func (cmd *PauseRequest) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		PauseRequest
	}{
		Command:      "pause",
		PauseRequest: *cmd,
	})
}
