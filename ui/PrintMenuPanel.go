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
	m.Grid().Attach(MustButtonImageStyle("Move",        "move.svg",           "color1", m.showMove),        0, 0, 1, 1)
	m.Grid().Attach(MustButtonImageStyle("Filament",    "filament-spool.svg", "color2", m.showFilament),    1, 0, 1, 1)
	m.Grid().Attach(MustButtonImageStyle("Extruder",    "extruder.svg",       "color3", m.showExtruder),    2, 0, 1, 1)
	m.Grid().Attach(MustButtonImageStyle("Temperature", "heat-up.svg",        "color4", m.showTemperature), 3, 0, 1, 1)

	m.Grid().Attach(MustButtonImageStyle("Fan",         "fan.svg",            "color1", m.showFan),         0, 1, 1, 1)
	m.Grid().Attach(MustButtonImageStyle("Network",     "network.svg",        "color2", m.showNetwork),     1, 1, 1, 1)
	m.Grid().Attach(MustButtonImageStyle("System",      "info.svg",           "color3", m.showSystem),      2, 1, 1, 1)
}

func (m *printMenuPanel) showMove() {
	m.UI.Add(MovePanel(m.UI, m))
}

func (m *printMenuPanel) showFilament() {
	m.UI.Add(FilamentPanel(m.UI, m))
}

func (m *printMenuPanel) showExtruder() {
	m.UI.Add(ExtruderPanel(m.UI, m))
}

func (m *printMenuPanel) showTemperature() {
	m.UI.Add(TemperaturePanel(m.UI, m))
}

func (m *printMenuPanel) showFan() {
	m.UI.Add(FanPanel(m.UI, m))
}

func (m *printMenuPanel) showNetwork() {
	m.UI.Add(NetworkPanel(m.UI, m))
}

func (m *printMenuPanel) showSystem() {
	m.UI.Add(SystemPanel(m.UI, m))
}
