package octoprint

import "encoding/json"

const URISettings = "/api/settings"

// SettingsRequest retrieves the current configuration of OctoPrint.
type SettingsRequest struct{}

// Do sends an API request and returns the API response.
func (cmd *SettingsRequest) Do(c *Client) (*Settings, error) {
	b, err := c.doJSONRequest("GET", URISettings, nil, nil)
	if err != nil {
		return nil, err
	}

	r := &Settings{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}
