package ui

import (
	// "fmt"
	// "strings"
	// "time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/interfaces"
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
	selectToolStepButton		*uiWidgets.SelectToolStepButton
}

func FilamentPanel(
	ui				*UI,
	parentPanel		interfaces.IPanel,
) *filamentPanel {
	if filamentPanelInstance == nil {
		instance := &filamentPanel {
			CommonPanel: NewCommonPanel(ui, parentPanel),
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
	this.selectToolStepButton = uiWidgets.CreateSelectToolStepButton(this.UI.Client, false)


	// First row
	this.filamentExtrudeButton = uiWidgets.CreateFilamentExtrudeButton(
		this.UI.window,
		this.UI.Client,
		this.amountToExtrudeStepButton,
		this.flowRateStepButton,
		this.selectToolStepButton,
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
		this.selectToolStepButton,
		false,
	)
	this.Grid().Attach(this.filamentRetractButton,		3, 0, 1, 1)


	// Second row
	this.filamentLoadButton = uiWidgets.CreateFilamentLoadButton(
		this.UI.window,
		this.UI.Client,
		this.flowRateStepButton,
		this.selectToolStepButton,
		true,
		int(this.UI.Settings.FilamentInLength),  // jab revisit
	)
	this.Grid().Attach(this.filamentLoadButton,			0, 1, 1, 1)

	this.temperatureStatusBox = uiWidgets.CreateTemperatureStatusBox(this.UI.Client, true, true)
	this.Grid().Attach(this.temperatureStatusBox,		1, 1, 2, 1)

	this.filamentUnloadButton = uiWidgets.CreateFilamentLoadButton(
		this.UI.window,
		this.UI.Client,
		this.flowRateStepButton,
		this.selectToolStepButton,
		false,
		int(this.UI.Settings.FilamentOutLength), // jab revisit
	)
	this.Grid().Attach(this.filamentUnloadButton,		3, 1, 1, 1)


	// Third row
	this.temperatureButton = utils.MustButtonImageStyle("Temperature", "heat-up.svg", "color1", this.showTemperaturePanel)
	this.Grid().Attach(this.temperatureButton, 0, 2, 1, 1)

	// The select tool step button is needed by some of the other controls (to get the name/ID of the tool
	// to send the command to), but only display it if multiple hotends are present.
	toolheadCount := utils.GetToolheadCount(this.UI.Client)
	if toolheadCount > 1 {
		this.Grid().Attach(this.selectToolStepButton, 1, 2, 1, 1)
	}
}

func (this *filamentPanel) showTemperaturePanel() {
	this.UI.GoToPanel(TemperaturePanel(this.UI, this))
}
