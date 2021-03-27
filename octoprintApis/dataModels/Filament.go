package dataModels


type Filament struct {
	// Length estimated of filament used, in mm
	Length uint32 `json:"length"`

	// Volume estimated of filament used, in cm³
	Volume float64 `json:"volume"`
}
