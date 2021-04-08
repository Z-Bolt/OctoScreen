package octoprintApis

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"io"
	// "strings"

	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// ToolFlowRateRequest changes the flow rate factor to apply to extrusion of the tool.
type ToolFlowRateRequest struct {
	// Factor is the new factor, percentage as integer, between 75 and 125%.
	Factor int `json:"factor"`
}

// Do sends an API request and returns an error if any.
func (cmd *ToolFlowRateRequest) Do(c *Client) error {
	buffer := bytes.NewBuffer(nil)
	if err := cmd.encode(buffer); err != nil {
		return err
	}

	_, err := c.doJsonRequest("POST", PrinterToolApiUri, buffer, PrintToolErrors, true)
	return err
}

func (cmd *ToolFlowRateRequest) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		ToolFlowRateRequest
	}{
		Command:             "flowrate",
		ToolFlowRateRequest: *cmd,
	})
}
