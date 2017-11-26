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
	m.Attach(NewButtonImage("X+", "move-x+.svg", nil), 1, 0, 1, 1)
	m.Attach(NewButtonImage("X-", "move-x-.svg", nil), 1, 1, 1, 1)
	m.Attach(NewButtonImage("Y+", "move-y+.svg", nil), 2, 0, 1, 1)
	m.Attach(NewButtonImage("Y-", "move-y-.svg", nil), 2, 1, 1, 1)
	m.Attach(NewButtonImage("Z+", "move-z+.svg", nil), 3, 0, 1, 1)
	m.Attach(NewButtonImage("Z-", "move-z-.svg", nil), 3, 1, 1, 1)
	m.Attach(NewButtonImage("10mm", "move-step.svg", nil), 4, 0, 1, 1)
	m.Attach(NewButtonImage("Back", "back.svg", m.Back), 4, 1, 1, 1)

}

func (m *MoveMenu) Back() {
	m.gui.Add(NewHomeMenu(m.gui).Grid)
}
