package dataModels


// ControlDefinition describe a system control.
type ControlDefinition struct {
	// Name of the control, will be displayed either on the button if it’s a
	// control sending a command or as a label for controls which only display
	// output.
	Name string `json:"name"`

	// Command a single GCODE command to send to the printer. Will be rendered
	// as a button which sends the command to the printer upon click. The button
	// text will be the value of the `name` attribute. Mutually exclusive with
	// `commands` and `script`. The rendered button be disabled if the printer
	// is currently offline or printing or alternatively if the requirements
	// defined via the `enabled` attribute are not met.
	Command string `json:"command"`

	// Command a list of GCODE commands to send to the printer. Will be rendered
	// as a button which sends the commands to the printer upon click. The
	// button text will be the value of the `name` attribute. Mutually exclusive
	// with `command` and `script`. The rendered button will be disabled if the
	// printer is currently offline or printing or alternatively if the
	// requirements defined via the `enabled` attribute are not met.
	Commands []string `json:"commands"`

	// Script is the name of a full blown GCODE script to send to the printer.
	// Will be rendered as a button which sends the script to the printer upon
	// click. The button text will be the value of the name attribute. Mutually
	// exclusive with `command` and `commands`. The rendered button will be
	// disabled if the printer is currently offline or printing or alternatively
	// if the requirements defined via the `enabled`` attribute are not met.
	//
	// Values of input parameters will be available in the template context
	// under the `parameter` variable (e.g. an input parameter speed will be
	// available in the script template as parameter.speed). On top of that all
	// other variables defined in the GCODE template context will be available.
	Script string `json:"script"`

	// JavaScript snippet to be executed when the button rendered for `command`
	// or `commands` is clicked. This allows to override the direct sending of
	// the command or commands to the printer with more sophisticated behavior.
	// The JavaScript snippet is `eval`’d and processed in a context where the
	// control it is part of is provided as local variable `data` and the
	// `ControlViewModel` is available as self.
	JavasScript string `json:"javascript"`

	// Enabled a JavaScript snippet to be executed when the button rendered for
	// `command` or `commands` is clicked. This allows to override the direct
	// sending of the command or commands to the printer with more sophisticated
	// behavior.  The JavaScript snippet is `eval`’d and processed in a context
	// where the control it is part of is provided as local variable `data` and
	// the `ControlViewModel` is available as `self`.
	IsEnabled bool `json:"enabled"`

	// Input a list of definitions of input parameters for a command or
	// commands, to be rendered as additional input fields.
	Input *ControlInput `json:"input"`

	// Regex a regular expression to match against lines received from the
	// printer to retrieve information from it (e.g. specific output). Together
	// with template this allows rendition of received data from the printer
	// within the UI.
	Regex string `json:"regex"`

	// Template to use for rendering the match of `regex`. May contain
	// placeholders in Python Format String Syntax[1] for either named groups
	// within the regex (e.g. `Temperature: {temperature}` for a regex
	// `T:\s*(?P<temperature>\d+(\.\d*)`) or positional groups within the regex
	// (e.g. `Position: X={0}, Y={1}, Z={2}, E={3}` for a regex
	// `X:([0-9.]+) Y:([0-9.]+) Z:([0-9.]+) E:([0-9.]+)`).
	// https://docs.python.org/2/library/string.html#format-string-syntax
	Template string `json:"template"`

	// Confirm a text to display to the user to confirm his button press. Can
	// be used with sensitive custom controls like changing EEPROM values in
	// order to prevent accidental clicks.
	Confirm string `json:"confirm"`
}
