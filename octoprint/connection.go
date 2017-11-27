package octoprint

import (
	"bytes"
	"encoding/json"
	"io"
)

const URIConnection = "/api/connection"

// ConnectionCommand Retrieve the current connection settings, including
// information regarding the available baudrates and serial ports and the
// current connection state.
type ConnectionCommand struct{}

// Do sends an API request and returns the API response.
func (cmd *ConnectionCommand) Do(c *Client) (*ConnectionResponse, error) {
	b, err := c.doRequest("GET", URIVersion, nil)
	if err != nil {
		return nil, err
	}

	r := &ConnectionResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}

// ConnectCommand sets the given target temperature on the printer’s tools.
type ConnectCommand struct {
	// Port specific port to connect to. If not set the current `portPreference`
	// will be used, or if no preference is available auto detection will be
	// attempted.
	Port string `json:"port"`
	// BaudRate specific baudrate to connect with. If not set the current
	// `baudratePreference` will be used, or if no preference is available auto
	// detection will be attempted.
	BaudRate int `json:"baudrate"`
	// PrinterProfile specific printer profile to use for connection. If not set
	// the current default printer profile will be used.
	PrinterProfile string `json:"printerProfile"`
	// Save whether to save the request’s port and baudrate settings as new
	// preferences.
	Save bool `json:"save"`
	// Autoconnect whether to automatically connect to the printer on
	// OctoPrint’s startup in the future.
	Autoconnect bool `json:"autoconnect"`
}

// Do sends an API request and returns an error if any.
func (cmd *ConnectCommand) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doRequest("POST", URIConnection, b)
	return err
}

func (cmd *ConnectCommand) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		ConnectCommand
	}{
		Command:        "connect",
		ConnectCommand: *cmd,
	})
}

// DisconnectCommand instructs OctoPrint to disconnect from the printer.
type DisconnectCommand struct{}

// Do sends an API request and returns an error if any.
func (cmd *DisconnectCommand) Do(c *Client) error {
	payload := map[string]string{"command": "disconnect"}

	b := bytes.NewBuffer(nil)
	if err := json.NewEncoder(b).Encode(payload); err != nil {
		return err
	}

	_, err := c.doRequest("POST", URIConnection, b)
	return err
}

// FakesACKCommand fakes an acknowledgment message for OctoPrint in case one got
// lost on the serial line and the communication with the printer since stalled.
//
// This should only be used in “emergencies” (e.g. to save prints), the reason
// for the lost acknowledgment should always be properly investigated and
// removed instead of depending on this “symptom solver”.
type FakesACKCommand struct{}

// Do sends an API request and returns an error if any.
func (cmd *FakesACKCommand) Do(c *Client) error {
	payload := map[string]string{"command": "fake_ack"}

	b := bytes.NewBuffer(nil)
	if err := json.NewEncoder(b).Encode(payload); err != nil {
		return err
	}

	_, err := c.doRequest("POST", URIConnection, b)
	return err
}
