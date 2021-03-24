package dataModels

import (
	"encoding/json"
)


// CommandDefinition describe a system command.
type CommandDefinition struct {
	// Name of the command to display in the System menu.
	Name string `json:"name"`

	// Command is the full command line to execute for the command.
	Command string `json:"command"`

	// Action is an identifier to refer to the command programmatically. The
	// special `action` string divider signifies a `divider` in the menu.
	Action string `json:"action"`

	// Confirm if present and set, this text will be displayed to the user in a
	// confirmation dialog they have to acknowledge in order to really execute
	// the command.
	RawConfirm json.RawMessage `json:"confirm"`
	Confirm    string          `json:"-"`

	// Async whether to execute the command asynchronously or wait for its
	// result before responding to the HTTP execution request.
	IsAsync bool `json:"async"`

	// Ignore whether to ignore the return code of the command’s execution.
	Ignore bool `json:"ignore"`

	// Source of the command definition.
	Source CommandSource `json:"source"`

	// Resource is the URL of the command to use for executing it.
	Resource string `json:"resource"`
}
