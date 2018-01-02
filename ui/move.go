package ui

import (
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

type MovePanel struct {
	CommonPanel
	step *StepButton
}

func NewMovePanel(ui *UI, parent Panel) Panel {
	m := &MovePanel{CommonPanel: NewCommonPanel(ui, parent)}
	m.initialize()
	return m
}

func (m *MovePanel) initialize() {
	defer m.Initialize()

	m.AddButton(m.createMoveButton("X+", "move-x+.svg", octoprint.XAxis, 1))
	m.AddButton(m.createMoveButton("Y+", "move-y+.svg", octoprint.YAxis, 1))
	m.AddButton(m.createMoveButton("Z+", "move-z+.svg", octoprint.ZAxis, 1))

	m.step = MustStepButton("move-step.svg",
		Step{"5mm", 5}, Step{"10mm", 10}, Step{"1mm", 1},
	)

	m.AddButton(m.step)
	m.AddButton(m.createMoveButton("X-", "move-x-.svg", octoprint.XAxis, -1))
	m.AddButton(m.createMoveButton("Y-", "move-y-.svg", octoprint.YAxis, -1))
	m.AddButton(m.createMoveButton("Z-", "move-z-.svg", octoprint.ZAxis, -1))
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

		Logger.Warningf("Jogging print head axis %s in %dmm",
			strings.ToUpper(string(a)), distance,
		)

		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})
}
