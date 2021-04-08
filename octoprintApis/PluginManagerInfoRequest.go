package octoprintApis

import (
	"bytes"
	"encoding/json"

	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


const pluginManagerRequestURI = "/api/plugin/pluginmanager"

// PluginManagerInfoRequest -
type PluginManagerInfoRequest struct {
	Command string `json:"command"`
}

// Do -
func (this *PluginManagerInfoRequest) Do(client *Client, uiState string) (*dataModels.PluginManagerInfoResponse, error) {
	this.Command = "get_settings"

	params := bytes.NewBuffer(nil)
	if err := json.NewEncoder(params).Encode(this); err != nil {
		logger.LogError("PluginManagerInfoRequest.Do()", "json.NewEncoder(params).Encode(this)", err)
		return nil, err
	}

	bytes, err := client.doJsonRequest("GET", pluginManagerRequestURI, params, ConnectionErrors, true)
	if err != nil {
		logger.LogError("PluginManagerInfoRequest.Do()", "client.doJsonRequest()", err)
		return nil, err
	}

	response := &dataModels.PluginManagerInfoResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		logger.LogError("PluginManagerInfoRequest.Do()", "json.Unmarshal()", err)
		return nil, err
	}

	return response, err
}
