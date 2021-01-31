package dataModels

type PauseAction string

const (
	// Pause the current job if it’s printing, does nothing if it’s already paused.
	Pause PauseAction = "pause"

	// Resume the current job if it’s paused, does nothing if it’s printing.
	Resume PauseAction = "resume"

	// Toggle the pause state of the job, pausing it if it’s printing and
	// resuming it if it’s currently paused.
	Toggle PauseAction = "toggle"
)
