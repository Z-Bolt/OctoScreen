package octoprint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const (
	URIPrinter       = "/api/printer"
	URIPrintHead     = "/api/printer/printhead"
	URIPrintTool     = "/api/printer/tool"
	URIPrintBed      = "/api/printer/bed"
	URIPrintSD       = "/api/printer/sd"
	URICommand       = "/api/printer/command"
	URICommandCustom = "/api/printer/command/custom"
)

var (
	PrintErrors = statusMapping{
		409: "Printer is not operational",
	}
	PrintHeadJobErrors = statusMapping{
		400: "Invalid axis specified, invalid value for travel amount for a jog command or factor for feed rate or otherwise invalid request",
		409: "Printer is not operational or currently printing",
	}
	PrintToolErrors = statusMapping{
		400: "Targets or offsets contains a property or tool contains a value not matching the format tool{n}, the target/offset temperature, extrusion amount or flow rate factor is not a valid number or outside of the supported range, or if the request is otherwise invalid",
		409: "Printer is not operational",
	}
	PrintBedErrors = statusMapping{
		409: "Printer is not operational or the selected printer profile does not have a heated bed.",
	}
	PrintSDErrors = statusMapping{
		404: "SD support has been disabled in OctoPrint’s settings.",
		409: "SD card has not been initialized.",
	}
)

// StateRequest retrieves the current state of the printer.
type StateRequest struct {
	// History if true retrieve the temperature history.
	History bool
	// Limit limtis amount of returned history data points.
	Limit int
	// Exclude list of fields to exclude from the response (e.g. if not
	// needed by the client). Valid values to supply here are `temperature`,
	// `sd` and `state`.
	Exclude []string
}

// Do sends an API request and returns the API response.
func (cmd *StateRequest) Do(c *Client) (*FullStateResponse, error) {
	uri := fmt.Sprintf("%s?history=%t&limit=%d&exclude=%s", URIPrinter,
		cmd.History, cmd.Limit, strings.Join(cmd.Exclude, ","),
	)

	b, err := c.doJSONRequest("GET", uri, nil, PrintErrors)
	if err != nil {
		return nil, err
	}

	r := &FullStateResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}

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
	Absolute bool `json:"absolute"`
	// Speed at which to move in mm/s. If not provided, minimum speed for all
	// selected axes from printer profile will be used.
	Speed int `json:"speed,omitempty"`
}

// Do sends an API request and returns an error if any.
func (cmd *PrintHeadJogRequest) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doJSONRequest("POST", URIPrintHead, b, PrintHeadJobErrors)

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

// PrintHeadHomeRequest homes the print head in all of the given axes.
type PrintHeadHomeRequest struct {
	// Axes is a list of axes which to home.
	Axes []Axis `json:"axes"`
}

// Do sends an API request and returns an error if any.
func (cmd *PrintHeadHomeRequest) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doJSONRequest("POST", URIPrintHead, b, PrintHeadJobErrors)
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

// ToolStateRequest retrieves the current temperature data (actual, target and
// offset) plus optionally a (limited) history (actual, target, timestamp) for
// all of the printer’s available tools.
type ToolStateRequest struct {
	// History if true retrieve the temperature history.
	History bool
	// Limit limtis amount of returned history data points.
	Limit int
}

// Do sends an API request and returns the API response.
func (cmd *ToolStateRequest) Do(c *Client) (*TemperatureState, error) {
	uri := fmt.Sprintf("%s?history=%t&limit=%d", URIPrintTool, cmd.History, cmd.Limit)
	b, err := c.doJSONRequest("GET", uri, nil, nil)
	if err != nil {
		return nil, err
	}

	r := &TemperatureState{}
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}

	return r, err
}

// ToolTargetRequest sets the given target temperature on the printer’s tools.
type ToolTargetRequest struct {
	// Target temperature(s) to set, key must match the format tool{n} with n
	// being the tool’s index starting with 0.
	Targets map[string]float64 `json:"targets"`
}

// Do sends an API request and returns an error if any.
func (cmd *ToolTargetRequest) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doJSONRequest("POST", URIPrintTool, b, PrintToolErrors)
	return err
}

func (cmd *ToolTargetRequest) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		ToolTargetRequest
	}{
		Command:           "target",
		ToolTargetRequest: *cmd,
	})
}

// ToolOffsetRequest sets the given temperature offset on the printer’s tools.
type ToolOffsetRequest struct {
	// Offset is offset(s) to set, key must match the format tool{n} with n
	// being the tool’s index starting with 0.
	Offsets map[string]float64 `json:"offsets"`
}

// Do sends an API request and returns an error if any.
func (cmd *ToolOffsetRequest) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doJSONRequest("POST", URIPrintTool, b, PrintToolErrors)
	return err
}

func (cmd *ToolOffsetRequest) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		ToolOffsetRequest
	}{
		Command:           "offset",
		ToolOffsetRequest: *cmd,
	})
}

// ToolExtrudeRequest extrudes the given amount of filament from the currently
// selected tool.
type ToolExtrudeRequest struct {
	// Amount is the amount of filament to extrude in mm. May be negative to
	// retract.
	Amount int `json:"amount"`
}

// Do sends an API request and returns an error if any.
func (cmd *ToolExtrudeRequest) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doJSONRequest("POST", URIPrintTool, b, PrintToolErrors)
	return err
}

func (cmd *ToolExtrudeRequest) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		ToolExtrudeRequest
	}{
		Command:            "extrude",
		ToolExtrudeRequest: *cmd,
	})
}

// ToolSelectRequest selects the printer’s current tool.
type ToolSelectRequest struct {
	// Tool to select, format tool{n} with n being the tool’s index starting
	// with 0.
	Tool string `json:"tool"`
}

// Do sends an API request and returns an error if any.
func (cmd *ToolSelectRequest) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doJSONRequest("POST", URIPrintTool, b, PrintToolErrors)
	return err
}

func (cmd *ToolSelectRequest) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		ToolSelectRequest
	}{
		Command:           "select",
		ToolSelectRequest: *cmd,
	})
}

// ToolFlowrateRequest changes the flow rate factor to apply to extrusion of
// the tool.
type ToolFlowrateRequest struct {
	// Factor is the new factor, percentage as integer, between 75 and 125%.
	Factor int `json:"factor"`
}

// Do sends an API request and returns an error if any.
func (cmd *ToolFlowrateRequest) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doJSONRequest("POST", URIPrintTool, b, PrintToolErrors)
	return err
}

func (cmd *ToolFlowrateRequest) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		ToolFlowrateRequest
	}{
		Command:             "flowrate",
		ToolFlowrateRequest: *cmd,
	})
}

// BedStateRequest retrieves the current temperature data (actual, target and
// offset) plus optionally a (limited) history (actual, target, timestamp) for
// the printer’s heated bed.
//
// It’s also possible to retrieve the temperature history by supplying the
// history query parameter set to true. The amount of returned history data
// points can be limited using the limit query parameter.
type BedStateRequest struct {
	// History if true retrieve the temperature history.
	History bool
	// Limit limtis amount of returned history data points.
	Limit int
}

// Do sends an API request and returns the API response.
func (cmd *BedStateRequest) Do(c *Client) (*TemperatureState, error) {
	uri := fmt.Sprintf("%s?history=%t&limit=%d", URIPrintBed, cmd.History, cmd.Limit)
	b, err := c.doJSONRequest("GET", uri, nil, PrintBedErrors)
	if err != nil {
		return nil, err
	}

	r := &TemperatureState{}
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}

	return r, err
}

// BedTargetRequest sets the given target temperature on the printer’s bed.
type BedTargetRequest struct {
	// Target temperature to set.
	Target float64 `json:"target"`
}

// Do sends an API request and returns an error if any.
func (cmd *BedTargetRequest) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doJSONRequest("POST", URIPrintBed, b, PrintBedErrors)
	return err
}

func (cmd *BedTargetRequest) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		BedTargetRequest
	}{
		Command:          "target",
		BedTargetRequest: *cmd,
	})
}

// BedOffsetRequest sets the given temperature offset on the printer’s bed.
type BedOffsetRequest struct {
	// Offset is offset to set.
	Offset int `json:"offset"`
}

// Do sends an API request and returns an error if any.
func (cmd *BedOffsetRequest) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doJSONRequest("POST", URIPrintTool, b, PrintToolErrors)
	return err
}

func (cmd *BedOffsetRequest) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		BedOffsetRequest
	}{
		Command:          "offset",
		BedOffsetRequest: *cmd,
	})
}

// CommandRequest sends any command to the printer via the serial interface.
// Should be used with some care as some commands can interfere with or even
// stop a running print job.
type CommandRequest struct {
	// Commands list of commands to send to the printer.
	Commands []string `json:"commands"`
}

// Do sends an API request and returns an error if any.
func (cmd *CommandRequest) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := json.NewEncoder(b).Encode(cmd); err != nil {
		return err
	}

	_, err := c.doJSONRequest("POST", URICommand, b, nil)
	return err
}

// CustomCommandsRequest retrieves all configured system controls.
type CustomCommandsRequest struct{}

// Do sends an API request and returns the API response.
func (cmd *CustomCommandsRequest) Do(c *Client) (*CustomCommandsResponse, error) {
	b, err := c.doJSONRequest("GET", URICommandCustom, nil, nil)
	if err != nil {
		return nil, err
	}

	r := &CustomCommandsResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}

// SDStateRequest retrieves the current state of the printer’s SD card. For this
// request no authentication is needed.
type SDStateRequest struct{}

// Do sends an API request and returns the API response.
func (cmd *SDStateRequest) Do(c *Client) (*SDState, error) {
	b, err := c.doJSONRequest("GET", URIPrintSD, nil, PrintSDErrors)
	if err != nil {
		return nil, err
	}

	r := &SDState{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}

// SDInitRequest initializes the printer’s SD card, making it available for use.
// This also includes an initial retrieval of the list of files currently stored
// on the SD card.
type SDInitRequest struct{}

// Do sends an API request and returns an error if any.
func (cmd *SDInitRequest) Do(c *Client) error {
	return doCommandRequest(c, URIPrintSD, "init", PrintSDErrors)
}

// SDRefreshRequest Refreshes the list of files stored on the printer’s SD card.
type SDRefreshRequest struct{}

// Do sends an API request and returns an error if any.
func (cmd *SDRefreshRequest) Do(c *Client) error {
	return doCommandRequest(c, URIPrintSD, "refresh", PrintSDErrors)
}

// SDReleaseRequest releases the SD card from the printer. The reverse operation
// to init. After issuing this command, the SD card won’t be available anymore,
// hence and operations targeting files stored on it will fail.
type SDReleaseRequest struct{}

// Do sends an API request and returns an error if any.
func (cmd *SDReleaseRequest) Do(c *Client) error {
	return doCommandRequest(c, URIPrintSD, "release", PrintSDErrors)
}

// doCommandRequest can be used in any operation where the only required field
// is the `command` field.
func doCommandRequest(c *Client, uri, command string, m statusMapping) error {
	v := map[string]string{"command": command}

	b := bytes.NewBuffer(nil)
	if err := json.NewEncoder(b).Encode(v); err != nil {
		return err
	}

	_, err := c.doJSONRequest("POST", uri, b, m)
	return err
}
