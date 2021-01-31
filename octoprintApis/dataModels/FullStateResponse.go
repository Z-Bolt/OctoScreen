package dataModels

// FullStateResponse contains informantion about the current state of the printer.
type FullStateResponse struct {
	// TemperatureStateResponse is the printer’s temperature state data.
	Temperature TemperatureStateResponse `json:"temperature"`

	// SD is the printer’s sd state data.
	SD SdState `json:"sd"`

	// State is the printer’s general state.
	State PrinterState `json:"state"`
}
