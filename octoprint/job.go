package octoprint

import (
	"bytes"
	"encoding/json"
	"io"
)

const JobTool = "/api/job"

// JobCommand retrieve information about the current job (if there is one).
type JobCommand struct{}

// Do sends an API request and returns the API response.
func (cmd *JobCommand) Do(c *Client) (*JobResponse, error) {
	b, err := c.doRequest("GET", JobTool, nil)
	if err != nil {
		return nil, err
	}

	r := &JobResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}

// StartCommand starts the print of the currently selected file.
type StartCommand struct{}

// Do sends an API request and returns an error if any.
func (cmd *TargetCommand) Do(c *Client) error {
	payload := map[string]string{"command": "start"}

	b := bytes.NewBuffer(nil)
	if err := json.NewEncoder(b).Encode(payload); err != nil {
		return err
	}

	_, err := c.doRequest("POST", URITool, b)
	return err
}

// CancelCommand cancels the current print job.
type CancelCommand struct{}

// Do sends an API request and returns an error if any.
func (cmd *CancelCommand) Do(c *Client) error {
	payload := map[string]string{"command": "cancel"}

	b := bytes.NewBuffer(nil)
	if err := json.NewEncoder(b).Encode(payload); err != nil {
		return err
	}

	_, err := c.doRequest("POST", URITool, b)
	return err
}

// RestartCommand restart the print of the currently selected file from the
// beginning. There must be an active print job for this to work and the print
// job must currently be paused
type RestartCommand struct{}

// Do sends an API request and returns an error if any.
func (cmd *RestartCommand) Do(c *Client) error {
	payload := map[string]string{"command": "restart"}

	b := bytes.NewBuffer(nil)
	if err := json.NewEncoder(b).Encode(payload); err != nil {
		return err
	}

	_, err := c.doRequest("POST", URITool, b)
	return err
}

type PauseAction string

const (
	// Pause the current job if it’s printing, does nothing if it’s already paused.
	Pause PauseAction = "pause"
	// Resume the current job if it’s paused, does nothing if it’s printing.
	Resume PauseAction = "resume"
	// Toggle the pause state of the job, pausing it if it’s printing and
	// resuming it if it’s currently paused.
	Toggle PauseAction = "toggle"
)

// PauseCommand pauses/resumes/toggles the current print job.
type PauseCommand struct {
	// Action specifies which action to take.
	// In order to stay backwards compatible to earlier iterations of this API,
	// the default action to take if no action parameter is supplied is to
	// toggle the print job status.
	Action PauseAction `json:"action"`
}

// Do sends an API request and returns an error if any.
func (cmd *PauseCommand) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doRequest("POST", URITool, b)
	return err
}

func (cmd *PauseCommand) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		PauseCommand
	}{
		Command:      "pause",
		PauseCommand: *cmd,
	})
}
