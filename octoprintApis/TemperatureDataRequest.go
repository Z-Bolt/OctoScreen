package octoprintApis

import (
	// "bytes"
	"encoding/json"
	"fmt"
	// "io"
	// "strings"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// FullStateRequest retrieves the current state of the printer.
type TemperatureDataRequest struct {
	// bytes if true retrieve the temperature history.
	// IncludeHistory bool

	// Limit limits the amount of returned history data points.
	// Limit int

	// Exclude list of fields to exclude from the response (e.g. if not
	// needed by the client). Valid values to supply here are `temperature`,
	// `sd` and `state`.
	// Exclude []string
}

// Do sends an API request and returns the API response.
func (cmd *TemperatureDataRequest) Do(c *Client) (*dataModels.TemperatureDataResponse, error) {
	uri := fmt.Sprintf(
		"%s?history=false&exclude=sd,state",
		URIPrinter,
	)

	// log.Printf("TODO-Remove: StateRequest (TemperatureDataResponse) uri is: %s", uri)
	//StateRequest uri is: %s /api/printer?history=false&exclude=sd,state
	/*
		{
			"temperature": {
				"bed": {
					"actual": 26.9,
					"offset": 0,
					"target": 0.0
				},
				"tool0": {
					"actual": 35.4,
					"offset": 0,
					"target": 0.0
				}
			}
		}
	*/


	bytes, err := c.doJsonRequest("GET", uri, nil, PrintErrors, true)
	if err != nil {
		return nil, err
	}

	response := &dataModels.TemperatureDataResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, err
}
