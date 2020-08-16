/*
I would have preferred that this code was in the utils folder, but due to Go's opinionated way
of coding and structure, that's not possible, so this file will house UI-based utility functions.
*/

package ui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)


func CreateToolheadButtonsAndAttachToGrid(toolheadCount int, grid *gtk.Grid) []*gtk.Button {
	toolheadButtons := CreateToolheadButtons(toolheadCount)
	AttachToolheadButtonsToGrid(toolheadButtons, grid)

	return toolheadButtons
}

func CreateChangeToolheadButtonsAndAttachToGrid(toolheadCount int, grid *gtk.Grid) []*gtk.Button {
	toolheadButtons := CreateToolheadButtons(toolheadCount)
	for index, toolheadButton := range toolheadButtons {
		toolheadButton.SetLabel(fmt.Sprintf("Change to Tool%d", index + 1))
	}

	AttachToolheadButtonsToGrid(toolheadButtons, grid)

	return toolheadButtons
}

func CreateToolheadButtons(toolheadCount int) []*gtk.Button {
	var toolheadButtons []*gtk.Button
	var toolheadButton *gtk.Button

	toolheadButton = CreateToolheadButton(0, toolheadCount, func() { })
	toolheadButtons = append(toolheadButtons, toolheadButton)

	if toolheadCount >= 2 {
		toolheadButton = CreateToolheadButton(1, toolheadCount, func() { })
		toolheadButtons = append(toolheadButtons, toolheadButton)
	}

	if toolheadCount >= 3 {
		toolheadButton = CreateToolheadButton(2, toolheadCount, func() { })
		toolheadButtons = append(toolheadButtons, toolheadButton)
	}

	if toolheadCount >= 4 {
		toolheadButton = CreateToolheadButton(3, toolheadCount, func() { })
		toolheadButtons = append(toolheadButtons, toolheadButton)
	}

	return toolheadButtons
}

func CreateToolheadButton(num, toolheadCount int, clicked func()) *gtk.Button {
	name := ""
	imageFileName := "toolhead.svg"
	if toolheadCount >= 2 {
		name = fmt.Sprintf("Tool%d", num + 1)
		imageFileName = fmt.Sprintf("toolhead-with-color-%d.svg", num + 1)
	}

	return MustButtonImageStyle(name, imageFileName, "", clicked)
}

func AttachToolheadButtonsToGrid(toolheadButtons []*gtk.Button, grid *gtk.Grid) {
	for index, toolheadButton := range toolheadButtons {
		grid.Attach(toolheadButton, index, 0, 1, 1)
	}
}
