package octoprintApis

import (
	// "bytes"
	"encoding/json"
	"fmt"
	// "io"
	// "strings"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// BedStateRequest retrieves the current temperature data (actual, target and
// offset) plus optionally a (limited) history (actual, target, timestamp) for
// the printer’s heated bed.
//
// It’s also possible to retrieve the temperature history by supplying the
// history query parameter set to true. The amount of returned history data
// points can be limited using the limit query parameter.
type BedStateRequest struct {
	// History if true retrieve the temperature history.
	IncludeHistory bool

	// Limit limtis amount of returned history data points.
	Limit int
}

// Do sends an API request and returns the API response.
func (cmd *BedStateRequest) Do(c *Client) (*dataModels.TemperatureStateResponse, error) {
	uri := fmt.Sprintf("%s?history=%t&limit=%d", PrinterBedApiUri, cmd.IncludeHistory, cmd.Limit)
	bytes, err := c.doJsonRequest("GET", uri, nil, PrintBedErrors, true)
	if err != nil {
		return nil, err
	}

	response := &dataModels.TemperatureStateResponse{}
	if err := json.Unmarshal(bytes, &response); err != nil {
		return nil, err
	}

	return response, err
}
