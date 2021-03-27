package octoprint

import (
	"encoding/json"
	"fmt"
)

var ExecuteErrors = StatusMapping{
	404: "The command could not be found for source and action",
	500: "The command didn’t define a command to execute, the command returned a non-zero return code and ignore was not true or some other internal server error occurred",
}

const URISystemCommands = "/api/system/commands"

// SystemCommandsRequest retrieves all configured system commands.
type SystemCommandsRequest struct{}

// Do sends an API request and returns the API response.
func (cmd *SystemCommandsRequest) Do(c *Client) (*SystemCommandsResponse, error) {
	b, err := c.doJSONRequest("GET", URISystemCommands, nil, nil)
	if err != nil {
		return nil, err
	}

	r := &SystemCommandsResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}
	for i := range r.Core {
		x := r.Core[i]
		if err2 := json.Unmarshal(x.RawConfirm, x.Confirm); err2 != nil {
		    x.Confirm = ""
		}
	}
	for i := range r.Custom {
		x := r.Custom[i]
		if err2 := json.Unmarshal(x.RawConfirm, x.Confirm); err2 != nil {
		    x.Confirm = ""
		}
	}

	return r, err
}

// SystemExecuteCommandRequest retrieves all configured system commands.
type SystemExecuteCommandRequest struct {
	// Source for which to list commands.
	Source CommandSource `json:"source"`

	// Action is the identifier of the command, action from its definition.
	Action string `json:"action"`
}

// Do sends an API request and returns an error if any.
func (cmd *SystemExecuteCommandRequest) Do(c *Client) error {
	uri := fmt.Sprintf("%s/%s/%s", URISystemCommands, cmd.Source, cmd.Action)
	_, err := c.doJSONRequest("POST", uri, nil, ExecuteErrors)
	return err
}
