package ui

import (
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

var movePanelInstance *movePanel

type movePanel struct {
	CommonPanel
	step *StepButton
}

func MovePanel(ui *UI, parent Panel) Panel {
	if movePanelInstance == nil {
		m := &movePanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.initialize()
		movePanelInstance = m
	}

	return movePanelInstance
}

func (m *movePanel) initialize() {
	defer m.Initialize()

	m.AddButton(m.createMoveButton("X+", "move-x+.svg", octoprint.XAxis, 1))
	m.AddButton(m.createMoveButton("Y-", "move-y+.svg", octoprint.YAxis, -1))
	m.AddButton(m.createMoveButton("Z+", "move-z+.svg", octoprint.ZAxis, 1))

	m.step = MustStepButton("move-step.svg",
		Step{"20mm", 20}, Step{"10mm", 10}, Step{"1mm", 1},
	)

	m.AddButton(m.step)
	m.AddButton(m.createMoveButton("X-", "move-x-.svg", octoprint.XAxis, -1))
	m.AddButton(m.createMoveButton("Y+", "move-y-.svg", octoprint.YAxis, 1))
	m.AddButton(m.createMoveButton("Z-", "move-z-.svg", octoprint.ZAxis, -1))
	
	m.Grid().Attach(MustButtonImage("Back", "back.svg", m.UI.GoHistory), 4, 1, 1, 1)
}

func (m *movePanel) createMoveButton(label, image string, a octoprint.Axis, dir int) gtk.IWidget {
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

		Logger.Warningf("Jogging %s axis in %dmm",
			strings.ToUpper(string(a)), distance,
		)

		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})
}
