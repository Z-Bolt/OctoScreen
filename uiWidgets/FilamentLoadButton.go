package uiWidgets

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type FilamentLoadButton struct {
	*gtk.Button

	parentWindow				*gtk.Window
	client						*octoprint.Client
	selectToolStepButton		*SelectToolStepButton
	isForward					bool
}

func CreateFilamentLoadButton(
	parentWindow				*gtk.Window,
	client						*octoprint.Client,
	selectToolStepButton		*SelectToolStepButton,
	isForward					bool,
) *FilamentLoadButton {
	var base *gtk.Button
	if isForward {
		base = utils.MustButtonImageStyle("Load", "filament-spool-load.svg", "", nil)
	} else {
		base = utils.MustButtonImageStyle("Unload", "filament-spool-unload.svg", "", nil)
	}

	instance := &FilamentLoadButton{
		Button:						base,
		parentWindow:				parentWindow,
		client:						client,
		selectToolStepButton:		selectToolStepButton,
		isForward:					isForward,
	}
	_, err := instance.Button.Connect("clicked", instance.handleClicked)
	if err != nil {
		panic(err)
	}

	return instance
}

func (this *FilamentLoadButton) handleClicked() {
	this.sendLoadCommand()
}

func (this *FilamentLoadButton) sendLoadCommand() {
	extruderId := this.selectToolStepButton.Value()
	var action string
	if this.isForward {
		action = "load"
	} else {
		action = "unload"
	}
	if utils.CurrentHotendTemperatureIsTooLow(this.client, extruderId, action, this.parentWindow) {
		utils.Logger.Error("FilamentLoadButton.sendLoadCommand() -  temperature is too low")
		return
	}

	// BUG: This does not work.  At least not on a Prusa i3.  Need to get this working with all printers.
	cmd := &octoprint.CommandRequest{}
	if this.isForward {
		cmd.Commands = []string{"G91", "G0 E600 F5000", "G0 E120 F500", "G90"}
	} else {
		cmd.Commands = []string{"G91", "G0 E-800 F5000", "G90"}
	}

	utils.Logger.Info("FilamentLoadButton.sendLoadCommand() - sending filament load request")
	if err := cmd.Do(this.client); err != nil {
		utils.LogError("FilamentLoadButton.sendLoadCommand()", "Do(CommandRequest)", err)
		return
	}
}
