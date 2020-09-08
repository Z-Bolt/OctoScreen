package uiWidgets

import (
	// "fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type CommandButton struct {
	*gtk.Button

	client				*octoprint.Client
	parentWindow		*gtk.Window
	commandDefinition	*octoprint.CommandDefinition
}

func CreateCommandButton(
	client				*octoprint.Client,
	parentWindow		*gtk.Window,
	commandDefinition	*octoprint.CommandDefinition,
	iconName			string,
) *CommandButton {
	base := utils.MustButtonImage(utils.StrEllipsisLen(commandDefinition.Name, 16), iconName + ".svg", nil)
	instance := &CommandButton {
		Button:				base,
		client:				client,
		parentWindow:		parentWindow,
		commandDefinition:	commandDefinition,
	}
	_, err := instance.Button.Connect("clicked", instance.handleClicked)
	if err != nil {
		panic(err)
	}

	return instance
}

func (this *CommandButton) handleClicked() {
	if len(this.commandDefinition.Confirm) != 0 {
		utils.MustConfirmDialogBox(this.parentWindow, this.commandDefinition.Confirm, this.sendCommand)
		return
	} else {
		this.sendCommand()
	}
}

func (this *CommandButton) sendCommand() {
	commandRequest := &octoprint.SystemExecuteCommandRequest{
		Source: octoprint.Custom,
		Action: this.commandDefinition.Action,
	}

	err := commandRequest.Do(this.client)
	if err != nil {
		utils.LogError("CommandButton.sendCommand()", "Do(SystemExecuteCommandRequest)", err)
	}
}
