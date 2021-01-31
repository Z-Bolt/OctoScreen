package octoprintApis

import (
	"encoding/json"
	"fmt"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// SettingsRequest retrieves the current configuration of OctoPrint.
type PrinterProfilesRequest struct {
	Id string
}

// Do sends an API request and returns the API response.
func (cmd *PrinterProfilesRequest) Do(c *Client) (*dataModels.PrinterProfileResponse, error) {
	uri := fmt.Sprintf("/api/printerprofiles/%s", cmd.Id)
	b, err := c.doJsonRequest("GET", uri, nil, nil)
	if err != nil {
		return nil, err
	}

	r := &dataModels.PrinterProfileResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}
