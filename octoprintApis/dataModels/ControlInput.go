package dataModels


// ControlInput a list of definitions of input parameters for a command or
// commands, to be rendered as additional input fields.
type ControlInput struct {
	// Name to display for the input field.
	Name string `json:"name"`

	// Parameter name for the input field, used as a placeholder in
	// command/commands.
	Parameter string `json:"parameter"`

	// Default value for the input field.
	Default interface{} `json:"default"`

	// Slider if defined instead of an input field a slider control will be
	// rendered.
	Slider struct {
		// Minimum value of the slider, defaults to 0.
		Min int `json:"min"`

		// Maximum value of the slider, defaults to 255.
		Maximum int `json:"max"`

		// Step size per slider “tick”, defaults to 1.
		Step int `json:"step"`
	} `json:"slider"`
}
