package uiWidgets

import (
	// "fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)


type FilamentLoadButton struct {
	*gtk.Button

	parentWindow				*gtk.Window
	client						*octoprintApis.Client
	flowRateStepButton			*FlowRateStepButton // The flow rate step button is optional.
	selectExtruderStepButton	*SelectToolStepButton
	isForward					bool
	length						int
}

func CreateFilamentLoadButton(
	parentWindow				*gtk.Window,
	client						*octoprintApis.Client,
	flowRateStepButton			*FlowRateStepButton, // The flow rate step button is optional.
	selectExtruderStepButton	*SelectToolStepButton,
	isForward					bool,
	length						int,
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
		flowRateStepButton:			flowRateStepButton,
		selectExtruderStepButton:	selectExtruderStepButton,
		isForward:					isForward,
		length:						length,
	}
	instance.Button.Connect("clicked", instance.handleClicked)

	return instance
}

func (this *FilamentLoadButton) handleClicked() {
	this.sendLoadCommand()
}

func (this *FilamentLoadButton) sendLoadCommand() {
	extruderId := this.selectExtruderStepButton.Value()

	flowRatePercentage := 100
	if this.flowRateStepButton != nil {
		flowRatePercentage = this.flowRateStepButton.Value()
	}

	utils.Extrude(
		this.client,
		this.isForward,
		extruderId,
		this.parentWindow,
		flowRatePercentage,
		this.length,
	)
}
