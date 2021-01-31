package dataModels


// TemperaturePreset describes the temperature preset for a given material.
type TemperaturePreset struct {
	Name     string  `json:"name"`
	Bed      float64 `json:"bed"`
	Extruder float64 `json:"extruder"`
}
