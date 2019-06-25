package ui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

var toolchangerPanelInstance *toolchangerPanel

type toolchangerPanel struct {
	CommonPanel
}

func ToolchangerPanel(ui *UI, parent Panel) Panel {
	if toolchangerPanelInstance == nil {
		m := &toolchangerPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.panelH = 3
		// m.b = NewBackgroundTask(time.Second, m.updateTemperatures)
		m.initialize()

		toolchangerPanelInstance = m
	}

	return toolchangerPanelInstance
}

func (m *toolchangerPanel) initialize() {
	defer m.Initialize()
	m.Grid().Attach(m.createChangeToolButton(0), 1, 0, 1, 1)
	m.Grid().Attach(m.createChangeToolButton(1), 2, 0, 1, 1)
	m.Grid().Attach(m.createChangeToolButton(2), 3, 0, 1, 1)
	m.Grid().Attach(m.createChangeToolButton(3), 4, 0, 1, 1)
	m.Grid().Attach(m.createMagnetOnButton(), 1, 2, 1, 1)
	m.Grid().Attach(m.createMagnetOffButton(), 2, 2, 1, 1)
}

func (m *toolchangerPanel) createChangeToolButton(num int) gtk.IWidget {
	style := fmt.Sprintf("color%d", num+1)
	name := fmt.Sprintf("Tool%d", num+1)
	gcode := fmt.Sprintf("T%d", num)
	return MustButtonImageStyle(name, "extruder.svg", style, func() {
		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{gcode}

		Logger.Info("Switching tool")
		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})
}

func (m *toolchangerPanel) createMagnetOnButton() gtk.IWidget {
	return MustButtonImageStyle("Magnet On", "magnet-on.svg", "color4", func() {
		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{"SET_PIN PIN=sol VALUE=1"}

		Logger.Info("Turn on magnet")
		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})
}

func (m *toolchangerPanel) createMagnetOffButton() gtk.IWidget {
	return MustButtonImageStyle("Magnet Off", "magnet-off.svg", "color3", func() {
		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{"SET_PIN PIN=sol VALUE=0"}

		Logger.Info("Turn off magnet")
		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})
}
