package octoprint

import (
	"encoding/json"
)

const URIVersion = "/api/version"

// VersionCommand retrieve information regarding server and API version.
type VersionCommand struct{}

// Do sends an API request and returns the API response.
func (cmd *VersionCommand) Do(c *Client) (*VersionResponse, error) {
	b, err := c.doRequest("GET", URIVersion, nil)
	if err != nil {
		return nil, err
	}

	r := &VersionResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}
