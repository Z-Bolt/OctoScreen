package octoprintApis

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"io"
	// "strings"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// PrintHeadHomeRequest homes the print head in all of the given axes.
type PrintHeadHomeRequest struct {
	// Axes is a list of axes which to home.
	Axes []dataModels.Axis `json:"axes"`
}

// Do sends an API request and returns an error if any.
func (cmd *PrintHeadHomeRequest) Do(c *Client) error {
	buffer := bytes.NewBuffer(nil)
	if err := cmd.encode(buffer); err != nil {
		return err
	}

	_, err := c.doJsonRequest("POST", PrinterPrintHeadApiUri, buffer, PrintHeadJobErrors, true)
	return err
}

func (cmd *PrintHeadHomeRequest) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		PrintHeadHomeRequest
	}{
		Command:              "home",
		PrintHeadHomeRequest: *cmd,
	})
}
