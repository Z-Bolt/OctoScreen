package ui

type DefaultPanel struct {
	CommonPanel
}

func NewDefaultPanel(ui *UI) Panel {
	m := &DefaultPanel{CommonPanel: NewCommonPanel(ui)}
	m.initialize()
	return m
}

func (m *DefaultPanel) initialize() {
	m.grid.Attach(MustButtonImage("Status", "status.svg", nil), 1, 0, 1, 1)
	m.grid.Attach(MustButtonImage("Heat Up", "heat-up.svg", m.showTemperature), 2, 0, 1, 1)
	m.grid.Attach(MustButtonImage("Move", "move.svg", m.showMove), 3, 0, 1, 1)
	m.grid.Attach(MustButtonImage("Home", "home.svg", m.showHome), 4, 0, 1, 1)
	m.grid.Attach(MustButtonImage("Filament", "filament.svg", m.showFilament), 1, 1, 1, 1)
	m.grid.Attach(MustButtonImage("HeatBed", "bed.svg", nil), 2, 1, 1, 1)
	m.grid.Attach(MustButtonImage("Fan", "fan.svg", nil), 3, 1, 1, 1)
	m.grid.Attach(MustButtonImage("Settings", "settings.svg", nil), 4, 1, 1, 1)
}

func (m *DefaultPanel) showHome() {
	m.UI.Add(NewHomePanel(m.UI))
}

func (m *DefaultPanel) showTemperature() {
	m.UI.Add(NewTemperaturePanel(m.UI))
}

func (m *DefaultPanel) showFilament() {
	m.UI.Add(NewFilamentPanel(m.UI))
}

func (m *DefaultPanel) showMove() {
	m.UI.Add(NewMovePanel(m.UI))
}
