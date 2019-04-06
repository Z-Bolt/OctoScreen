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
	m.Grid().Attach(MustButtonImage("Status", "status.svg", m.showStatus), 1, 0, 1, 1)
	m.Grid().Attach(MustButtonImage("Files", "files.svg", m.showFiles), 2, 0, 1, 1)
	m.Grid().Attach(MustButtonImage("Heat Up", "heat-up.svg", m.showTemperature), 3, 0, 1, 1)
	m.Grid().Attach(MustButtonImage("Filament", "filament_clean.svg", m.showFilament), 4, 0, 1, 1)
	m.Grid().Attach(MustButtonImage("Home", "home.svg", m.showHome), 1, 1, 1, 1)
	m.Grid().Attach(MustButtonImage("Move", "move.svg", m.showMove), 2, 1, 1, 1)
	m.Grid().Attach(MustButtonImage("Control", "fan-on.svg", m.showControl), 3, 1, 1, 1)
	m.Grid().Attach(MustButtonImage("System", "settings.svg", m.showSystem), 4, 1, 1, 1)
}

func (m *defaultPanel) showStatus() {
	m.UI.Add(StatusPanel(m.UI, m))
}

func (m *defaultPanel) showHome() {
	m.UI.Add(HomePanel(m.UI, m))
}

func (m *defaultPanel) showTemperature() {
	m.UI.Add(ProfilesPanel(m.UI, m))
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
