package ui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

var bedLevelPanelInstance *bedLevelPanel

type bedLevelPanel struct {
	CommonPanel
	points map[string][]float64
}

func BedLevelPanel(ui *UI, parent Panel) Panel {
	if bedLevelPanelInstance == nil {
		m := &bedLevelPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.initialize()
		bedLevelPanelInstance = m
	}

	return bedLevelPanelInstance
}

func (m *bedLevelPanel) initialize() {
	defer m.Initialize()

	m.loadLevelingPoints()

	m.Grid().Attach(m.createLevelButton("t-l"), 2, 0, 1, 1)
	m.Grid().Attach(m.createLevelButton("t-r"), 3, 0, 1, 1)
	m.Grid().Attach(m.createLevelButton("b-l"), 2, 1, 1, 1)
	m.Grid().Attach(m.createLevelButton("b-r"), 3, 1, 1, 1)
}

func (m *bedLevelPanel) loadLevelingPoints() {
	c, err := (&octoprint.ConnectionRequest{}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return
	}

	Logger.Info(c.Current.PrinterProfile)

	profile, err := (&octoprint.PrinterProfilesRequest{Id: c.Current.PrinterProfile}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return
	}

	xMax := profile.Volume.Width
	yMax := profile.Volume.Depth
	xOffset := xMax * 0.1
	yOffset := yMax * 0.1

	m.points = map[string][]float64{
		"t-l": {xOffset, yMax - yOffset},
		"t-r": {xMax - xOffset, yMax - yOffset},
		"b-l": {xOffset, yOffset},
		"b-r": {xMax - xOffset, yOffset},
	}
}

func (m *bedLevelPanel) createLevelButton(p string) *gtk.Button {
	img := fmt.Sprintf("bed-level-%s.svg", p)
	b := MustButtonImage("", img, func() {
		gcode := fmt.Sprintf("G0 X%f Y%f", m.points[p][0], m.points[p][1])

		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{
			"G0 Z30",
			gcode,
			"G0 Z0",
		}

		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})
	return b
}
