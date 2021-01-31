package octoprintApis

import (
	"bytes"
	"encoding/json"
	// "io"

	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// DisconnectRequest instructs OctoPrint to disconnect from the printer.
type DisconnectRequest struct{}

// Do sends an API request and returns an error if any.
func (cmd *DisconnectRequest) Do(c *Client) error {
	payload := map[string]string{"command": "disconnect"}

	b := bytes.NewBuffer(nil)
	if err := json.NewEncoder(b).Encode(payload); err != nil {
		return err
	}

	_, err := c.doJsonRequest("POST", ConnectionApiUri, b, ConnectionErrors)
	return err
}
