package ui

import (
	"github.com/gotk3/gotk3/gtk"

	// "github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)


type temperaturePanel struct {
	CommonPanel

	backgroundTask					*utils.BackgroundTask

	// First row
	decreaseButton					*uiWidgets.TemperatureIncreaseButton
	selectHotendStepButton			*uiWidgets.SelectToolStepButton
	temperatureAmountStepButton		*uiWidgets.TemperatureAmountStepButton
	increaseButton					*uiWidgets.TemperatureIncreaseButton

	// Second row
	coolDownButton					*uiWidgets.CoolDownButton
	temperatureStatusBox			*uiWidgets.TemperatureStatusBox

	// Third row
	presetsButton					*gtk.Button
}

var temperaturePanelInstance *temperaturePanel

func GetTemperaturePanelInstance(
	ui				*UI,
) *temperaturePanel {
	if temperaturePanelInstance == nil {
		temperaturePanelInstance = &temperaturePanel {
			CommonPanel: CreateCommonPanel("TemperaturePanel", ui),
		}
		temperaturePanelInstance.initialize()
		temperaturePanelInstance.createBackgroundTask()
	}

	return temperaturePanelInstance
}

func (this *temperaturePanel) initialize() {
	defer this.Initialize()

	// Create the step buttons first, since they are needed by some of the other controls.
	this.temperatureAmountStepButton = uiWidgets.CreateTemperatureAmountStepButton()
	this.selectHotendStepButton = uiWidgets.CreateSelectHotendStepButton(this.UI.Client, true)


	// First row
	this.decreaseButton = uiWidgets.CreateTemperatureIncreaseButton(
		this.UI.Client,
		this.temperatureAmountStepButton,
		this.selectHotendStepButton,
		false,
	)
	this.Grid().Attach(this.decreaseButton, 0, 0, 1, 1)

	this.Grid().Attach(this.selectHotendStepButton, 1, 0, 1, 1)

	this.Grid().Attach(this.temperatureAmountStepButton, 2, 0, 1, 1)

	this.increaseButton = uiWidgets.CreateTemperatureIncreaseButton(
		this.UI.Client,
		this.temperatureAmountStepButton,
		this.selectHotendStepButton,
		true,
	)
	this.Grid().Attach(this.increaseButton, 3, 0, 1, 1)


	// Second row
	this.coolDownButton = uiWidgets.CreateCoolDownButton(this.UI.Client, nil)
	this.Grid().Attach(this.coolDownButton, 0, 1, 1, 1)

	this.temperatureStatusBox = uiWidgets.CreateTemperatureStatusBox(this.UI.Client, true, true)
	this.Grid().Attach(this.temperatureStatusBox, 1, 1, 2, 1)


	// Third row
	this.presetsButton = utils.MustButtonImageStyle("Presets", "heat-up.svg",  "color2", this.showTemperaturePresetsPanel)
	this.Grid().Attach(this.presetsButton, 0, 2, 1, 1)
}

func (this *temperaturePanel) createBackgroundTask() {
	logger.TraceEnter("TemperaturePanel.createBackgroundTask()")

	// Default timeout of 1 second.
	duration := utils.GetExperimentalFrequency(1, "EXPERIMENTAL_IDLE_UPDATE_FREQUENCY")
	this.backgroundTask = utils.CreateBackgroundTask(duration, this.update)
	// Update the UI every second, but the data is only updated once every 10 seconds.
	// See OctoPrintResponseManager.update(). 
	this.backgroundTask.Start()

	logger.TraceLeave("TemperaturePanel.createBackgroundTask()")
}

func (this *temperaturePanel) update() {
	logger.TraceEnter("TemperaturePanel.update()")

	this.updateTemperature()

	logger.TraceLeave("TemperaturePanel.update()")
}

func (this *temperaturePanel) updateTemperature() {
	logger.TraceEnter("TemperaturePanel.updateTemperature()")

	octoPrintResponseManager := GetOctoPrintResponseManagerInstance(this.UI)
	if octoPrintResponseManager.IsConnected() != true {
		// If not connected, do nothing and leave.
		logger.TraceLeave("TemperaturePanel.updateTemperature() (not connected)")
		return
	}

	this.temperatureStatusBox.UpdateTemperatureData(octoPrintResponseManager.FullStateResponse.Temperature.CurrentTemperatureData)

	logger.TraceLeave("TemperaturePanel.updateTemperature()")
}

func (this *temperaturePanel) showTemperaturePresetsPanel() {
	this.UI.GoToPanel(GetTemperaturePresetsPanelInstance(this.UI, this.selectHotendStepButton))
}
