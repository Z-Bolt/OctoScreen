package octoprintApis

import (
	"bytes"
	"encoding/json"
	"io"

	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)

// TODO: FakesAckRequest doesn't seem to be used anywhere... maybe remove it?

// FakesAckRequest fakes an acknowledgment message for OctoPrint in case one got
// lost on the serial line and the communication with the printer since stalled.
//
// This should only be used in “emergencies” (e.g. to save prints), the reason
// for the lost acknowledgment should always be properly investigated and
// removed instead of depending on this “symptom solver”.
type FakesAckRequest struct{}

// Do sends an API request and returns an error if any.
func (cmd *FakesAckRequest) Do(client *Client) error {
	LogMessage("entering FakesAckRequest.Do()")

	bytes := bytes.NewBuffer(nil)
	if err := cmd.encode(bytes); err != nil {
		LogError(err, "FakesAckRequest.go, cmd.encode() failed")
		LogMessage("leaving FakesAckRequest.Do()")
		return err
	}

	_, err := client.doJsonRequest("POST", ConnectionApiUri, bytes, ConnectionErrors)
	if err != nil {
		LogError(err, "FakesAckRequest.go, client.doJsonRequest(POST) failed")
	}

	LogMessage("leaving FakesAckRequest.Do()")
	return err
}

func (cmd *FakesAckRequest) encode(w io.Writer) error {
	payload := map[string]string {
		"command": "fake_ack",
	}

	return json.NewEncoder(w).Encode(payload)
}
