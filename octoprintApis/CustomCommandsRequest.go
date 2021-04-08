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
	bytes, err := c.doJsonRequest("GET", PrinterCommandCustomApiUri, nil, nil, true)
	if err != nil {
		return nil, err
	}

	response := &dataModels.CustomCommandsResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, err
}
