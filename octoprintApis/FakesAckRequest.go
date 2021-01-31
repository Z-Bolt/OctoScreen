package octoprintApis

import (
	"bytes"
	"encoding/json"
	// "io"

	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// FakesAckRequest fakes an acknowledgment message for OctoPrint in case one got
// lost on the serial line and the communication with the printer since stalled.
//
// This should only be used in “emergencies” (e.g. to save prints), the reason
// for the lost acknowledgment should always be properly investigated and
// removed instead of depending on this “symptom solver”.
type FakesAckRequest struct{}

// Do sends an API request and returns an error if any.
func (cmd *FakesAckRequest) Do(c *Client) error {
	payload := map[string]string{"command": "fake_ack"}

	b := bytes.NewBuffer(nil)
	if err := json.NewEncoder(b).Encode(payload); err != nil {
		return err
	}

	_, err := c.doJsonRequest("POST", ConnectionApiUri, b, ConnectionErrors)
	return err
}
