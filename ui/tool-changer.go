package ui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/utils"
)

var toolchangerPanelInstance *toolchangerPanel

type toolchangerPanel struct {
	CommonPanel
	//activeTool int
}

func ToolchangerPanel(ui *UI, parent Panel) Panel {
	if toolchangerPanelInstance == nil {
		m := &toolchangerPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.initialize()
		toolchangerPanelInstance = m
	}

	return toolchangerPanelInstance
}

func (m *toolchangerPanel) initialize() {
	defer m.Initialize()

	m.createToolheadButtons()

	m.Grid().Attach(m.createHomeButton(),         0, 1, 1, 1)

	m.Grid().Attach(m.createMagnetOnButton(),     2, 1, 1, 1)
	m.Grid().Attach(m.createMagnetOffButton(),    3, 1, 1, 1)
	m.Grid().Attach(m.createZCalibrationButton(), 1, 2, 1, 1)
}

func (m *toolchangerPanel) createZCalibrationButton() gtk.IWidget {
	b := MustButtonImageStyle("Z Offsets", "z-calibration.svg", "color2", func() {
		m.UI.Add(ZOffsetCalibrationPanel(m.UI, m))
	})

	return b
}

func (m *toolchangerPanel) createHomeButton() gtk.IWidget {
	return MustButtonImageStyle("Home XYZ", "home.svg", "", func() {
		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{
			"G28 Z",
			"G28 X",
			"G28 Y",
		}

		if err := cmd.Do(m.UI.Printer); err != nil {
			utils.LogError("tool-changer.createHomeButton()", "Do(CommandRequest)", err)
		}
	})
}

func (m *toolchangerPanel) createMagnetOnButton() gtk.IWidget {
	return MustButtonImageStyle("Magnet On", "magnet-on.svg", "", func() {
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
	return MustButtonImageStyle("Magnet Off", "magnet-off.svg", "", func() {
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
	toolheadButtons := CreateChangeToolheadButtonsAndAttachToGrid(toolheadCount, m.Grid())
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
