package gui

import (
	"github.com/gotk3/gotk3/gtk"
)

type Menu struct {
	*gtk.Grid
	gui *GUI
}

func NewMenu(gui *GUI) *Menu {
	m := &Menu{Grid: MustGrid(), gui: gui}
	m.initialize()
	return m
}

func (m *Menu) initialize() {
	m.Attach(MustButtonImage("Status", "status.svg", nil), 1, 0, 1, 1)
	m.Attach(MustButtonImage("Heat Up", "heat-up.svg", nil), 2, 0, 1, 1)
	m.Attach(MustButtonImage("Move", "move.svg", m.ShowMove), 3, 0, 1, 1)
	m.Attach(MustButtonImage("Home", "home.svg", m.ShowHome), 4, 0, 1, 1)
	m.Attach(MustButtonImage("Extruct", "extruct.svg", nil), 1, 1, 1, 1)
	m.Attach(MustButtonImage("HeatBed", "bed.svg", nil), 2, 1, 1, 1)
	m.Attach(MustButtonImage("Fan", "fan.svg", nil), 3, 1, 1, 1)
	m.Attach(MustButtonImage("Settings", "settings.svg", nil), 4, 1, 1, 1)
}

func (m *Menu) Back() {
	m.gui.Add(NewHomeMenu(m.gui).Grid)
}

func (m *Menu) ShowHome() {
	m.gui.Add(NewHomeMenu(m.gui).Grid)
}

func (m *Menu) ShowMove() {
	m.gui.Add(NewMoveMenu(m.gui).Grid)
}
