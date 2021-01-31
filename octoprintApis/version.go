package octoprintApis

import (
	"encoding/json"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


const VersionApiUri = "/api/version"


// VersionRequest retrieve information regarding server and API version.
type VersionRequest struct{}

// Do sends an API request and returns the API response.
func (cmd *VersionRequest) Do(c *Client) (*dataModels.VersionResponse, error) {
	b, err := c.doJsonRequest("GET", VersionApiUri, nil, nil)
	if err != nil {
		return nil, err
	}

	r := &dataModels.VersionResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}
