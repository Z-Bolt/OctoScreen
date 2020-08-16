package ui

import (
)

var configurationPanelInstance *configurationPanel

type configurationPanel struct {
	CommonPanel
}

func ConfigurationPanel(ui *UI, parent Panel) Panel {
	if configurationPanelInstance == nil {
		m := &configurationPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.initialize()
		configurationPanelInstance = m
	}

	return configurationPanelInstance
}

func (m *configurationPanel) initialize() {
	defer m.Initialize()

	bedlevelButton := MustButtonImageStyle("Bed Level", "bed-level.svg", "color1", m.handleBedLevelClicked)
	m.Grid().Attach(bedlevelButton, 0, 0, 1, 1)

	zOffsetCalibrationButton := MustButtonImageStyle("Z-Offset Calibration", "z-offset-increase.svg", "color2", m.handleZOffsetCalibrationClicked)
	m.Grid().Attach(zOffsetCalibrationButton, 1, 0, 1, 1)

	networkButton := MustButtonImageStyle("Network", "network.svg", "color3", m.handleNetworkClicked)
	m.Grid().Attach(networkButton, 2, 0, 1, 1)

	systemButton := MustButtonImageStyle("System", "info.svg", "color4", m.handleSystemClicked)
	m.Grid().Attach(systemButton, 3, 0, 1, 1)
}

func (m *configurationPanel) handleBedLevelClicked() {
	m.UI.Add(BedLevelPanel(m.UI, m))
}

func (m *configurationPanel) handleZOffsetCalibrationClicked() {
	m.UI.Add(ZOffsetCalibrationPanel(m.UI, m))
}

func (m *configurationPanel) handleNetworkClicked() {
	m.UI.Add(NetworkPanel(m.UI, m))
}

func (m *configurationPanel) handleSystemClicked() {
	m.UI.Add(SystemPanel(m.UI, m))
}
