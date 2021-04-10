package octoprintApis

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"io"
	// "strings"

	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// BedOffsetRequest sets the given temperature offset on the printer’s bed.
type BedOffsetRequest struct {
	// Offset is offset to set.
	Offset int `json:"offset"`
}

// Do sends an API request and returns an error if any.
func (cmd *BedOffsetRequest) Do(c *Client) error {
	buffer := bytes.NewBuffer(nil)
	if err := cmd.encode(buffer); err != nil {
		return err
	}

	_, err := c.doJsonRequest("POST", PrinterToolApiUri, buffer, PrintToolErrors, true)
	return err
}

func (cmd *BedOffsetRequest) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		BedOffsetRequest
	}{
		Command:          "offset",
		BedOffsetRequest: *cmd,
	})
}
