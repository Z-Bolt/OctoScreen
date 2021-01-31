package octoprintApis

import (
	// "bytes"
	"encoding/json"
	// "io"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// ConnectionRequest Retrieve the current connection settings, including
// information regarding the available baudrates and serial ports and the
// current connection state.
type ConnectionRequest struct{}

// Do sends an API request and returns the API response.
func (cmd *ConnectionRequest) Do(c *Client) (*dataModels.ConnectionResponse, error) {
	b, err := c.doJsonRequest("GET", ConnectionApiUri, nil, nil)
	if err != nil {
		return nil, err
	}

	r := &dataModels.ConnectionResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}
