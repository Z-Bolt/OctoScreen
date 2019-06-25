package ui

var defaultPanelInstance *defaultPanel

type defaultPanel struct {
	CommonPanel
}

func DefaultPanel(ui *UI) Panel {
	if defaultPanelInstance == nil {
		m := &defaultPanel{CommonPanel: NewCommonPanel(ui, nil)}
		m.initialize()
		defaultPanelInstance = m
	}

	return defaultPanelInstance
}

func (m *defaultPanel) initialize() {
	m.Grid().Attach(MustButtonImageStyle("Status", "status.svg", "color1", m.showStatus), 1, 0, 1, 1)
	m.Grid().Attach(MustButtonImageStyle("Heat Up", "heat-up.svg", "color2", m.showTemperature), 2, 0, 1, 1)
	m.Grid().Attach(MustButtonImageStyle("Move", "move.svg", "color3", m.showMove), 3, 0, 1, 1)
	m.Grid().Attach(MustButtonImageStyle("Home", "home.svg", "color4", m.showHome), 4, 0, 1, 1)

	m.Grid().Attach(MustButtonImageStyle("Print", "print.svg", "color2", m.showFiles), 1, 1, 1, 1)
	m.Grid().Attach(MustButtonImageStyle("Filament", "filament.svg", "color1", m.showFilament), 2, 1, 1, 1)
	m.Grid().Attach(MustButtonImageStyle("Control", "control.svg", "color4", m.showControl), 3, 1, 1, 1)
	m.Grid().Attach(MustButtonImageStyle("System", "settings.svg", "color3", m.showSystem), 4, 1, 1, 1)
}

func (m *defaultPanel) showStatus() {
	m.UI.Add(StatusPanel(m.UI, m))
}

func (m *defaultPanel) showHome() {
	m.UI.Add(HomePanel(m.UI, m))
}

func (m *defaultPanel) showTemperature() {
	m.UI.Add(TemperaturePanel(m.UI, m))
}

func (m *defaultPanel) showFilament() {
	m.UI.Add(FilamentPanel(m.UI, m))
}

func (m *defaultPanel) showMove() {
	m.UI.Add(MovePanel(m.UI, m))
}

func (m *defaultPanel) showControl() {
	m.UI.Add(ControlPanel(m.UI, m))
}

func (m *defaultPanel) showFiles() {
	m.UI.Add(FilesPanel(m.UI, m))
}

func (m *defaultPanel) showSystem() {
	m.UI.Add(SystemPanel(m.UI, m))
}
