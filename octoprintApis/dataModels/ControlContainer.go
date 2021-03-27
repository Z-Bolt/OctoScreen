package dataModels


// ControlContainer describes a control container.
type ControlContainer struct {
	// Name to display above the container, basically a section header.
	Name string `json:"name"`

	// Gcode command
	Command string `json:"command"`

	// Script that will be run on click
	Script string `json:"script"`

	// Children a list of children controls or containers contained within this
	// container.
	Children []*ControlDefinition `json:"children"`

	// Layout  to use for laying out the contained children, either from top to
	// bottom (`vertical`) or from left to right (`horizontal``). Defaults to a
	// vertical layout.
	Layout string `json:"layout"`
}
