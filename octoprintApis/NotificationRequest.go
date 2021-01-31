package octoprintApis

import (
	// "bytes"
	"encoding/json"
	"fmt"
	// "log"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


type NotificationRequest struct {
	Command string `json:"command"`
}

func (cmd *NotificationRequest) Do(c *Client) (*dataModels.NotificationResponse, error) {
	target := fmt.Sprintf("%s?command=get_notification", PluginZBoltOctoScreenApiUri)
	bytes, err := c.doJsonRequest("GET", target, nil, ConnectionErrors)
	if err != nil {
		return nil, err
	}

	response := &dataModels.NotificationResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, err
}
