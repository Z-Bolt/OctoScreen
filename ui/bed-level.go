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
	homed  bool
}

func BedLevelPanel(ui *UI, parent Panel) Panel {
	if bedLevelPanelInstance == nil {
		m := &bedLevelPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.initialize()
		bedLevelPanelInstance = m
	}

	bedLevelPanelInstance.homed = false
	return bedLevelPanelInstance
}

func (m *bedLevelPanel) initialize() {
	defer m.Initialize()

	m.defineLevelingPoints()

	m.Grid().Attach(m.createLevelButton("t-l"), 2, 0, 1, 1)
	m.Grid().Attach(m.createLevelButton("t-r"), 3, 0, 1, 1)
	m.Grid().Attach(m.createLevelButton("b-l"), 2, 1, 1, 1)
	m.Grid().Attach(m.createLevelButton("b-r"), 3, 1, 1, 1)
}

func (m *bedLevelPanel) defineLevelingPoints() {
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
		m.createHoveIfRequire()

		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{
			"G0 Z10 F2000",
			fmt.Sprintf("G0 X%f Y%f F10000", m.points[p][0], m.points[p][1]),
			"G0 Z0 F400",
		}

		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})
	return b
}

func (m *bedLevelPanel) createHoveIfRequire() {
	if m.homed {
		return
	}

	cmd := &octoprint.CommandRequest{}
	cmd.Commands = []string{
		"G28",
	}

	if err := cmd.Do(m.UI.Printer); err != nil {
		Logger.Error(err)
		return
	}

	m.homed = true
}
