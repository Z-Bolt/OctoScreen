package ui

var printMenuPanelInstance *printMenuPanel

type printMenuPanel struct {
	CommonPanel
}

func PrintMenuPanel(ui *UI, parent Panel) Panel {
	if printMenuPanelInstance == nil {
		m := &printMenuPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.initialize()
		printMenuPanelInstance = m
	}

	return printMenuPanelInstance
}

func (m *printMenuPanel) initialize() {
	defer m.Initialize()
	m.Grid().Attach(MustButtonImageStyle("Temperature", "heat-up.svg", "color4", m.showTemperature), 1, 0, 1, 1)
	m.Grid().Attach(MustButtonImageStyle("Fan", "fan.svg", "color2", m.showFan), 2, 0, 1, 1)
	m.Grid().Attach(MustButtonImageStyle("Network", "network.svg", "color1", m.showNetwork), 3, 0, 1, 1)
	m.Grid().Attach(MustButtonImageStyle("System", "info.svg", "color3", m.showSystem), 4, 0, 1, 1)
	m.Grid().Attach(MustButtonImageStyle("Filament", "filament.svg", "color3", m.showFilament), 1, 1, 1, 1)
}

func (m *printMenuPanel) showTemperature() {
	m.UI.Add(TemperaturePanel(m.UI, m))
}

func (m *printMenuPanel) showNetwork() {
	m.UI.Add(NetworkPanel(m.UI, m))
}

func (m *printMenuPanel) showSystem() {
	m.UI.Add(SystemPanel(m.UI, m))
}

func (m *printMenuPanel) showFan() {
	m.UI.Add(FanPanel(m.UI, m))
}

func (m *printMenuPanel) showFilament() {
	m.UI.Add(FilamentPanel(m.UI, m))
}

// func (m *printMenuPanel) showControl() {
// 	m.UI.Add(ControlPanel(m.UI, m))
// }
