package ui

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

type MovePanel struct {
	CommonPanel
	step *StepButton
}

func NewMovePanel(ui *UI) Panel {
	m := &MovePanel{CommonPanel: NewCommonPanel(ui)}
	m.initialize()
	return m
}

func (m *MovePanel) initialize() {
	m.grid.Attach(m.createMoveButton("X+", "move-x+.svg", octoprint.XAxis, 1), 1, 0, 1, 1)
	m.grid.Attach(m.createMoveButton("X-", "move-x-.svg", octoprint.XAxis, -1), 1, 1, 1, 1)
	m.grid.Attach(m.createMoveButton("Y+", "move-y+.svg", octoprint.YAxis, 1), 2, 0, 1, 1)
	m.grid.Attach(m.createMoveButton("Y-", "move-y-.svg", octoprint.YAxis, -1), 2, 1, 1, 1)
	m.grid.Attach(m.createMoveButton("Z+", "move-z+.svg", octoprint.ZAxis, 1), 3, 0, 1, 1)
	m.grid.Attach(m.createMoveButton("Z-", "move-z-.svg", octoprint.ZAxis, -1), 3, 1, 1, 1)

	m.step = MustStepButton("move-step.svg",
		Step{"5mm", 5}, Step{"10mm", 10}, Step{"1mm", 1},
	)
	m.grid.Attach(m.step, 4, 0, 1, 1)

	m.grid.Attach(MustButtonImage("Back", "back.svg", m.UI.ShowDefaultPanel), 4, 1, 1, 1)
}

func (m *MovePanel) createMoveButton(label, image string, a octoprint.Axis, dir int) gtk.IWidget {
	return MustButtonImage(label, image, func() {
		distance := m.step.Value().(int) * dir

		cmd := &octoprint.PrintHeadJogRequest{}
		switch a {
		case octoprint.XAxis:
			cmd.X = distance
		case octoprint.YAxis:
			cmd.Y = distance
		case octoprint.ZAxis:
			cmd.Z = distance
		}

		if err := cmd.Do(m.UI.Printer); err != nil {
			panic(err)
		}
	})
}
