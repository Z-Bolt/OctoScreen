package octoprintApis

import (
	"encoding/json"
	"fmt"

	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


const SystemCommandsApiUri = "/api/system/commands"


var ExecuteErrors = StatusMapping {
	404: "The command could not be found for source and action",
	500: "The command didn’t define a command to execute, the command returned a non-zero return code and ignore was not true or some other internal server error occurred",
}


// SystemCommandsRequest retrieves all configured system commands.
type SystemCommandsRequest struct{}

// Do sends an API request and returns the API response.
func (cmd *SystemCommandsRequest) Do(c *Client) (*dataModels.SystemCommandsResponse, error) {
	bytes, err := c.doJsonRequest("GET", SystemCommandsApiUri, nil, nil)
	if err != nil {
		return nil, err
	}

	response := &dataModels.SystemCommandsResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	for i := range response.Core {
		commandDefinition := response.Core[i]
		err = json.Unmarshal(commandDefinition.RawConfirm, &commandDefinition.Confirm)
		if err != nil {
			logger.LogError("SystemCommandsRequest.Do()", "json.Unmarshal(Core)", err)
			commandDefinition.Confirm = ""
			return nil, err
		}
	}

	for i := range response.Custom {
		commandDefinition := response.Custom[i]
		err = json.Unmarshal(commandDefinition.RawConfirm, &commandDefinition.Confirm)
		if err != nil {
			logger.LogError("SystemCommandsRequest.Do()", "json.Unmarshal(Custom)", err)
			commandDefinition.Confirm = ""
			return nil, err
		}
	}

	return response, err
}

// SystemExecuteCommandRequest retrieves all configured system commands.
type SystemExecuteCommandRequest struct {
	// Source for which to list commands.
	Source dataModels.CommandSource `json:"source"`

	// Action is the identifier of the command, action from its definition.
	Action string `json:"action"`
}

// Do sends an API request and returns an error if any.
func (cmd *SystemExecuteCommandRequest) Do(c *Client) error {
	uri := fmt.Sprintf("%s/%s/%s", SystemCommandsApiUri, cmd.Source, cmd.Action)
	_, err := c.doJsonRequest("POST", uri, nil, ExecuteErrors)
	return err
}
