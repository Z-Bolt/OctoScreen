package octoprintApis

import (
	// "encoding/json"
	"fmt"

	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


var ExecuteErrors = StatusMapping {
	404: "The command could not be found for source and action",
	500: "The command didn’t define a command to execute, the command returned a non-zero return code and ignore was not true or some other internal server error occurred",
}

// SystemExecuteCommandRequest retrieves all configured system commands.
type SystemExecuteCommandRequest struct {
	// Source for which to list commands.
	Source dataModels.CommandSource `json:"source"`

	// Action is the identifier of the command, action from its definition.
	Action string `json:"action"`
}

// Do sends an API request and returns an error if any.
func (this *SystemExecuteCommandRequest) Do(client *Client) error {
	uri := fmt.Sprintf("%s/%s/%s", SystemCommandsApiUri, this.Source, this.Action)
	_, err := client.doJsonRequest("POST", uri, nil, ExecuteErrors, true)
	if err != nil {
		logger.LogError("SystemExecuteCommandRequest.Do()", "client.doJsonRequest(POST)", err)
	}

	return err
}
