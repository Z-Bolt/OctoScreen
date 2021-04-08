package octoprintApis

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"io"
	// "strings"

	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// BedTargetRequest sets the given target temperature on the printer’s bed.
type BedTargetRequest struct {
	// Target temperature to set.
	Target float64 `json:"target"`
}

// Do sends an API request and returns an error if any.
func (cmd *BedTargetRequest) Do(c *Client) error {
	buffer := bytes.NewBuffer(nil)
	if err := cmd.encode(buffer); err != nil {
		return err
	}

	_, err := c.doJsonRequest("POST", PrinterBedApiUri, buffer, PrintBedErrors, true)
	return err
}

func (cmd *BedTargetRequest) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		BedTargetRequest
	}{
		Command:          "target",
		BedTargetRequest: *cmd,
	})
}
