package octoprintApis

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"io"
	// "strings"

	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// PrintHeadJogRequest jogs the print head (relatively) by a defined amount in
// one or more axes.
type PrintHeadJogRequest struct {
	// X is the amount distance to travel in mm or coordinate to jog print head
	// on x axis.
	X float64 `json:"x,omitempty"`

	// Y is the amount distance to travel in mm or coordinate to jog print head
	// on y axis.
	Y float64 `json:"y,omitempty"`

	// Z is the amount distance to travel in mm.or coordinate to jog print head
	// on x axis.
	Z float64 `json:"z,omitempty"`

	// Absolute is whether to move relative to current position (provided axes
	// values are relative amounts) or to absolute position (provided axes
	// values are coordinates)
	IsAbsolute bool `json:"absolute"`

	// Speed at which to move in mm/s. If not provided, minimum speed for all
	// selected axes from printer profile will be used.
	Speed int `json:"speed,omitempty"`
}

// Do sends an API request and returns an error if any.
func (cmd *PrintHeadJogRequest) Do(c *Client) error {
	buffer := bytes.NewBuffer(nil)
	if err := cmd.encode(buffer); err != nil {
		return err
	}

	_, err := c.doJsonRequest("POST", PrinterPrintHeadApiUri, buffer, PrintHeadJobErrors, true)

	return err
}

func (cmd *PrintHeadJogRequest) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		PrintHeadJogRequest
	}{
		Command:             "jog",
		PrintHeadJogRequest: *cmd,
	})
}
