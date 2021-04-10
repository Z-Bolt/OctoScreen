package octoprintApis

import (
	// "bytes"
	"encoding/json"
	"fmt"
	// "io"
	// "strings"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// ToolStateRequest retrieves the current temperature data (actual, target and
// offset) plus optionally a (limited) history (actual, target, timestamp) for
// all of the printer’s available tools.
type ToolStateRequest struct {
	// History if true retrieve the temperature history.
	IncludeHistory bool

	// Limit limtis amount of returned history data points.
	Limit int
}

// Do sends an API request and returns the API response.
func (cmd *ToolStateRequest) Do(c *Client) (*dataModels.TemperatureStateResponse, error) {
	uri := fmt.Sprintf("%s?history=%t&limit=%d", PrinterToolApiUri, cmd.IncludeHistory, cmd.Limit)

	// log.Printf("TODO-Remove: ToolStateRequest uri is: %s", uri)
	//ToolStateRequest uri is: %s /api/printer/tool?history=true&limit=1
	/*
		{
			"history": [
				{
					"tool0": {
						"actual": 38.0,
						"target": 0.0
					}
				}
			],
			"tool0": {
				"actual": 38.0,
				"offset": 0,
				"target": 0.0
			}
		}
	*/

	bytes, err := c.doJsonRequest("GET", uri, nil, nil, true)
	if err != nil {
		return nil, err
	}

	response := &dataModels.TemperatureStateResponse{}
	if err := json.Unmarshal(bytes, &response); err != nil {
		return nil, err
	}

	return response, err
}
