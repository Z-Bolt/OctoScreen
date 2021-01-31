package dataModels

// FullStateResponse contains informantion about the current state of the printer.
type TemperatureDataResponse struct {
	// TemperatureStateResponse is the printer’s temperature state data.
	TemperatureStateResponse TemperatureStateResponse `json:"temperature"`

	// SD is the printer’s sd state data.
	// SD SDState `json:"sd"`

	// State is the printer’s general state.
	// State PrinterState `json:"state"`
}
