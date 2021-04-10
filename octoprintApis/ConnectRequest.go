package octoprintApis

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/Z-Bolt/OctoScreen/logger"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// ConnectRequest sets the given target temperature on the printer’s tools.
type ConnectRequest struct {
	// Port specific port to connect to. If not set the current `portPreference`
	// will be used, or if no preference is available auto detection will be
	// attempted.
	Port string `json:"port,omitempty"`

	// BaudRate specific baudrate to connect with. If not set the current
	// `baudratePreference` will be used, or if no preference is available auto
	// detection will be attempted.
	BaudRate int `json:"baudrate,omitempty"`

	// PrinterProfile specific printer profile to use for connection. If not set
	// the current default printer profile will be used.
	PrinterProfile string `json:"printerProfile,omitempty"`

	// Save whether to save the request’s port and baudrate settings as new
	// preferences.
	Save bool `json:"save"`

	// Autoconnect whether to automatically connect to the printer on
	// OctoPrint’s startup in the future.
	Autoconnect bool `json:"autoconnect"`
}

// Do sends an API request and returns an error if any.
func (cmd *ConnectRequest) Do(client *Client) error {
	logger.TraceEnter("ConnectRequest.Do()")

	buffer := bytes.NewBuffer(nil)
	if err := cmd.encode(buffer); err != nil {
		logger.LogError("ConnectRequest.Do()", "cmd.encode()", err)
		logger.TraceLeave("ConnectRequest.Do()")
		return err
	}

	_, err := client.doJsonRequest("POST", ConnectionApiUri, buffer, ConnectionErrors, true)
	if err != nil {
		logger.LogError("ConnectRequest.go()", "client.doJsonRequest(POST)", err)
	}

	logger.TraceLeave("ConnectRequest.Do()")
	return err
}

func (cmd *ConnectRequest) encode(w io.Writer) error {
	payload := struct {
		Command string `json:"command"`
		ConnectRequest
	}{
		Command:        "connect",
		ConnectRequest: *cmd,
	}

	return json.NewEncoder(w).Encode(payload)
}
