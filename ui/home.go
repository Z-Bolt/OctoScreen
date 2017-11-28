package ui

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

type HomePanel struct {
	CommonPanel
}

func NewHomePanel(ui *UI) *HomePanel {
	m := &HomePanel{CommonPanel: NewCommonPanel(ui)}
	m.initialize()
	return m
}

func (m *HomePanel) initialize() {
	m.grid.Attach(m.createMoveButton("Home All", "home.svg",
		octoprint.XAxis, octoprint.YAxis, octoprint.ZAxis,
	), 1, 0, 1, 1)

	m.grid.Attach(m.createMoveButton("Home X", "home-x.svg", octoprint.XAxis), 2, 0, 1, 1)
	m.grid.Attach(m.createMoveButton("Home Y", "home-y.svg", octoprint.YAxis), 3, 0, 1, 1)
	m.grid.Attach(m.createMoveButton("Home Z", "home-z.svg", octoprint.ZAxis), 4, 0, 1, 1)
	m.grid.Attach(MustButtonImage("Back", "back.svg", m.UI.ShowDefaultPanel), 4, 1, 1, 1)
}

func (m *HomePanel) createMoveButton(label, image string, axes ...octoprint.Axis) gtk.IWidget {
	return MustButtonImage(label, image, func() {
		cmd := &octoprint.PrintHeadHomeRequest{Axes: axes}
		if err := cmd.Do(m.UI.Printer); err != nil {
			panic(err)
		}
	})
}
