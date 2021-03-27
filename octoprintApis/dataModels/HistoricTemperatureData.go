package dataModels

import (
	"encoding/json"
)

// HistoricTemperatureData is temperature historic stats for a tool.
type HistoricTemperatureData historicTemperatureData

type historicTemperatureData struct {
	// Time of this data point.
	Time JsonTime `json:"time"`

	// Tools is temperature stats a set of tools.
	Tools map[string]TemperatureData `json:"tools"`
}

func (h *HistoricTemperatureData) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	ts := raw["time"]
	delete(raw, "time")
	b, _ = json.Marshal(map[string]interface{}{
		"time":  ts,
		"tools": raw,
	})

	i := &historicTemperatureData{}
	if err := json.Unmarshal(b, i); err != nil {
		return err
	}

	*h = HistoricTemperatureData(*i)
	return nil
}
