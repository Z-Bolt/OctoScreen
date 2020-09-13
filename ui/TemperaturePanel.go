package ui

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)

var temperaturePanelInstance *temperaturePanel

type temperaturePanel struct {
	CommonPanel

	// First row
	decreaseButton					*uiWidgets.TemperatureIncreaseButton
	temperatureAmountStepButton		*uiWidgets.TemperatureAmountStepButton
	increaseButton					*uiWidgets.TemperatureIncreaseButton

	// Second row
	temperatureStatusBox			*uiWidgets.TemperatureStatusBox

	// Third row
	presetsButton					*gtk.Button
	selectToolStepButton			*uiWidgets.SelectToolStepButton
}

func TemperaturePanel(
	ui				*UI,
	parentPanel		interfaces.IPanel,
) *temperaturePanel {
	if temperaturePanelInstance == nil {
		temperaturePanelInstance = &temperaturePanel{
			CommonPanel:	NewCommonPanel(ui, parentPanel),
		}
		temperaturePanelInstance.initialize()
	}

	return temperaturePanelInstance
}

func (this *temperaturePanel) initialize() {
	defer this.Initialize()

	// Create the step buttons first, since they are needed by some of the other controls.
	this.temperatureAmountStepButton = uiWidgets.CreateTemperatureAmountStepButton()
	this.selectToolStepButton = uiWidgets.CreateSelectToolStepButton(this.UI.Printer, true)


	// First row
	this.decreaseButton = uiWidgets.CreateTemperatureIncreaseButton(
		this.UI.Printer,
		this.temperatureAmountStepButton,
		this.selectToolStepButton,
		false,
	)
	this.Grid().Attach(this.decreaseButton, 0, 0, 1, 1)

	this.Grid().Attach(this.temperatureAmountStepButton, 1, 0, 1, 1)

	this.increaseButton = uiWidgets.CreateTemperatureIncreaseButton(
		this.UI.Printer,
		this.temperatureAmountStepButton,
		this.selectToolStepButton,
		true,
	)
	this.Grid().Attach(this.increaseButton, 2, 0, 1, 1)


	// Second row
	this.temperatureStatusBox = uiWidgets.CreateTemperatureStatusBox(this.UI.Printer, true, true)
	this.Grid().Attach(this.temperatureStatusBox, 1, 1, 2, 1)


	// Third row
	this.presetsButton = utils.MustButtonImageStyle("Presets", "heat-up.svg",  "color2", this.showTemperaturePresetsPanel)
	this.Grid().Attach(this.presetsButton, 0, 2, 1, 1)

	this.Grid().Attach(this.selectToolStepButton, 1, 2, 1, 1)
}

func (this *temperaturePanel) showTemperaturePresetsPanel() {
	temperaturePresetsPanel := TemperaturePresetsPanel(this.UI, this, this.selectToolStepButton)
	this.UI.GoToPanel(temperaturePresetsPanel)
}
