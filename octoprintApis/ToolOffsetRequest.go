package octoprintApis

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"io"
	// "strings"

	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// ToolOffsetRequest sets the given temperature offset on the printer’s tools.
type ToolOffsetRequest struct {
	// Offset is offset(s) to set, key must match the format tool{n} with n
	// being the tool’s index starting with 0.
	Offsets map[string]float64 `json:"offsets"`
}

// Do sends an API request and returns an error if any.
func (cmd *ToolOffsetRequest) Do(c *Client) error {
	buffer := bytes.NewBuffer(nil)
	if err := cmd.encode(buffer); err != nil {
		return err
	}

	_, err := c.doJsonRequest("POST", PrinterToolApiUri, buffer, PrintToolErrors, true)
	return err
}

func (cmd *ToolOffsetRequest) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		ToolOffsetRequest
	}{
		Command:           "offset",
		ToolOffsetRequest: *cmd,
	})
}
