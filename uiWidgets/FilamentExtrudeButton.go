package uiWidgets

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type FilamentExtrudeButton struct {
	*gtk.Button

	parentWindow				*gtk.Window
	client						*octoprintApis.Client
	amountToExtrudeStepButton	*AmountToExtrudeStepButton
	flowRateStepButton			*FlowRateStepButton // The flow rate step button is optional.
	selectToolStepButton		*SelectToolStepButton
	isForward					bool
}

func CreateFilamentExtrudeButton(
	parentWindow				*gtk.Window,
	client						*octoprintApis.Client,
	amountToExtrudeStepButton	*AmountToExtrudeStepButton,
	flowRateStepButton			*FlowRateStepButton, // The flow rate step button is optional.
	selectToolStepButton		*SelectToolStepButton,
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
		selectToolStepButton:		selectToolStepButton,
		isForward:					isForward,
	}
	_, err := instance.Button.Connect("clicked", instance.handleClicked)
	if err != nil {
		panic(err)
	}

	return instance
}

func (this *FilamentExtrudeButton) handleClicked() {
	this.sendExtrudeCommand(this.amountToExtrudeStepButton.Value())
}

func (this *FilamentExtrudeButton) sendExtrudeCommand(amount int) {
	// The flow rate step button is optional.
	if this.flowRateStepButton != nil {
		err := this.flowRateStepButton.SendChangeFlowRate()
		if err != nil {
			utils.LogError("FilamentExtrudeButton.sendExtrudeCommand()", "SendChangeFlowRate()", err)
			return
		}
	}

	extruderId := this.selectToolStepButton.Value()
	var action string
	if this.isForward {
		action = "extrude"
	} else {
		action = "retract"
	}
	if utils.CurrentHotendTemperatureIsTooLow(this.client, extruderId, action, this.parentWindow) {
		utils.Logger.Error("FilamentExtrudeButton.sendExtrudeCommand() -  temperature is too low")
		return
	}

	cmd := &octoprintApis.ToolExtrudeRequest{}
	if this.isForward {
		cmd.Amount = amount
	} else {
		cmd.Amount = -amount
	}

	utils.Logger.Infof("FilamentExtrudeButton.sendExtrudeCommand() - sending extrude request with amount: %d", cmd.Amount)
	if err := cmd.Do(this.client); err != nil {
		utils.LogError("FilamentExtrudeButton.sendExtrudeCommand()", "Do(ToolExtrudeRequest)", err)
	}
}
