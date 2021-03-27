package octoprint

import (
	"encoding/json"
	// "strconv"
	// "strings"
	// "time"
)


// TODO: add request
// TODO: add Do()


// TemperatureState is the printer’s temperature state data.
type TemperatureStateResponse temperatureStateResponse

type temperatureStateResponse struct {
	// Current temperature stats.
	CurrentTemperatureData map[string]TemperatureData `json:"current"`

	// Temperature history.
	History []*HistoricTemperatureData `json:"history"`
}

func (r *TemperatureStateResponse) UnmarshalJSON(bytes []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	history := raw["history"]
	delete(raw, "history")
	bytes, _ = json.Marshal(map[string]interface{}{
		"current": raw,
		"history": history,
	})

	i := &temperatureStateResponse{}
	if err := json.Unmarshal(bytes, i); err != nil {
		return err
	}

	*r = TemperatureStateResponse(*i)

	return nil
}
