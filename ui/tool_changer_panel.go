package ui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"

	// "github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)


type toolChangerPanel struct {
	CommonPanel
	//activeTool int
}

var toolChangerPanelInstance *toolChangerPanel

func GetToolChangerPanelInstance(
	ui				*UI,
) *toolChangerPanel {
	if toolChangerPanelInstance == nil {
		this := &toolChangerPanel {
			CommonPanel: CreateCommonPanel("ToolChangerPanel", ui),
		}
		this.initialize()
		toolChangerPanelInstance = this
	}

	return toolChangerPanelInstance
}

func (this *toolChangerPanel) initialize() {
	defer this.Initialize()

	this.createToolheadButtons()

	homeAllButton := uiWidgets.CreateHomeAllButton(this.UI.Client)
	this.Grid().Attach(homeAllButton,                   0, 1, 1, 1)

	this.Grid().Attach(this.createMagnetOnButton(),     2, 1, 1, 1)
	this.Grid().Attach(this.createMagnetOffButton(),    3, 1, 1, 1)
	this.Grid().Attach(this.createZCalibrationButton(), 1, 2, 1, 1)
}

func (this *toolChangerPanel) createZCalibrationButton() gtk.IWidget {
	button := utils.MustButtonImageStyle("Z Offsets", "z-calibration.svg", "color2", func() {
		this.UI.GoToPanel(GetZOffsetCalibrationPanelInstance(this.UI))
	})

	return button
}

func (this *toolChangerPanel) createMagnetOnButton() gtk.IWidget {
	return utils.MustButtonImageStyle("Magnet On", "magnet-on.svg", "", func() {
		cmd := &octoprintApis.CommandRequest{}
		cmd.Commands = []string{"SET_PIN PIN=sol VALUE=1"}

		logger.Info("Turn on magnet")
		if err := cmd.Do(this.UI.Client); err != nil {
			logger.LogError("tool-changer.createMagnetOnButton()", "Do(CommandRequest)", err)
			return
		}
	})
}

func (this *toolChangerPanel) createMagnetOffButton() gtk.IWidget {
	return utils.MustButtonImageStyle("Magnet Off", "magnet-off.svg", "", func() {
		cmd := &octoprintApis.CommandRequest{}
		cmd.Commands = []string{"SET_PIN PIN=sol VALUE=0"}

		logger.Info("Turn off magnet")
		if err := cmd.Do(this.UI.Client); err != nil {
			logger.LogError("tool-changer.createMagnetOffButton()", "Do(CommandRequest)", err)
			return
		}
	})
}

func (this *toolChangerPanel) createToolheadButtons() {
	extruderCount := utils.GetExtruderCount(this.UI.Client)
	toolheadButtons := utils.CreateChangeToolheadButtonsAndAttachToGrid(extruderCount, this.Grid())
	this.setToolheadButtonClickHandlers(toolheadButtons)
}

func (this *toolChangerPanel) setToolheadButtonClickHandlers(toolheadButtons []*gtk.Button) {
	for index, toolheadButton := range toolheadButtons {
		this.setToolheadButtonClickHandler(toolheadButton, index)
	}
}

func (this *toolChangerPanel) setToolheadButtonClickHandler(toolheadButton *gtk.Button, toolheadIndex int) {
	toolheadButton.Connect("clicked", func() {
		logger.Infof("Changing tool to tool%d", toolheadIndex)

		gcode := fmt.Sprintf("T%d", toolheadIndex)
		this.command(gcode)
	})
}
