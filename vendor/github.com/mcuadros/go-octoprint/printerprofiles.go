package octoprint

import (
	"encoding/json"
	"fmt"
)

// SettingsRequest retrieves the current configuration of OctoPrint.
type PrinterProfilesRequest struct {
	Id string
}

// Do sends an API request and returns the API response.
func (cmd *PrinterProfilesRequest) Do(c *Client) (*PrinterProfile, error) {
	uri := fmt.Sprintf("/api/printerprofiles/%s", cmd.Id)
	b, err := c.doJSONRequest("GET", uri, nil, nil)
	if err != nil {
		return nil, err
	}

	r := &PrinterProfile{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}
