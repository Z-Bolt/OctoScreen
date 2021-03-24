package octoprintApis

import (
	// "bytes"
	"encoding/json"
	// "fmt"
	// "io"
	// "strings"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// CustomCommandsRequest retrieves all configured system controls.
type CustomCommandsRequest struct{}

// Do sends an API request and returns the API response.
func (cmd *CustomCommandsRequest) Do(c *Client) (*dataModels.CustomCommandsResponse, error) {
	b, err := c.doJsonRequest("GET", PrinterCommandCustomApiUri, nil, nil)
	if err != nil {
		return nil, err
	}

	r := &dataModels.CustomCommandsResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}
