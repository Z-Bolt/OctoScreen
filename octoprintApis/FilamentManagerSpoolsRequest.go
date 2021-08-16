//package apis
package octoprintApis

import (
	"encoding/json"

	"github.com/Z-Bolt/OctoScreen/logger"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)

// Fetch the spools from the FilamentManager Plugin
type FilamentManagerSpoolsRequest struct {}

func (request *FilamentManagerSpoolsRequest) Do(c *Client) (*dataModels.FilamentManagerSpools, error) {
	bytes, err := c.doJsonRequest("GET", FilamentManagerSpoolsUri, nil, nil, true)
	if err != nil {
		logger.LogError("FilamentManagerSpoolsRequest.Do()", "client.doJsonRequest(GET)", err)
		return nil, err
	}

	response := &dataModels.FilamentManagerSpools{}
	if err := json.Unmarshal(bytes, response); err != nil {
		logger.LogError("FilamentManagerSpoolsRequest.Do()", "json.Unmarshal()", err)
		return nil, err
	}

	return response, err
}
