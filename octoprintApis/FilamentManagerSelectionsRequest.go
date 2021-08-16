//package apis
package octoprintApis

import (
	"encoding/json"

	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)



// Fetch the current selections from the FilamentManager Plugin
type FilamentManagerSelectionsRequest struct {}

// Do sends an API request and returns the API response
func (request *FilamentManagerSelectionsRequest) Do(c *Client) (*dataModels.FilamentManagerSelections, error) {
	bytes, err := c.doJsonRequest("GET", FilamentManagerSelectionsUri, nil, nil, true)
	if err != nil {
		logger.LogError("FilamentManagerSelectionsRequest.Do()", "client.doJsonRequest(GET)", err)
		return nil, err
	}

	response := &dataModels.FilamentManagerSelections{}
	if err := json.Unmarshal(bytes, response); err != nil {
		logger.LogError("FilamentManagerSelectionsRequest.Do()", "json.Unmarshal()", err)
		return nil, err
	}

	return response, err
}
