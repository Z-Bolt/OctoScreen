package ui

import (
	// "fmt"
	// "strings"
	// "time"

	"github.com/gotk3/gotk3/gtk"

	// "github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)


type filamentPanel struct {
	CommonPanel

	backgroundTask				*utils.BackgroundTask

	// First row
	filamentExtrudeButton		*uiWidgets.FilamentExtrudeButton
	flowRateStepButton			*uiWidgets.FlowRateStepButton
	amountToExtrudeStepButton	*uiWidgets.AmountToExtrudeStepButton
	filamentRetractButton		*uiWidgets.FilamentExtrudeButton

	// Second row
	filamentLoadButton			*uiWidgets.FilamentLoadButton
	temperatureStatusBox		*uiWidgets.TemperatureStatusBox
	filamentUnloadButton		*uiWidgets.FilamentLoadButton

	// Third row
	selectExtruderStepButton	*uiWidgets.SelectToolStepButton
	temperatureButton			*gtk.Button
	filamentManagerButton		*gtk.Button
}

var filamentPanelInstance *filamentPanel

func GetFilamentPanelInstance(
	ui				*UI,
) *filamentPanel {
	if filamentPanelInstance == nil {
		filamentPanelInstance = &filamentPanel {
			CommonPanel: CreateCommonPanel("FilamentPanel", ui),
		}
		filamentPanelInstance.initialize()
		filamentPanelInstance.createBackgroundTask()
	}

	return filamentPanelInstance
}

func (this *filamentPanel) initialize() {
	defer this.Initialize()

	// Create the step buttons first, since they are needed by some of the other controls.
	this.flowRateStepButton = uiWidgets.CreateFlowRateStepButton(this.UI.Client, 1, nil)
	this.amountToExtrudeStepButton = uiWidgets.CreateAmountToExtrudeStepButton(2, nil)
	this.selectExtruderStepButton = uiWidgets.CreateSelectExtruderStepButton(this.UI.Client, false, 3, nil)

	// First row
	this.filamentExtrudeButton = uiWidgets.CreateFilamentExtrudeButton(
		this.UI.window,
		this.UI.Client,
		this.amountToExtrudeStepButton,
		this.flowRateStepButton,
		this.selectExtruderStepButton,
		true,
	)
	this.Grid().Attach(this.filamentExtrudeButton,		0, 0, 1, 1)

	this.Grid().Attach(this.flowRateStepButton,			1, 0, 1, 1)

	this.Grid().Attach(this.amountToExtrudeStepButton,	2, 0, 1, 1)

	this.filamentRetractButton = uiWidgets.CreateFilamentExtrudeButton(
		this.UI.window,
		this.UI.Client,
		this.amountToExtrudeStepButton,
		this.flowRateStepButton,
		this.selectExtruderStepButton,
		false,
	)
	this.Grid().Attach(this.filamentRetractButton,		3, 0, 1, 1)


	// Second row
	this.filamentLoadButton = uiWidgets.CreateFilamentLoadButton(
		this.UI.window,
		this.UI.Client,
		this.flowRateStepButton,
		this.selectExtruderStepButton,
		true,
		int(this.UI.Settings.FilamentInLength),
	)
	this.Grid().Attach(this.filamentLoadButton,			0, 1, 1, 1)

	this.temperatureStatusBox = uiWidgets.CreateTemperatureStatusBox(this.UI.Client, true, true)
	this.Grid().Attach(this.temperatureStatusBox,		1, 1, 2, 1)

	this.filamentUnloadButton = uiWidgets.CreateFilamentLoadButton(
		this.UI.window,
		this.UI.Client,
		this.flowRateStepButton,
		this.selectExtruderStepButton,
		false,
		int(this.UI.Settings.FilamentOutLength),
	)
	this.Grid().Attach(this.filamentUnloadButton,		3, 1, 1, 1)


	// Third row
	// The select tool step button is needed by some of the other controls (to get the name/ID of the tool
	// to send the command to), but only display it if multiple extruders are present.
	extruderCount := utils.GetExtruderCount(this.UI.Client)
	if extruderCount > 1 {
		this.Grid().Attach(this.selectExtruderStepButton, 0, 2, 1, 1)
	}

	this.addTemperatureButton(1)

	if utils.FilamentManagerPluginIsInstalled(this.UI.Client) {
		this.addFilamentManagerButton(2)
	}
}

func (this *filamentPanel) createBackgroundTask() {
	logger.TraceEnter("FilamentPanel.createBackgroundTask()")

	// Default timeout of 1 second.
	duration := utils.GetExperimentalFrequency(1, "EXPERIMENTAL_IDLE_UPDATE_FREQUENCY")
	this.backgroundTask = utils.CreateBackgroundTask(duration, this.update)
	// Update the UI every second, but the data is only updated once every 10 seconds.
	// See OctoPrintResponseManager.update(). 
	this.backgroundTask.Start()

	logger.TraceLeave("FilamentPanel.createBackgroundTask()")
}

func (this *filamentPanel) update() {
	logger.TraceEnter("FilamentPanel.update()")

	this.updateTemperature()

	logger.TraceLeave("FilamentPanel.update()")
}

func (this *filamentPanel) updateTemperature() {
	logger.TraceEnter("FilamentPanel.updateTemperature()")

	octoPrintResponseManager := GetOctoPrintResponseManagerInstance(this.UI)
	if octoPrintResponseManager.IsConnected() != true {
		// If not connected, do nothing and leave.
		logger.TraceLeave("FilamentPanel.updateTemperature() (not connected)")
		return
	}

	this.temperatureStatusBox.UpdateTemperatureData(octoPrintResponseManager.FullStateResponse.Temperature.CurrentTemperatureData)

	logger.TraceLeave("FilamentPanel.updateTemperature()")
}

func (this *filamentPanel) showTemperaturePanel() {
	this.UI.GoToPanel(GetTemperaturePanelInstance(this.UI))
}

func (this *filamentPanel) addTemperatureButton(column int) {
	this.temperatureButton = utils.MustButtonImageStyle("Temperature", "heat-up.svg", "color1", this.showTemperaturePanel)
	this.Grid().Attach(this.temperatureButton, column, 2, 1, 1)
}

func (this *filamentPanel) addFilamentManagerButton(column int) {
	// If we only have 1 tool head, skip the view for multiple toolheads.
	filamentManagerSelectionsRequest := &octoprintApis.FilamentManagerSelectionsRequest {}
	filamentManagerSelectionsResponse, _ := filamentManagerSelectionsRequest.Do(this.UI.Client)

	filamentManagerSpoolsRequest := &octoprintApis.FilamentManagerSpoolsRequest {}
	filamentManagerSpoolsResponse, _ := filamentManagerSpoolsRequest.Do(this.UI.Client)


	if filamentManagerSelectionsResponse != nil {
		this.filamentManagerButton = utils.MustButtonImageStyle(
			"Filament Manager", "filament-manager.svg", "color2",
			func() {
				filamentManagerPanelInstance := GetFilamentManagerPanelInstance(
					this.UI,
					filamentManagerSelectionsResponse,
					filamentManagerSpoolsResponse,
				)
				this.UI.GoToPanel(filamentManagerPanelInstance)
			})
		this.Grid().Attach(this.filamentManagerButton, column, 2, 1, 1)
	}
}
