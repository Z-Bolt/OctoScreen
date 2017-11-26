package gui

import "github.com/gotk3/gotk3/gtk"

type HomeMenu struct {
	*gtk.Grid
	gui *GUI
}

func NewHomeMenu(gui *GUI) *HomeMenu {
	grid, _ := gtk.GridNew()

	m := &HomeMenu{Grid: grid, gui: gui}
	m.initialize()
	return m
}

func (m *HomeMenu) initialize() {
	m.Attach(NewButtonImage("Status", "status.svg", nil), 1, 0, 1, 1)
	m.Attach(NewButtonImage("Heat Up", "heat-up.svg", nil), 2, 0, 1, 1)
	m.Attach(NewButtonImage("Move", "move.svg", m.ShowMove), 3, 0, 1, 1)
	m.Attach(NewButtonImage("Home", "home.svg", nil), 4, 0, 1, 1)
	m.Attach(NewButtonImage("Extruct", "extruct.svg", nil), 1, 1, 1, 1)
	m.Attach(NewButtonImage("HeatBed", "bed.svg", nil), 2, 1, 1, 1)
	m.Attach(NewButtonImage("Fan", "fan.svg", nil), 3, 1, 1, 1)
	m.Attach(NewButtonImage("Settings", "settings.svg", nil), 4, 1, 1, 1)
}

func (m *HomeMenu) ShowMove() {
	m.gui.Add(NewMoveMenu(m.gui).Grid)
}
