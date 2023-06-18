package octoprintApis

import (
	"encoding/json"

	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// Fetch the current selections from the FilamentManager Plugin.
type FilamentManagerSelectionsRequest struct {}

func (this *FilamentManagerSelectionsRequest) Do(client *Client) (*dataModels.FilamentManagerSelectionsResponse, error) {
	logger.TraceEnter("FilamentManagerSelectionsRequest.Do()")

	bytes, err := client.doJsonRequest("GET", FilamentManagerSelectionsUri, nil, nil, true)
	if err != nil {
		logger.LogError("FilamentManagerSelectionsRequest.Do()", "client.doJsonRequest(GET)", err)
		logger.TraceLeave("FilamentManagerSelectionsRequest.Do()")
		return nil, err
	}

	response := &dataModels.FilamentManagerSelectionsResponse {}
	if err := json.Unmarshal(bytes, response); err != nil {
		logger.LogError("FilamentManagerSelectionsRequest.Do()", "json.Unmarshal()", err)
		logger.TraceLeave("FilamentManagerSelectionsRequest.Do()")
		return nil, err
	}

	logger.TraceLeave("FilamentManagerSelectionsRequest.Do()")
	return response, err
}
