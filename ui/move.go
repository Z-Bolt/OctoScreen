package ui

import (
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
		m.panelH = 3
		m.initialize()
		movePanelInstance = m
	}

	return movePanelInstance
}

func (m *movePanel) initialize() {
	defer m.Initialize()
	m.Grid().Attach(m.createMoveButton("X-", "move-x-.svg", octoprint.XAxis, -1), 1, 1, 1, 1)
	m.Grid().Attach(m.createMoveButton("X+", "move-x+.svg", octoprint.XAxis, 1), 3, 1, 1, 1)
	m.Grid().Attach(m.createMoveButton("Y+", "move-y+.svg", octoprint.YAxis, 1), 2, 0, 1, 1)
	m.Grid().Attach(m.createMoveButton("Y-", "move-y-.svg", octoprint.YAxis, -1), 2, 2, 1, 1)

	if m.UI.Settings != nil && m.UI.Settings.ZAxisInverted {
		m.Grid().Attach(m.createMoveButton("Z-", "move-z-.svg", octoprint.ZAxis, -1), 4, 0, 1, 1)
		m.Grid().Attach(m.createMoveButton("Z+", "move-z+.svg", octoprint.ZAxis, 1), 4, 1, 1, 1)
	} else {
		m.Grid().Attach(m.createMoveButton("Z+", "move-z-.svg", octoprint.ZAxis, 1), 4, 0, 1, 1)
		m.Grid().Attach(m.createMoveButton("Z-", "move-z+.svg", octoprint.ZAxis, -1), 4, 1, 1, 1)
	}

	m.step = MustStepButton("move-step.svg",
		Step{"10mm", 10.0}, Step{"1mm", 1.0}, Step{"0.1mm", 0.1}, Step{"0.02mm", 0.02},
	)

	m.Grid().Attach(m.step, 3, 2, 1, 1)

	m.Grid().Attach(m.createHomeButton(), 1, 2, 1, 1)
}

func (m *movePanel) createMoveButton(label, image string, a octoprint.Axis, dir float64) gtk.IWidget {

	return MustPressedButton(label, image, func() {
		distance := m.step.Value().(float64) * dir

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
			Logger.Error(err)
			return
		}

	}, 200)
}

func (m *movePanel) createHomeButton() gtk.IWidget {
	return MustButtonImage("Home All", "home.svg", func() {
		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{
			"G28",
		}

		Logger.Info("Sending filament unload request")
		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})

}
