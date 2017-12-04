package ui

import (
	"github.com/gotk3/gotk3/gtk"
)

type ToolsPanel struct {
	CommonPanel
}

func NewToolsPanel(ui *UI) *ToolsPanel {
	m := &ToolsPanel{CommonPanel: NewCommonPanel(ui)}
	m.initialize()
	return m
}

func (m *ToolsPanel) initialize() {
	m.grid.Attach(m.createMotorOff(), 1, 0, 1, 1)
	m.grid.Attach(m.createFanOn(), 2, 0, 1, 1)
	m.grid.Attach(m.createFanOff(), 3, 0, 1, 1)
	m.grid.Attach(m.createCalibrate(), 4, 0, 1, 1)
	m.grid.Attach(MustButtonImage("Back", "back.svg", m.UI.ShowDefaultPanel), 4, 1, 1, 1)
}

func (m *ToolsPanel) createMotorOff() gtk.IWidget {
	return MustButtonImage("Motor Off", "motor-off.svg", nil)
}

func (m *ToolsPanel) createFanOn() gtk.IWidget {
	return MustButtonImage("Fan On", "fan-on.svg", nil)
}

func (m *ToolsPanel) createFanOff() gtk.IWidget {
	return MustButtonImage("Fan Off", "fan.svg", nil)
}

func (m *ToolsPanel) createCalibrate() gtk.IWidget {
	b := MustButtonImage("Calibrate", "calibrate.svg", nil)
	b.SetSensitive(false)

	return b
}
