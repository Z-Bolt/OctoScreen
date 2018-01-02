package ui

var defaultPanel *DefaultPanel

type DefaultPanel struct {
	CommonPanel
}

func NewDefaultPanel(ui *UI) Panel {
	if defaultPanel == nil {
		m := &DefaultPanel{CommonPanel: NewCommonPanel(ui, nil)}
		m.initialize()
		defaultPanel = m
	}

	return defaultPanel
}

func (m *DefaultPanel) initialize() {
	m.Grid().Attach(MustButtonImage("Status", "status.svg", m.showStatus), 1, 0, 1, 1)
	m.Grid().Attach(MustButtonImage("Heat Up", "heat-up.svg", m.showTemperature), 2, 0, 1, 1)
	m.Grid().Attach(MustButtonImage("Move", "move.svg", m.showMove), 3, 0, 1, 1)
	m.Grid().Attach(MustButtonImage("Home", "home.svg", m.showHome), 4, 0, 1, 1)
	m.Grid().Attach(MustButtonImage("Filament", "filament.svg", m.showFilament), 1, 1, 1, 1)
	m.Grid().Attach(MustButtonImage("Control", "control.svg", m.showControl), 2, 1, 1, 1)
	m.Grid().Attach(MustButtonImage("Files", "files.svg", m.showFiles), 3, 1, 1, 1)
	m.Grid().Attach(MustButtonImage("System", "settings.svg", m.showSystem), 4, 1, 1, 1)
}

func (m *DefaultPanel) showStatus() {
	m.UI.Add(NewStatusPanel(m.UI, m))
}

func (m *DefaultPanel) showHome() {
	m.UI.Add(NewHomePanel(m.UI, m))
}

func (m *DefaultPanel) showTemperature() {
	m.UI.Add(NewTemperaturePanel(m.UI, m))
}

func (m *DefaultPanel) showFilament() {
	m.UI.Add(NewFilamentPanel(m.UI, m))
}

func (m *DefaultPanel) showMove() {
	m.UI.Add(NewMovePanel(m.UI, m))
}

func (m *DefaultPanel) showControl() {
	m.UI.Add(NewControlPanel(m.UI, m))
}

func (m *DefaultPanel) showFiles() {
	m.UI.Add(NewFilesPanel(m.UI, m))
}

func (m *DefaultPanel) showSystem() {
	m.UI.Add(NewSystemPanel(m.UI, m))
}
