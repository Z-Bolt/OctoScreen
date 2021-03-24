package dataModels


// CommandSource is the source of the command definition.
type CommandSource string

const (
	// Core for system actions defined by OctoPrint itself.
	Core CommandSource = "core"

	// Custom for custom system commands defined by the user through `config.yaml`.
	Custom CommandSource = "custom"
)
