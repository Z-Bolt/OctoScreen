package octoprint

import "encoding/json"

const URISettings = "/api/settings"

// SettingsRequest retrieves the current configuration of OctoPrint.
type SettingsRequest struct{}

// Do sends an API request and returns the API response.
func (cmd *SettingsRequest) Do(c *Client) (*SettingsResponse, error) {
	bytes, err := c.doJSONRequest("GET", URISettings, nil, nil)
	if err != nil {
		return nil, err
	}

	response := &SettingsResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, err
}
