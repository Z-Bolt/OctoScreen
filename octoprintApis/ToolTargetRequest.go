package octoprintApis

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"io"
	// "strings"

	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// ToolTargetRequest sets the given target temperature on the printer’s tools.
type ToolTargetRequest struct {
	// Target temperature(s) to set, key must match the format tool{n} with n
	// being the tool’s index starting with 0.
	Targets map[string]float64 `json:"targets"`
}

// Do sends an API request and returns an error if any.
func (cmd *ToolTargetRequest) Do(c *Client) error {
	buffer := bytes.NewBuffer(nil)
	if err := cmd.encode(buffer); err != nil {
		return err
	}

	_, err := c.doJsonRequest("POST", PrinterToolApiUri, buffer, PrintToolErrors, true)
	return err
}

func (cmd *ToolTargetRequest) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		ToolTargetRequest
	}{
		Command:           "target",
		ToolTargetRequest: *cmd,
	})
}
