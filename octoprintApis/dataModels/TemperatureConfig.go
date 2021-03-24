package dataModels


// TemperatureConfig temperature profiles which will be displayed in the
// temperature tab.
type TemperatureConfig struct {
	// Graph cutoff in minutes.
	Cutoff int `json:"cutoff"`

	// Profiles  which will be displayed in the temperature tab.
	TemperaturePresets []*TemperaturePreset `json:"profiles"`

	// SendAutomatically enable this to have temperature fine adjustments you
	// do via the + or - button be sent to the printer automatically.
	SendAutomatically bool `json:"sendAutomatically"`

	// SendAutomaticallyAfter OctoPrint will use this delay to limit the number
	// of sent temperature commands should you perform multiple fine adjustments
	// in a short time.
	SendAutomaticallyAfter float64 `json:"sendAutomaticallyAfter"`
}
