package octoprint

import (
	// "encoding/json"
	// "strconv"
	// "strings"
	// "time"
)

// TemperatureData is temperature stats for a tool.
type TemperatureData struct {
	// Actual current temperature.
	Actual float64 `json:"actual"`

	// Target temperature, may be nil if no target temperature is set.
	Target float64 `json:"target"`

	// Offset currently configured temperature offset to apply, will be left
	// out for historic temperature information.
	Offset float64 `json:"offset"`
}
