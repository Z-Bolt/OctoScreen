package octoprintApis

import (
	"encoding/json"
	"fmt"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


type OctoScreenSettingsRequest struct {
	Command string `json:"command"`
}

func (cmd *OctoScreenSettingsRequest) Do(c *Client) (*dataModels.OctoScreenSettingsResponse, error) {
	target := fmt.Sprintf("%s?command=get_settings", PluginZBoltOctoScreenApiUri)
	bytes, err := c.doJsonRequest("GET", target, nil, ConnectionErrors)
	if err != nil {
		return nil, err
	}

	response := &dataModels.OctoScreenSettingsResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, err
}
