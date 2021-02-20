package octoprintApis

import (
	"bytes"
	"encoding/json"
	"io"

	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// TODO: DisconnectRequest doesn't seem to be used anywhere... maybe remove it?

// DisconnectRequest instructs OctoPrint to disconnect from the printer.
type DisconnectRequest struct{}

// Do sends an API request and returns an error if any.
func (cmd *DisconnectRequest) Do(client *Client) error {
	LogMessage("entering DisconnectRequest.Do()")

	bytes := bytes.NewBuffer(nil)
	if err := cmd.encode(bytes); err != nil {
		LogError(err, "DisconnectRequest.go, cmd.encode() failed")
		LogMessage("leaving DisconnectRequest.Do()")
		return err
	}

	_, err := client.doJsonRequest("POST", ConnectionApiUri, bytes, ConnectionErrors)
	if err != nil {
		LogError(err, "DisconnectRequest.go, client.doJsonRequest(POST) failed")
	}

	LogMessage("leaving DisconnectRequest.Do()")
	return err
}

func (cmd *DisconnectRequest) encode(w io.Writer) error {
	payload := map[string]string {
		"command": "disconnect",
	}

	return json.NewEncoder(w).Encode(payload)
}
