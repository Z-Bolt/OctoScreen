package ui

import (
	// "github.com/gotk3/gotk3/gtk"

	// "github.com/Z-Bolt/OctoScreen/interfaces"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	// "github.com/Z-Bolt/OctoScreen/utils"
)

type movePanel struct {
	CommonPanel
	amountToMoveStepButton    *uiWidgets.AmountToMoveStepButton
}

var movePanelInstance *movePanel

func GetMovePanelInstance(
	ui				*UI,
) *movePanel {
	if movePanelInstance == nil {
		instance := &movePanel {
			CommonPanel: CreateCommonPanel("MovePanel", ui),
		}
		instance.initialize()
		movePanelInstance = instance
	}

	return movePanelInstance
}

func (this *movePanel) initialize() {
	defer this.Initialize()

	// Create the step button first, since it is needed by some of the other controls.
	this.amountToMoveStepButton = uiWidgets.CreateAmountToMoveStepButton()

	xAxisInverted, yAxisInverted, zAxisInverted := false, false, false
	if this.UI.Settings != nil {
		xAxisInverted = this.UI.Settings.XAxisInverted
		yAxisInverted = this.UI.Settings.YAxisInverted
		zAxisInverted = this.UI.Settings.ZAxisInverted
	}

	if xAxisInverted {
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "X-", "move-x-.svg", dataModels.XAxis,  1), 0, 1, 1, 1)
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "X+", "move-x+.svg", dataModels.XAxis, -1), 2, 1, 1, 1)
	} else {
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "X-", "move-x-.svg", dataModels.XAxis, -1), 0, 1, 1, 1)
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "X+", "move-x+.svg", dataModels.XAxis,  1), 2, 1, 1, 1)
	}

	if yAxisInverted {
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "Y+", "move-y+.svg", dataModels.YAxis, -1), 1, 0, 1, 1)
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "Y-", "move-y-.svg", dataModels.YAxis,  1), 1, 2, 1, 1)
	} else {
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "Y+", "move-y+.svg", dataModels.YAxis,  1), 1, 0, 1, 1)
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "Y-", "move-y-.svg", dataModels.YAxis, -1), 1, 2, 1, 1)
	}

	if zAxisInverted {
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "Z+", "move-z+.svg", dataModels.ZAxis, -1), 3, 0, 1, 1)
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "Z-", "move-z-.svg", dataModels.ZAxis,  1), 3, 1, 1, 1)
	} else {
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "Z+", "move-z+.svg", dataModels.ZAxis,  1), 3, 0, 1, 1)
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "Z-", "move-z-.svg", dataModels.ZAxis, -1), 3, 1, 1, 1)
	}

	this.Grid().Attach(this.amountToMoveStepButton, 1, 1, 1, 1)
}
