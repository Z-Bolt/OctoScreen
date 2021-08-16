package ui

import (
	// "fmt"
	// "strings"
	// "time"

	"github.com/gotk3/gotk3/gtk"

	// "github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)


var filamentPanelInstance *filamentPanel

type filamentPanel struct {
	CommonPanel

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
	temperatureButton			*gtk.Button
	filamentManagerButton		*gtk.Button
	selectExtruderStepButton	*uiWidgets.SelectToolStepButton
}

func FilamentPanel(
	ui				*UI,
) *filamentPanel {
	if filamentPanelInstance == nil {
		instance := &filamentPanel {
			CommonPanel: NewCommonPanel("FilamentPanel", ui),
		}
		instance.initialize()
		filamentPanelInstance = instance
	}

	return filamentPanelInstance
}

func (this *filamentPanel) initialize() {
	defer this.Initialize()

	// Create the step buttons first, since they are needed by some of the other controls.
	this.flowRateStepButton = uiWidgets.CreateFlowRateStepButton(this.UI.Client)
	this.amountToExtrudeStepButton = uiWidgets.CreateAmountToExtrudeStepButton()
	this.selectExtruderStepButton = uiWidgets.CreateSelectExtruderStepButton(this.UI.Client, false)


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
	column := 0
	this.temperatureButton = utils.MustButtonImageStyle("Temperature", "heat-up.svg", "color1", this.showTemperaturePanel)
	this.Grid().Attach(this.temperatureButton, column, 2, 1, 1)
	column++

	// The select tool step button is needed by some of the other controls (to get the name/ID of the tool
	// to send the command to), but only display it if multiple extruders are present.
	extruderCount := utils.GetExtruderCount(this.UI.Client)
	if extruderCount > 1 {
		this.Grid().Attach(this.selectExtruderStepButton, column, 2, 1, 1)
		column++
	}

	if utils.FilamentManagerPluginIsInstalled(this.UI.Client) {
		this.addFilamentManagerButton(column)
		column++
	}
}

func (this *filamentPanel) showTemperaturePanel() {
	this.UI.GoToPanel(TemperaturePanel(this.UI))
}

func (this *filamentPanel) addFilamentManagerButton(column int) {
	// if we only have 1 tool head, skip the view for multiple toolheads
	request := &octoprintApis.FilamentManagerSelectionsRequest {}
	response, _ := request.Do(this.UI.Client)

	if response != nil {
		if len(response.Selections) == 1 {
			this.filamentManagerButton = utils.MustButtonImageStyle(
				"Filament Manager", "printing-control.svg", "",
				func() {
					this.UI.GoToPanel(FilamentManagerToolPanel(
						this.UI,
						response.Selections[0]))
				})
		} else {
			this.filamentManagerButton = utils.MustButtonImageStyle(
				"Filament Manager", "printing-control.svg", "",
				func() {
					this.UI.GoToPanel(FilamentManagerPanel(this.UI))
				})
		}
		this.Grid().Attach(this.filamentManagerButton, column, 2, 1, 1)
	}
}