package ui

import "github.com/gotk3/gotk3/gtk"

var idleConfigurationMenuPanelInstance *idleConfigurationMenuPanel

type idleConfigurationMenuPanel struct {
	CommonPanel
}

func IdleConfigurationMenuPanel(ui *UI, parent Panel) Panel {
	if idleConfigurationMenuPanelInstance == nil {
		m := &idleConfigurationMenuPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.initialize()
		idleConfigurationMenuPanelInstance = m
	}

	return idleConfigurationMenuPanelInstance
}

func (m *idleConfigurationMenuPanel) initialize() {
	defer m.Initialize()

	var buttons = []gtk.IWidget{
		MustButtonImageStyle("Bed Level", "bed-level.svg", "color4", m.showCalibrate),
		MustButtonImageStyle("ZOffsets", "z-offset-increase.svg", "color2", m.showNozzleCalibration),
		MustButtonImageStyle("Network", "network.svg", "color1", m.showNetwork),
		MustButtonImageStyle("System", "info.svg", "color3", m.showSystem),
	}

	m.arrangeButtons(buttons)
}

func (m *idleConfigurationMenuPanel) showNetwork() {
	m.UI.Add(NetworkPanel(m.UI, m))
}

func (m *idleConfigurationMenuPanel) showSystem() {
	m.UI.Add(SystemPanel(m.UI, m))
}

func (m *idleConfigurationMenuPanel) showCalibrate() {
	m.UI.Add(BedLevelPanel(m.UI, m))
}

func (m *idleConfigurationMenuPanel) showNozzleCalibration() {
	m.UI.Add(NozzleCalibrationPanel(m.UI, m))
}
