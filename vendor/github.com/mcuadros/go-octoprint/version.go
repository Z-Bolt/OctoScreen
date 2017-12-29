package octoprint

import (
	"encoding/json"
)

const URIVersion = "/api/version"

// VersionRequest retrieve information regarding server and API version.
type VersionRequest struct{}

// Do sends an API request and returns the API response.
func (cmd *VersionRequest) Do(c *Client) (*VersionResponse, error) {
	b, err := c.doJSONRequest("GET", URIVersion, nil, nil)
	if err != nil {
		return nil, err
	}

	r := &VersionResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}
