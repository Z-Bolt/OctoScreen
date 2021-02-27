package utils

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)


// ****************************************************************************
// Button Routines
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
	imageFileName := "toolhead-typeB.svg"
	if toolheadCount >= 2 {
		name = fmt.Sprintf("Tool%d", num + 1)
		imageFileName = fmt.Sprintf("toolhead-typeB-%d.svg", num + 1)
	}

	return MustButtonImageStyle(name, imageFileName, "", clicked)
}

func AttachToolheadButtonsToGrid(toolheadButtons []*gtk.Button, grid *gtk.Grid) {
	for index, toolheadButton := range toolheadButtons {
		grid.Attach(toolheadButton, index, 0, 1, 1)
	}
}






// ****************************************************************************
// DialogBox Routines
func MustConfirmDialogBox(parent *gtk.Window, msg string, cb func()) func() {
	return func() {
		win := gtk.MessageDialogNewWithMarkup(
			parent,
			gtk.DIALOG_MODAL,
			gtk.MESSAGE_INFO,
			gtk.BUTTONS_YES_NO,
			"",
		)

		win.SetMarkup(CleanHTML(msg))
		defer win.Destroy()

		box, _ := win.GetContentArea()
		box.SetMarginStart(15)
		box.SetMarginEnd(15)
		box.SetMarginTop(15)
		box.SetMarginBottom(15)

		ctx, _ := win.GetStyleContext()
		ctx.AddClass("dialog")

		if win.Run() == int(gtk.RESPONSE_YES) {
			cb()
		}
	}
}

func InfoMessageDialogBox(parentWindow *gtk.Window, message string) {
	messageDialogBox(parentWindow, gtk.MESSAGE_INFO, message)
}

func WarningMessageDialogBox(parentWindow *gtk.Window, message string) {
	messageDialogBox(parentWindow, gtk.MESSAGE_WARNING, message)
}

func ErrorMessageDialogBox(parentWindow *gtk.Window, message string) {
	messageDialogBox(parentWindow, gtk.MESSAGE_ERROR, message)
}

func messageDialogBox(parentWindow *gtk.Window, messageType gtk.MessageType, message string) {
	dialogBox := gtk.MessageDialogNewWithMarkup(
		parentWindow,
		gtk.DIALOG_MODAL,
		messageType,
		// gtk.BUTTONS_OK,
		gtk.BUTTONS_NONE,
		"",
	)

	dialogBox.AddButton("Continue", gtk.RESPONSE_OK)

	dialogBox.SetMarkup(CleanHTML(message))
	defer dialogBox.Destroy()

	box, _ := dialogBox.GetContentArea()
	box.SetMarginStart(25)
	box.SetMarginEnd(25)
	box.SetMarginTop(20)
	box.SetMarginBottom(10)

	ctx, _ := dialogBox.GetStyleContext()
	ctx.AddClass("message")

	dialogBox.Run()
}


// func hotendTemperatureIsTooLow(temperatureData octoprint.TemperatureData, action string, parentWindow *gtk.Window) bool {
// 	targetTemperature := temperatureData.Target
// 	Logger.Infof("ui_utils.HotendTemperatureIsTooLow() - targetTemperature is %.2f", targetTemperature)

// 	actualTemperature := temperatureData.Actual
// 	Logger.Infof("ui_utils.HotendTemperatureIsTooLow() - actualTemperature is %.2f", actualTemperature)

// 	if targetTemperature <= 150.0 || actualTemperature <= 150.0 {
// 		return true
// 	}

// 	return false
// }


func EmptyTheContainer(container *gtk.Container) {
	children := container.GetChildren()
	defer children.Free()

	children.Foreach(func(i interface{}) {
		container.Remove(i.(gtk.IWidget))
	})
}
