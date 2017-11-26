package gui

import "github.com/gotk3/gotk3/gtk"

type HomeMenu struct {
	*gtk.Grid
	gui *GUI
}

func NewHomeMenu(gui *GUI) *HomeMenu {
	m := &HomeMenu{Grid: MustGrid(), gui: gui}
	m.initialize()
	return m
}

func (m *HomeMenu) initialize() {
	m.Attach(MustButtonImage("Status", "status.svg", nil), 1, 0, 1, 1)
	m.Attach(MustButtonImage("Heat Up", "heat-up.svg", nil), 2, 0, 1, 1)
	m.Attach(MustButtonImage("Move", "move.svg", m.ShowMove), 3, 0, 1, 1)
	m.Attach(MustButtonImage("Home", "home.svg", nil), 4, 0, 1, 1)
	m.Attach(MustButtonImage("Extruct", "extruct.svg", nil), 1, 1, 1, 1)
	m.Attach(MustButtonImage("HeatBed", "bed.svg", nil), 2, 1, 1, 1)
	m.Attach(MustButtonImage("Fan", "fan.svg", nil), 3, 1, 1, 1)
	m.Attach(MustButtonImage("Settings", "settings.svg", nil), 4, 1, 1, 1)
}

func (m *HomeMenu) ShowMove() {
	m.gui.Add(NewMoveMenu(m.gui).Grid)
}
