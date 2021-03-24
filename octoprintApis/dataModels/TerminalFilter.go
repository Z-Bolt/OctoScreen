package dataModels


// TerminalFilter to display in the terminal tab for filtering certain lines
// from the display terminal log.
type TerminalFilter struct {
	Name  string `json:"name"`
	RegEx string `json:"regex"`
}
