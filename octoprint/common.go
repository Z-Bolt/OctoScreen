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

type CurrentState struct {
	State
	Offset float64 `json:"offset"`
}

type State struct {
	Actual float64 `json:"actual"`
	Target float64 `json:"target"`
}

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

type SDState struct {
	Ready bool `json:"ready"`
}

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
