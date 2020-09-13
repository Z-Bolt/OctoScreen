package ui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)

var toolchangerPanelInstance *toolchangerPanel

type toolchangerPanel struct {
	CommonPanel
	//activeTool int
}

func ToolchangerPanel(
	ui				*UI,
	parentPanel		interfaces.IPanel,
) *toolchangerPanel {
	if toolchangerPanelInstance == nil {
		m := &toolchangerPanel {
			CommonPanel: NewCommonPanel(ui, parentPanel),
		}
		m.initialize()
		toolchangerPanelInstance = m
	}

	return toolchangerPanelInstance
}

func (this *toolchangerPanel) initialize() {
	defer this.Initialize()

	this.createToolheadButtons()

	homeAllButton := uiWidgets.CreateHomeAllButton(this.UI.Printer)
	this.Grid().Attach(homeAllButton,                   0, 1, 1, 1)

	this.Grid().Attach(this.createMagnetOnButton(),     2, 1, 1, 1)
	this.Grid().Attach(this.createMagnetOffButton(),    3, 1, 1, 1)
	this.Grid().Attach(this.createZCalibrationButton(), 1, 2, 1, 1)
}

func (m *toolchangerPanel) createZCalibrationButton() gtk.IWidget {
	b := utils.MustButtonImageStyle("Z Offsets", "z-calibration.svg", "color2", func() {
		m.UI.GoToPanel(ZOffsetCalibrationPanel(m.UI, m))
	})

	return b
}

func (m *toolchangerPanel) createMagnetOnButton() gtk.IWidget {
	return utils.MustButtonImageStyle("Magnet On", "magnet-on.svg", "", func() {
		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{"SET_PIN PIN=sol VALUE=1"}

		utils.Logger.Info("Turn on magnet")
		if err := cmd.Do(m.UI.Printer); err != nil {
			utils.LogError("tool-changer.createMagnetOnButton()", "Do(CommandRequest)", err)
			return
		}
	})
}

func (m *toolchangerPanel) createMagnetOffButton() gtk.IWidget {
	return utils.MustButtonImageStyle("Magnet Off", "magnet-off.svg", "", func() {
		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{"SET_PIN PIN=sol VALUE=0"}

		utils.Logger.Info("Turn off magnet")
		if err := cmd.Do(m.UI.Printer); err != nil {
			utils.LogError("tool-changer.createMagnetOffButton()", "Do(CommandRequest)", err)
			return
		}
	})
}

func (m *toolchangerPanel) createToolheadButtons() {
	toolheadCount := utils.GetToolheadCount(m.UI.Printer)
	toolheadButtons := utils.CreateChangeToolheadButtonsAndAttachToGrid(toolheadCount, m.Grid())
	m.setToolheadButtonClickHandlers(toolheadButtons)
}

func (m *toolchangerPanel) setToolheadButtonClickHandlers(toolheadButtons []*gtk.Button) {
	for index, toolheadButton := range toolheadButtons {
		m.setToolheadButtonClickHandler(toolheadButton, index)
	}
}

func (m *toolchangerPanel) setToolheadButtonClickHandler(toolheadButton *gtk.Button, toolheadIndex int) {
	toolheadButton.Connect("clicked", func() {
		utils.Logger.Infof("Changing tool to tool%d", toolheadIndex)

		gcode := fmt.Sprintf("T%d", toolheadIndex)
		m.command(gcode)
	})
}
