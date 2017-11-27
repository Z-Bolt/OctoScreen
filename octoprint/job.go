package octoprint

import (
	"encoding/json"
)

const JobTool = "/api/job"

// StateCommand retrieves the current state of the printer.
type JobCommand struct {
}

func (cmd *JobCommand) Do(p *Printer) (*StateResponse, error) {
	b, err := p.doRequest("GET", JobTool, nil)
	if err != nil {
		return nil, err
	}

	r := &StateResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}
