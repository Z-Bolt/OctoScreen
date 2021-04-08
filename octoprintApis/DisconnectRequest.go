package octoprintApis

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/Z-Bolt/OctoScreen/logger"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// TODO: DisconnectRequest doesn't seem to be used anywhere... maybe remove it?

// DisconnectRequest instructs OctoPrint to disconnect from the printer.
type DisconnectRequest struct{}

// Do sends an API request and returns an error if any.
func (this *DisconnectRequest) Do(client *Client) error {
	logger.TraceEnter("DisconnectRequest.Do()")

	buffer := bytes.NewBuffer(nil)
	if err := this.encode(buffer); err != nil {
		logger.LogError("DisconnectRequest.Do()", "this.encode(bytes)", err)
		logger.TraceLeave("DisconnectRequest.Do()")
		return err
	}

	_, err := client.doJsonRequest("POST", ConnectionApiUri, buffer, ConnectionErrors, true)
	if err != nil {
		logger.LogError("DisconnectRequest.Do()", "client.doJsonRequest(POST)", err)
	}

	logger.TraceLeave("DisconnectRequest.Do()")
	return err
}

func (cmd *DisconnectRequest) encode(w io.Writer) error {
	payload := map[string]string {
		"command": "disconnect",
	}

	return json.NewEncoder(w).Encode(payload)
}
