package octoprint

import (
	"encoding/json"
	"time"
)

type Axis string

const (
	XAxis Axis = "x"
	YAxis Axis = "y"
	ZAxis Axis = "z"
)

// CurrentState of a tool.
type CurrentState struct {
	State
	Offset float64 `json:"offset"`
}

// State of a tool.
type State struct {
	Actual float64 `json:"actual"`
	Target float64 `json:"target"`
}

// PrinterState current state of the printer.
type PrinterState struct {
	Text  string `json:"text"`
	Flags struct {
		Operations    bool `json:"operational"`
		Puased        bool `json:"paused"`
		Printing      bool `json:"printing"`
		SDReady       bool `json:"sdReady"`
		Error         bool `json:"error"`
		Ready         bool `json:"ready"`
		ClosedOnError bool `json:"closedOrError"`
	} `json:"flags"`
}

// SDState is the state of the sd reader.
type SDState struct {
	Ready bool `json:"ready"`
}

// History of a tool.
type History history
type history struct {
	Time  time.Time        `json:"time"`
	Tools map[string]State `json:"tools"`
}

func (h *History) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	ts := raw["time"]
	delete(raw, "time")
	b, _ = json.Marshal(map[string]interface{}{
		"time":  time.Unix(int64(ts.(float64)), 0),
		"tools": raw,
	})

	i := &history{}
	if err := json.Unmarshal(b, i); err != nil {
		return err
	}

	*h = History(*i)
	return nil
}
