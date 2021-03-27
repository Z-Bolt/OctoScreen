package uiWidgets

import (
	// "fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type CommandButton struct {
	*gtk.Button

	client				*octoprintApis.Client
	parentWindow		*gtk.Window
	commandDefinition	*dataModels.CommandDefinition
}

func CreateCommandButton(
	client				*octoprintApis.Client,
	parentWindow		*gtk.Window,
	commandDefinition	*dataModels.CommandDefinition,
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
		logger.LogError("PANIC!!! - CreateCommandButton()", "instance.Button.Connect()", err)
		panic(err)
	}

	return instance
}

func (this *CommandButton) handleClicked() {
	var handleClick func()
	
	if len(this.commandDefinition.Confirm) != 0 {
		logger.Debugf("CommandButton.handleClicked().MustConfirmDialogBox: %s", this.commandDefinition.Name)
		handleClick = utils.MustConfirmDialogBox(this.parentWindow, this.commandDefinition.Confirm, this.sendCommand)
	} else {
		logger.Debugf("CommandButton.handleClicked().sendCommand: %s", this.commandDefinition.Name)
		handleClick = this.sendCommand
	}
	
	handleClick()
}

func (this *CommandButton) sendCommand() {
	logger.Debugf("CommandButton.sendCommand(): %q", this.commandDefinition.Name)
	
	commandRequest := &octoprintApis.SystemExecuteCommandRequest{
		Source: dataModels.Custom,
		Action: this.commandDefinition.Action,
	}
	
	err := commandRequest.Do(this.client)
	if err != nil {
		logger.LogError("CommandButton.sendCommand()", "Do(SystemExecuteCommandRequest)", err)
	}
}
