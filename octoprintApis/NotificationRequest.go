package octoprintApis

import (
	// "bytes"
	"encoding/json"
	"fmt"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


type NotificationRequest struct {
	Command string `json:"command"`
}

func (this *NotificationRequest) Do(client *Client, uiState string) (*dataModels.NotificationResponse, error) {
	if uiState != "splash" && uiState != "idle" {
		return nil, nil
	}

	target := fmt.Sprintf("%s?command=get_notification", PluginZBoltOctoScreenApiUri)
	bytes, err := client.doJsonRequest("GET", target, nil, ConnectionErrors, true)
	if err != nil {
		return nil, err
	}

	response := &dataModels.NotificationResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, err
}
