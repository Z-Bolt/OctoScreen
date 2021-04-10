package octoprintApis

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"io"
	// "strings"

	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// ToolSelectRequest selects the printer’s current tool.
type ToolSelectRequest struct {
	// Tool to select, format tool{n} with n being the tool’s index starting
	// with 0.
	Tool string `json:"tool"`
}

// Do sends an API request and returns an error if any.
func (cmd *ToolSelectRequest) Do(c *Client) error {
	buffer := bytes.NewBuffer(nil)
	if err := cmd.encode(buffer); err != nil {
		return err
	}

	_, err := c.doJsonRequest("POST", PrinterToolApiUri, buffer, PrintToolErrors, true)
	return err
}

func (cmd *ToolSelectRequest) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		ToolSelectRequest
	}{
		Command:           "select",
		ToolSelectRequest: *cmd,
	})
}
