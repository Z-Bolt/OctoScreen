package octoprintApis

import (
	"encoding/json"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


const SettingsApiUri = "/api/settings"


// SettingsRequest retrieves the current configuration of OctoPrint.
type SettingsRequest struct{}

// Do sends an API request and returns the API response.
func (cmd *SettingsRequest) Do(c *Client) (*dataModels.SettingsResponse, error) {
	bytes, err := c.doJsonRequest("GET", SettingsApiUri, nil, nil, true)
	if err != nil {
		return nil, err
	}

	response := &dataModels.SettingsResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, err
}
