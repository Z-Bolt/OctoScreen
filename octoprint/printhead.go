package octoprint

import (
	"bytes"
	"encoding/json"
	"io"
)

const URIPrintHead = "/api/printer/printhead"

// JogCommand jogs the print head (relatively) by a defined amount in one or
// more axes.
type JogCommand struct {
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
func (cmd *JogCommand) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doRequest("POST", URIPrintHead, b)
	return err
}

func (cmd *JogCommand) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		JogCommand
	}{
		Command:    "jog",
		JogCommand: *cmd,
	})
}

// HomeCommand homes the print head in all of the given axes.
type HomeCommand struct {
	// Axes is a list of axes which to home.
	Axes []Axis `json:"axes"`
}

// Do sends an API request and returns an error if any.
func (cmd *HomeCommand) Do(c *Client) error {
	b := bytes.NewBuffer(nil)
	if err := cmd.encode(b); err != nil {
		return err
	}

	_, err := c.doRequest("POST", URIPrintHead, b)
	return err
}

func (cmd *HomeCommand) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		HomeCommand
	}{
		Command:     "home",
		HomeCommand: *cmd,
	})
}
