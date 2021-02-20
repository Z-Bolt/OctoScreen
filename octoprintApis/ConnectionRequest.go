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
func (cmd *ConnectionRequest) Do(client *Client) (*dataModels.ConnectionResponse, error) {
	LogMessage("entering ConnectionRequest.Do()")

	bytes, err := client.doJsonRequest("GET", ConnectionApiUri, nil, nil)
	if err != nil {
		LogError(err, "ConnectionRequest.go, client.doJsonRequest(GET) failed")
		LogMessage("leaving ConnectionRequest.Do()")
		return nil, err
	}

	response := &dataModels.ConnectionResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		LogError(err, "ConnectionRequest.go, json.Unmarshal() failed")
		LogMessage("leaving ConnectionRequest.Do()")
		return nil, err
	}

	LogMessage("leaving ConnectionRequest.Do()")

	return response, err
}
