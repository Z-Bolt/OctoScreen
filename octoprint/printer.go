package octoprint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const (
	URIPrinter   = "/api/printer"
	URIPrintHead = "/api/printer/printhead"
	URITool      = "/api/printer/tool"
	URIBed       = "/api/printer/bed"
	URICommand   = "/api/printer/command"
)

// StateCommand retrieves the current state of the printer.
type StateCommand struct {
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
func (cmd *StateCommand) Do(c *Client) (*FullStateResponse, error) {
	uri := fmt.Sprintf("%s?history=%t&limit=%d&exclude=%s", URIPrinter,
		cmd.History, cmd.Limit, strings.Join(cmd.Exclude, ","),
	)

	b, err := c.doRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	r := &FullStateResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}

// PrintHeadJogCommand jogs the print head (relatively) by a defined amount in
// one or more axes.
type PrintHeadJogCommand struct {
	// X is the amount distance to travel in mm or coordinate to jog print head
	// on x axis.
	X int `json:"x,omitempty"`
	// Y is the amount distance to travel in mm or coordinate to jog print head
	// on y axis.
	Y int `json:"y,omitempty"`
	// Z is the amount distance to travel in mm.or coordinate to jog print head
	// on x axis.
	Z int `json:"z,omitempty"`
	// Absolute is whether to move relative to current position (provided axes
	// values are relative amounts) or to absolute position (provided axes
	// values are coordinates)
	Absolute bool `json:"absolute"`
	// Speed at which to move in mm/s. If not provided, minimum speed for all
	// selected axes from printer profile will be used.
	Speed int `json:"speed,omitempty"`
}

// Do sends an API request and returns an error if any.
func (cmd *PrintHeadJogCommand) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doRequest("POST", URIPrintHead, b)
	return err
}

func (cmd *PrintHeadJogCommand) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		PrintHeadJogCommand
	}{
		Command:             "jog",
		PrintHeadJogCommand: *cmd,
	})
}

// PrintHeadHomeCommand homes the print head in all of the given axes.
type PrintHeadHomeCommand struct {
	// Axes is a list of axes which to home.
	Axes []Axis `json:"axes"`
}

// Do sends an API request and returns an error if any.
func (cmd *PrintHeadHomeCommand) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doRequest("POST", URIPrintHead, b)
	return err
}

func (cmd *PrintHeadHomeCommand) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		PrintHeadHomeCommand
	}{
		Command:              "home",
		PrintHeadHomeCommand: *cmd,
	})
}

// ToolStateCommand retrieves the current temperature data (actual, target and
// offset) plus optionally a (limited) history (actual, target, timestamp) for
// all of the printer’s available tools.
type ToolStateCommand struct {
	// History if true retrieve the temperature history.
	History bool
	// Limit limtis amount of returned history data points.
	Limit int
}

// Do sends an API request and returns the API response.
func (cmd *ToolStateCommand) Do(c *Client) (*TemperatureState, error) {
	uri := fmt.Sprintf("%s?history=%t&limit=%d", URITool, cmd.History, cmd.Limit)
	b, err := c.doRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	r := &TemperatureState{}
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}

	return r, err
}

// ToolTargetCommand sets the given target temperature on the printer’s tools.
type ToolTargetCommand struct {
	// Target temperature(s) to set, key must match the format tool{n} with n
	// being the tool’s index starting with 0.
	Target map[string]int `json:"target"`
}

// Do sends an API request and returns an error if any.
func (cmd *ToolTargetCommand) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doRequest("POST", URITool, b)
	return err
}

func (cmd *ToolTargetCommand) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		ToolTargetCommand
	}{
		Command:           "target",
		ToolTargetCommand: *cmd,
	})
}

// ToolOffsetCommand sets the given temperature offset on the printer’s tools.
type ToolOffsetCommand struct {
	// Offset is offset(s) to set, key must match the format tool{n} with n
	// being the tool’s index starting with 0.
	Offsets map[string]int `json:"offsets"`
}

// Do sends an API request and returns an error if any.
func (cmd *ToolOffsetCommand) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doRequest("POST", URITool, b)
	return err
}

func (cmd *ToolOffsetCommand) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		ToolOffsetCommand
	}{
		Command:           "offset",
		ToolOffsetCommand: *cmd,
	})
}

// ToolExtrudeCommand extrudes the given amount of filament from the currently
// selected tool.
type ToolExtrudeCommand struct {
	// Amount is the amount of filament to extrude in mm. May be negative to
	// retract.
	Amount int `json:"amount"`
}

// Do sends an API request and returns an error if any.
func (cmd *ToolExtrudeCommand) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doRequest("POST", URITool, b)
	return err
}

func (cmd *ToolExtrudeCommand) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		ToolExtrudeCommand
	}{
		Command:            "extrude",
		ToolExtrudeCommand: *cmd,
	})
}

// ToolSelectCommand selects the printer’s current tool.
type ToolSelectCommand struct {
	// Tool to select, format tool{n} with n being the tool’s index starting
	// with 0.
	Tool string `json:"tool"`
}

// Do sends an API request and returns an error if any.
func (cmd *ToolSelectCommand) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doRequest("POST", URITool, b)
	return err
}

func (cmd *ToolSelectCommand) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		ToolSelectCommand
	}{
		Command:           "select",
		ToolSelectCommand: *cmd,
	})
}

// ToolFlowrateCommand changes the flow rate factor to apply to extrusion of
// the tool.
type ToolFlowrateCommand struct {
	// Factor is the new factor, percentage as integer, between 75 and 125%.
	Factor string `json:"factor"`
}

// Do sends an API request and returns an error if any.
func (cmd *ToolFlowrateCommand) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doRequest("POST", URITool, b)
	return err
}

func (cmd *ToolFlowrateCommand) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		ToolFlowrateCommand
	}{
		Command:             "flowrate",
		ToolFlowrateCommand: *cmd,
	})
}

// BedStateCommand retrieves the current temperature data (actual, target and
// offset) plus optionally a (limited) history (actual, target, timestamp) for
// the printer’s heated bed.
//
// It’s also possible to retrieve the temperature history by supplying the
// history query parameter set to true. The amount of returned history data
// points can be limited using the limit query parameter.
type BedStateCommand struct {
	// History if true retrieve the temperature history.
	History bool
	// Limit limtis amount of returned history data points.
	Limit int
}

// Do sends an API request and returns the API response.
func (cmd *BedStateCommand) Do(c *Client) (*TemperatureState, error) {
	uri := fmt.Sprintf("%s?history=%t&limit=%d", URIBed, cmd.History, cmd.Limit)
	b, err := c.doRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	r := &TemperatureState{}
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}

	return r, err
}

// BedTargetCommand sets the given target temperature on the printer’s bed.
type BedTargetCommand struct {
	// Target temperature to set.
	Target int `json:"target"`
}

// Do sends an API request and returns an error if any.
func (cmd *BedTargetCommand) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doRequest("POST", URIBed, b)
	return err
}

func (cmd *BedTargetCommand) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		BedTargetCommand
	}{
		Command:          "target",
		BedTargetCommand: *cmd,
	})
}

// BedOffsetCommand sets the given temperature offset on the printer’s bed.
type BedOffsetCommand struct {
	// Offset is offset to set.
	Offset int `json:"offset"`
}

// Do sends an API request and returns an error if any.
func (cmd *BedOffsetCommand) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doRequest("POST", URITool, b)
	return err
}

func (cmd *BedOffsetCommand) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		BedOffsetCommand
	}{
		Command:          "offset",
		BedOffsetCommand: *cmd,
	})
}

// GCodeCommand sends any command to the printer via the serial interface.
// Should be used with some care as some commands can interfere with or even
// stop a running print job.
type GCodeCommand struct {
	// Commands list of commands to send to the printer.
	Commands []string `json:"commands"`
}

// Do sends an API request and returns an error if any.
func (cmd *GCodeCommand) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := json.NewEncoder(b).Encode(cmd); err != nil {
		return err
	}

	_, err := c.doRequest("POST", URICommand, b)
	return err
}
