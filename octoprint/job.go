package octoprint

import (
	"encoding/json"
)

const JobTool = "/api/job"

// StateCommand retrieves the current state of the printer.
type JobCommand struct{}

// Do sends an API request and returns the API response.
func (cmd *JobCommand) Do(c *Client) (*StateResponse, error) {
	b, err := c.doRequest("GET", JobTool, nil)
	if err != nil {
		return nil, err
	}

	r := &StateResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}
