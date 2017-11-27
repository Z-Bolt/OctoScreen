package octoprint

import (
	"encoding/json"
	"fmt"
	"strings"
)

const PrinterTool = "/api/printer"

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

type StateResponse struct {
	Temperature ToolResponse `json:"temperature"`
	SD SDState `json:"sd"`
	State PrinterState `json:"state"`
}

func (cmd *StateCommand) Do(p *Printer) (*StateResponse, error) {
	uri := fmt.Sprintf("%s?history=%t&limit=%d&exclude=%s", PrinterTool,
		cmd.History, cmd.Limit, strings.Join(cmd.Exclude, ","),
	)

	b, err := p.doRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	r := &StateResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}
