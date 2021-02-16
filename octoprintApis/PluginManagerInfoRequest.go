package octoprintApis

import (
	"bytes"
	"encoding/json"
	"log"

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
		log.Println("plugin_manager.Do() - Encode() failed")
		return nil, err
	}

	b, err := client.doJsonRequest("GET", pluginManagerRequestURI, params, ConnectionErrors)
	if err != nil {
		log.Println("plugin_manager.Do() - doJsonRequest() failed")
		return nil, err
	}

	r := &dataModels.PluginManagerInfoResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		log.Println("plugin_manager.Do() - Unmarshal() failed")
		return nil, err
	}

	return r, err
}
