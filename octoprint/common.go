package octoprint

type CurrentState struct {
	State
	Offset float64 `json:"offset"`
}

type State struct {
	Actual float64 `json:"actual"`
	Target float64 `json:"target"`
}
