package octoprintApis

import (
	"encoding/json"

	"github.com/Z-Bolt/OctoScreen/logger"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// Fetch the spools from the FilamentManager Plugin.
type FilamentManagerSpoolsRequest struct {}

func (this *FilamentManagerSpoolsRequest) Do(client *Client) (*dataModels.FilamentManagerSpoolsResponse, error) {
	logger.TraceEnter("FilamentManagerSpoolsRequest.Do()")

	bytes, err := client.doJsonRequest("GET", FilamentManagerSpoolsUri, nil, nil, true)
	if err != nil {
		logger.LogError("FilamentManagerSpoolsRequest.Do()", "client.doJsonRequest(GET)", err)
		logger.TraceLeave("FilamentManagerSpoolsRequest.Do()")
		return nil, err
	}

	response := &dataModels.FilamentManagerSpoolsResponse {}
	if err := json.Unmarshal(bytes, response); err != nil {
		logger.LogError("FilamentManagerSpoolsRequest.Do()", "json.Unmarshal()", err)
		logger.TraceLeave("FilamentManagerSpoolsRequest.Do()")
		return nil, err
	}

	logger.TraceLeave("FilamentManagerSpoolsRequest.Do()")
	return response, err
}
