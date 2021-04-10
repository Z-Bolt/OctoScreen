package octoprintApis

import (
	"encoding/json"
	"fmt"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// TODO: OctoScreenSettingsRequest seems like it's practically the same as PluginManagerInfoRequest
// Need to clean up and consolidate, or add comments as to why the two different classes.

type OctoScreenSettingsRequest struct {
	Command string `json:"command"`
}

func (this *OctoScreenSettingsRequest) Do(client *Client, uiState string) (*dataModels.OctoScreenSettingsResponse, error) {
	target := fmt.Sprintf("%s?command=get_settings", PluginZBoltOctoScreenApiUri)
	bytes, err := client.doJsonRequest("GET", target, nil, ConnectionErrors, false)
	if err != nil {
		return nil, err
	}

	response := &dataModels.OctoScreenSettingsResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, err
}
