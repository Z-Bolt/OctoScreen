package uiWidgets

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type TemperatureIncreaseButton struct {
	*gtk.Button

	client							*octoprintApis.Client
	temperatureAmountStepButton		*TemperatureAmountStepButton
	selectHotendStepButton			*SelectToolStepButton
	isIncrease						bool
}

func CreateTemperatureIncreaseButton(
	client							*octoprintApis.Client,
	temperatureAmountStepButton		*TemperatureAmountStepButton,
	selectHotendStepButton			*SelectToolStepButton,
	isIncrease						bool,
) *TemperatureIncreaseButton {
	var base *gtk.Button
	if isIncrease {
		base = utils.MustButtonImageStyle("Increase", "increase.svg", "", nil)
	} else {
		base = utils.MustButtonImageStyle("Decrease", "decrease.svg", "", nil)
	}

	instance := &TemperatureIncreaseButton{
		Button:							base,
		client:							client,
		temperatureAmountStepButton:	temperatureAmountStepButton,
		selectHotendStepButton:			selectHotendStepButton,
		isIncrease:						isIncrease,
	}
	instance.Button.Connect("clicked", instance.handleClicked)

	return instance
}

func (this *TemperatureIncreaseButton) handleClicked() {
	value := this.temperatureAmountStepButton.Value()
	tool := this.selectHotendStepButton.Value()
	target, err := utils.GetToolTarget(this.client, tool)
	if err != nil {
		logger.LogError("TemperatureIncreaseButton.handleClicked()", "GetToolTarget()", err)
		return
	}

	if this.isIncrease {
		target += value
	} else {
		target -= value
	}

	if target < 0 {
		target = 0
	}

	// TODO: should the target be checked for a max temp?
	// If so, how to calculate what the max should be?

	logger.Infof("TemperatureIncreaseButton.handleClicked() - setting target temperature for %s to %1.f°C.", tool, target)

	err = utils.SetToolTarget(this.client, tool, target)
	if err != nil {
		logger.LogError("TemperatureIncreaseButton.handleClicked()", "GetToolTarget()", err)
	}
}
