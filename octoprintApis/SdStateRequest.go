package octoprintApis

import (
	// "bytes"
	"encoding/json"
	// "fmt"
	// "io"
	// "strings"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// SdStateRequest retrieves the current state of the printer’s SD card. For this
// request no authentication is needed.
type SdStateRequest struct{}

// Do sends an API request and returns the API response.
func (cmd *SdStateRequest) Do(c *Client) (*dataModels.SdState, error) {
	b, err := c.doJsonRequest("GET", PrinterSdApiUri, nil, PrintSdErrors)
	if err != nil {
		return nil, err
	}

	r := &dataModels.SdState{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}
