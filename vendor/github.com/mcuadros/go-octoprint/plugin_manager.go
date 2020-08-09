package octoprint

import (
	"bytes"
	"encoding/json"
	"log"
)

const pluginManagerRequestURI = "/api/plugin/pluginmanager"


// Plugin -
type Plugin struct {
	Author                     string   `json:"author"`
	Blacklisted                bool     `json:"blacklisted"`
	Bundled                    bool     `json:"bundled"`
	Description                string   `json:"description"`
	// DisablingDiscouraged       bool     `json:"disabling_discouraged"`
	Enabled                    bool     `json:"enabled"`
	ForcedDisabled             bool     `json:"forced_disabled"`
	Incompatible               bool     `json:"incompatible"`
	Key                        string   `json:"key"`
	License                    string   `json:"license"`
	Managable                  bool     `json:"managable"`
	Name                       string   `json:"name"`
	// notifications: []
	Origin                     string   `json:"origin"`
	PendingDisable             bool     `json:"pending_disable"`
	PendingEnable              bool     `json:"pending_enable"`
	PendingInstall             bool     `json:"pending_install"`
	PendingUninstall           bool     `json:"pending_uninstall"`
	Python                     string   `json:"python"`
	SafeModeVictim             bool     `json:"safe_mode_victim"`
	URL                        string   `json:"url"`
	Version                    string   `json:"version"`
}

// GetPluginManagerInfoRequest -
type GetPluginManagerInfoRequest struct {
	Command string `json:"command"`
}

// GetPluginManagerInfoResponse -
type GetPluginManagerInfoResponse struct {
	Octoprint     string        `json:"octoprint"`
    Online        bool          `json:"online"`
    //orphan_data: { }
	OS            string        `json:"os"`
	//pip: {}
	Plugins       []Plugin      `json:"plugins"`
}

// Do -
func (cmd *GetPluginManagerInfoRequest) Do(c *Client) (*GetPluginManagerInfoResponse, error) {
	cmd.Command = "get_settings"

	params := bytes.NewBuffer(nil)
	if err := json.NewEncoder(params).Encode(cmd); err != nil {
		log.Println("plugin_manager.Do() - Encode() failed")
		return nil, err
	}

	b, err := c.doJSONRequest("GET", pluginManagerRequestURI, params, ConnectionErrors)
	if err != nil {
		log.Println("plugin_manager.Do() - doJSONRequest() failed")
		return nil, err
	}

	r := &GetPluginManagerInfoResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		log.Println("plugin_manager.Do() - Unmarshal() failed")
		return nil, err
	}

	return r, err
}
