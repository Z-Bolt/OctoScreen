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
	m.Initialize()

	m.AddButton(m.createMoveButton("Home All", "home.svg",
		octoprint.XAxis, octoprint.YAxis, octoprint.ZAxis,
	))

	m.AddButton(m.createMoveButton("Home X", "home-x.svg", octoprint.XAxis))
	m.AddButton(m.createMoveButton("Home Y", "home-y.svg", octoprint.YAxis))
	m.AddButton(m.createMoveButton("Home Z", "home-z.svg", octoprint.ZAxis))
}

func (m *HomePanel) createMoveButton(label, image string, axes ...octoprint.Axis) gtk.IWidget {
	return MustButtonImage(label, image, func() {
		cmd := &octoprint.PrintHeadHomeRequest{Axes: axes}
		Logger.Warningf("Homing the print head in %s axes", axes)
		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})
}
