package octoprintApis

import (
	"encoding/json"
	// "fmt"

	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


const SystemCommandsApiUri = "/api/system/commands"

// SystemCommandsRequest retrieves all configured system commands.
type SystemCommandsRequest struct{}

// Do sends an API request and returns the API response.
func (cmd *SystemCommandsRequest) Do(c *Client) (*dataModels.SystemCommandsResponse, error) {
	bytes, err := c.doJsonRequest("GET", SystemCommandsApiUri, nil, nil, true)
	if err != nil {
		return nil, err
	}

	response := &dataModels.SystemCommandsResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	for i := range response.Core {
		commandDefinition := response.Core[i]
		convertRawConfirm(commandDefinition)
	}

	for i := range response.Custom {
		commandDefinition := response.Custom[i]
		convertRawConfirm(commandDefinition)
	}

	return response, err
}

func convertRawConfirm(commandDefinition *dataModels.CommandDefinition) {
	if commandDefinition == nil || commandDefinition.RawConfirm == nil || len(commandDefinition.RawConfirm) < 1 {
		return
	}

	err := json.Unmarshal(commandDefinition.RawConfirm, &commandDefinition.Confirm)
	if err != nil {
		logger.LogError("SystemCommandsRequest.convertRawConfirm()", "json.Unmarshal(Custom)", err)
		commandDefinition.Confirm = ""
	}
}
