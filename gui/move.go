package gui

import (
	"github.com/mcuadros/OctoPrint-TFT/octoprint"

	"github.com/gotk3/gotk3/gtk"
)

type step struct {
	L string
	V int
}

type MoveMenu struct {
	gui   *GUI
	step  int
	steps []step

	*gtk.Grid
}

func NewMoveMenu(gui *GUI) *MoveMenu {
	m := &MoveMenu{Grid: MustGrid(),
		gui: gui,
		steps: []step{
			{"5mm", 5},
			{"10mm", 10},
			{"1mm", 1},
		},
	}

	m.initialize()
	return m
}

func (m *MoveMenu) initialize() {
	m.Attach(m.createMoveButton("X+", "move-x+.svg", octoprint.XAxis, 1), 1, 0, 1, 1)
	m.Attach(m.createMoveButton("X-", "move-x-.svg", octoprint.XAxis, -1), 1, 1, 1, 1)
	m.Attach(m.createMoveButton("Y+", "move-y+.svg", octoprint.YAxis, 1), 2, 0, 1, 1)
	m.Attach(m.createMoveButton("Y-", "move-y-.svg", octoprint.YAxis, -1), 2, 1, 1, 1)
	m.Attach(m.createMoveButton("Z+", "move-z+.svg", octoprint.ZAxis, 1), 3, 0, 1, 1)
	m.Attach(m.createMoveButton("Z-", "move-z-.svg", octoprint.ZAxis, -1), 3, 1, 1, 1)
	m.Attach(m.createStepButton(), 4, 0, 1, 1)
	m.Attach(MustButtonImage("Back", "back.svg", m.gui.ShowMenu), 4, 1, 1, 1)
}

func (m *MoveMenu) createStepButton() gtk.IWidget {
	b := MustButtonImage(m.steps[m.step].L, "move-step.svg", nil)
	b.Connect("clicked", func() {
		m.step++
		if m.step >= len(m.steps) {
			m.step = 0
		}

		b.SetLabel(m.steps[m.step].L)
	})

	return b
}

func (m *MoveMenu) createMoveButton(label, image string, a octoprint.Axis, dir int) gtk.IWidget {
	return MustButtonImage(label, image, func() {
		distance := m.steps[m.step].V * dir

		cmd := &octoprint.JogCommand{}
		switch a {
		case octoprint.XAxis:
			cmd.X = distance
		case octoprint.YAxis:
			cmd.Y = distance
		case octoprint.ZAxis:
			cmd.Z = distance
		}

		if err := cmd.Do(m.gui.Printer); err != nil {
			panic(err)
		}
	})
}
