package octoprintApis

import (
	// "bytes"
	"encoding/json"
	// "io"

	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// ConnectionRequest Retrieve the current connection settings, including
// information regarding the available baudrates and serial ports and the
// current connection state.
type ConnectionRequest struct{}

// Do sends an API request and returns the API response.
func (cmd *ConnectionRequest) Do(client *Client) (*dataModels.ConnectionResponse, error) {
	logger.TraceEnter("ConnectionRequest.Do()")

	bytes, err := client.doJsonRequest("GET", ConnectionApiUri, nil, nil, true)
	if err != nil {
		logger.LogError("ConnectionRequest.Do()", "client.doJsonRequest(GET)", err)
		logger.TraceLeave("ConnectionRequest.Do()")
		return nil, err
	}

	response := &dataModels.ConnectionResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		logger.LogError("ConnectionRequest.Do()", "json.Unmarshal()", err)
		logger.TraceLeave("ConnectionRequest.Do()")
		return nil, err
	}

	logger.TraceLeave("ConnectionRequest.Do()")
	return response, err
}
