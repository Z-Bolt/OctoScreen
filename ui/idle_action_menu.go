package ui

import "github.com/gotk3/gotk3/gtk"

var idleActionMenuPanelInstance *idleActionMenuPanel

type idleActionMenuPanel struct {
	CommonPanel
}

func IdleActionMenuPanel(ui *UI, parent Panel) Panel {
	if idleActionMenuPanelInstance == nil {
		m := &idleActionMenuPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.initialize()
		idleActionMenuPanelInstance = m
	}

	return idleActionMenuPanelInstance
}

func (m *idleActionMenuPanel) initialize() {
	defer m.Initialize()

	var buttons = []gtk.IWidget{
		MustButtonImageStyle("Move", "move.svg", "color1", m.showMove),
		MustButtonImageStyle("Extrude", "filament.svg", "color1", m.showExtrude),
		MustButtonImageStyle("Fan", "fan.svg", "color2", m.showFan),
		MustButtonImageStyle("Temperature", "heat-up.svg", "color4", m.showTemperature),
		MustButtonImageStyle("Control", "control.svg", "color4", m.showControl),
	}

	if m.UI.Settings != nil && m.UI.Settings.ToolChanger {
		buttons = append(buttons, MustButtonImageStyle("ToolChanger", "toolchanger.svg", "color2", m.showToolchanger))
	}

	m.arrangeButtons(buttons)
}

func (m *idleActionMenuPanel) showTemperature() {
	m.UI.Add(TemperaturePanel(m.UI, m))
}

func (m *idleActionMenuPanel) showExtrude() {
	m.UI.Add(ExtrudePanel(m.UI, m))
}

func (m *idleActionMenuPanel) showControl() {
	m.UI.Add(ControlPanel(m.UI, m))
}

func (m *idleActionMenuPanel) showToolchanger() {
	m.UI.Add(ToolchangerPanel(m.UI, m))
}

func (m *idleActionMenuPanel) showMove() {
	m.UI.Add(MovePanel(m.UI, m))
}

func (m *idleActionMenuPanel) showFan() {
	m.UI.Add(FanPanel(m.UI, m))
}
