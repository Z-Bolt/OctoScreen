package uiWidgets

import (
	// "fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type ControlButton struct {
	*gtk.Button

	client				*octoprint.Client
	parentWindow		*gtk.Window
	controlDefinition	*octoprint.ControlDefinition
}

func CreateControlButton(
	client				*octoprint.Client,
	parentWindow		*gtk.Window,
	controlDefinition	*octoprint.ControlDefinition,
	iconName			string,
) *ControlButton {
	base := utils.MustButtonImage(utils.StrEllipsisLen(controlDefinition.Name, 16), iconName + ".svg", nil)
	instance := &ControlButton {
		Button:				base,
		client:				client,
		parentWindow:		parentWindow,
		controlDefinition:	controlDefinition,
	}
	_, err := instance.Button.Connect("clicked", instance.handleClicked)
	if err != nil {
		panic(err)
	}

	return instance
}

func (this *ControlButton) handleClicked() {
	if len(this.controlDefinition.Confirm) != 0 {
		utils.MustConfirmDialogBox(this.parentWindow, this.controlDefinition.Confirm, this.sendCommand)
		return
	} else {
		this.sendCommand()
	}
}

func (this *ControlButton) sendCommand() {
	commandRequest := &octoprint.CommandRequest{
		Commands: this.controlDefinition.Commands,
	}

	if len(this.controlDefinition.Command) != 0 {
		commandRequest.Commands = []string{this.controlDefinition.Command}
	}

	utils.Logger.Infof("Executing command %q", this.controlDefinition.Name)
	err := commandRequest.Do(this.client)
	if err != nil {
		utils.LogError("ControlButton.sendCommand()", "Do(CommandRequest)", err)
	}
}
