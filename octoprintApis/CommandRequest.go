package octoprintApis

import (
	"bytes"
	"encoding/json"
	// "fmt"
	// "io"
	// "strings"

	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// CommandRequest sends any command to the printer via the serial interface.
// Should be used with some care as some commands can interfere with or even
// stop a running print job.
type CommandRequest struct {
	// Commands list of commands to send to the printer.
	Commands []string `json:"commands"`
}

// Do sends an API request and returns an error if any.
func (cmd *CommandRequest) Do(c *Client) error {
	buffer := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buffer).Encode(cmd); err != nil {
		return err
	}

	_, err := c.doJsonRequest("POST", PrinterCommandApiUri, buffer, nil, true)
	return err
}
