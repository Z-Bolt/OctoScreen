package uiWidgets

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type FilamentExtrudeButton struct {
	*gtk.Button

	parentWindow				*gtk.Window
	client						*octoprintApis.Client
	amountToExtrudeStepButton	*AmountToExtrudeStepButton
	flowRateStepButton			*FlowRateStepButton // The flow rate step button is optional.
	selectExtruderStepButton	*SelectToolStepButton
	isForward					bool
}

func CreateFilamentExtrudeButton(
	parentWindow				*gtk.Window,
	client						*octoprintApis.Client,
	amountToExtrudeStepButton	*AmountToExtrudeStepButton,
	flowRateStepButton			*FlowRateStepButton, // The flow rate step button is optional.
	selectExtruderStepButton	*SelectToolStepButton,
	isForward					bool,
) *FilamentExtrudeButton {
	var base *gtk.Button
	if isForward {
		base = utils.MustButtonImageStyle("Extrude", "extruder-extrude.svg", "", nil)
	} else {
		base = utils.MustButtonImageStyle("Retract", "extruder-retract.svg", "", nil)
	}

	instance := &FilamentExtrudeButton{
		Button:						base,
		parentWindow:				parentWindow,
		client:						client,
		amountToExtrudeStepButton:	amountToExtrudeStepButton,
		flowRateStepButton:			flowRateStepButton,
		selectExtruderStepButton:	selectExtruderStepButton,
		isForward:					isForward,
	}
	instance.Button.Connect("clicked", instance.handleClicked)

	return instance
}

func (this *FilamentExtrudeButton) handleClicked() {
	this.sendExtrudeCommand(this.amountToExtrudeStepButton.Value())
}

func (this *FilamentExtrudeButton) sendExtrudeCommand(length int) {
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
		length,
	)
}
