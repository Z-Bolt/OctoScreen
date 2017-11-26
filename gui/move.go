package gui

import "github.com/gotk3/gotk3/gtk"

type MoveMenu struct {
	*gtk.Grid
	gui *GUI
}

func NewMoveMenu(gui *GUI) *MoveMenu {
	grid, _ := gtk.GridNew()

	m := &MoveMenu{Grid: grid, gui: gui}
	m.initialize()
	return m
}

func (m *MoveMenu) initialize() {
	m.Attach(NewButtonImage("Status", "status.svg", nil), 1, 0, 1, 1)
}
