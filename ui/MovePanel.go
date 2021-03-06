package ui

import (
	// "github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	// "github.com/Z-Bolt/OctoScreen/utils"
)

var movePanelInstance *movePanel

type movePanel struct {
	CommonPanel
	amountToMoveStepButton    *uiWidgets.AmountToMoveStepButton
}

func MovePanel(
	ui				*UI,
	parentPanel		interfaces.IPanel,
) *movePanel {
	if movePanelInstance == nil {
		instance := &movePanel {
			CommonPanel: NewCommonPanel(ui, parentPanel),
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
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "X-", "move-x-.svg", octoprint.XAxis,  1), 0, 1, 1, 1)
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "X+", "move-x+.svg", octoprint.XAxis, -1), 2, 1, 1, 1)
	} else {
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "X-", "move-x-.svg", octoprint.XAxis, -1), 0, 1, 1, 1)
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "X+", "move-x+.svg", octoprint.XAxis,  1), 2, 1, 1, 1)
	}

	if yAxisInverted {
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "Y+", "move-y+.svg", octoprint.YAxis, -1), 1, 0, 1, 1)
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "Y-", "move-y-.svg", octoprint.YAxis,  1), 1, 2, 1, 1)
	} else {
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "Y+", "move-y+.svg", octoprint.YAxis,  1), 1, 0, 1, 1)
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "Y-", "move-y-.svg", octoprint.YAxis, -1), 1, 2, 1, 1)
	}

	if zAxisInverted {
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "Z+", "move-z+.svg", octoprint.ZAxis, -1), 3, 0, 1, 1)
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "Z-", "move-z-.svg", octoprint.ZAxis,  1), 3, 1, 1, 1)
	} else {
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "Z+", "move-z+.svg", octoprint.ZAxis,  1), 3, 0, 1, 1)
		this.Grid().Attach(uiWidgets.CreateMoveButton(this.UI.Client, this.amountToMoveStepButton, "Z-", "move-z-.svg", octoprint.ZAxis, -1), 3, 1, 1, 1)
	}

	this.Grid().Attach(this.amountToMoveStepButton, 1, 1, 1, 1)
}
