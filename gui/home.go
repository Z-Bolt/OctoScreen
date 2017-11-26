package gui

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/OctoPrint-TFT/octoprint"
)

type HomeMenu struct {
	*gtk.Grid
	gui *GUI
}

func NewHomeMenu(gui *GUI) *HomeMenu {
	m := &HomeMenu{Grid: MustGrid(), gui: gui}
	m.initialize()
	return m
}

func (m *HomeMenu) initialize() {
	m.Attach(m.createMoveButton("Home All", "home.svg",
		octoprint.XAxis, octoprint.YAxis, octoprint.ZAxis,
	), 1, 0, 1, 1)

	m.Attach(m.createMoveButton("Home X", "home-x.svg", octoprint.XAxis), 2, 0, 1, 1)
	m.Attach(m.createMoveButton("Home Y", "home-y.svg", octoprint.YAxis), 3, 0, 1, 1)
	m.Attach(m.createMoveButton("Home Z", "home-z.svg", octoprint.ZAxis), 4, 0, 1, 1)
	m.Attach(MustButtonImage("Back", "back.svg", m.gui.ShowMenu), 4, 1, 1, 1)
}

func (m *HomeMenu) createMoveButton(label, image string, axes ...octoprint.Axis) gtk.IWidget {
	return MustButtonImage(label, image, func() {
		cmd := &octoprint.HomeCommand{Axes: axes}
		if err := cmd.Do(m.gui.Printer); err != nil {
			panic(err)
		}
	})
}
