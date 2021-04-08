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
	bytes, err := c.doJsonRequest("GET", VersionApiUri, nil, nil, true)
	if err != nil {
		return nil, err
	}

	response := &dataModels.VersionResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, err
}
