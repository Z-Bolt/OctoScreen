package ui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/utils"
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

	m.addBedLevelCornerButton(true, true)
	m.addBedLevelCornerButton(false, true)
	m.addBedLevelCornerButton(true, false)
	m.addBedLevelCornerButton(false, false)

	if m.UI.Settings != nil && m.UI.Settings.GCodes.AutoBedLevel != "" {
		m.Grid().Attach(m.createAutoLevelButton(m.UI.Settings.GCodes.AutoBedLevel), 3, 0, 1, 1)
	}
}

func (m *bedLevelPanel) addBedLevelCornerButton(isLeft, isTop bool) {
	x := 1
	y := 1
	placement := "b-"
	if isTop {
		placement = "t-"
		y = 0
	}

	if isLeft {
		placement += "l"
	} else {
		placement += "r"
		x = 2
	}

	button := m.createLevelButton(placement)
	if isLeft {
		button.SetHAlign(gtk.ALIGN_END)
	} else {
		button.SetHAlign(gtk.ALIGN_START)
	}

	if isTop {
		button.SetVAlign(gtk.ALIGN_END)
		button.SetImagePosition(gtk.POS_BOTTOM)
	} else {
		button.SetVAlign(gtk.ALIGN_START)
		button.SetImagePosition(gtk.POS_TOP)
	}

	styleContext, _ := button.GetStyleContext()
	styleContext.AddClass("no-margin")
	styleContext.AddClass("no-border")
	styleContext.AddClass("no-padding")

	if isLeft {
		styleContext.AddClass("padding-left-20")
	} else {
		styleContext.AddClass("padding-right-20")
	}

	m.Grid().Attach(button, x, y, 1, 1)
}

func (m *bedLevelPanel) defineLevelingPoints() {
	c, err := (&octoprint.ConnectionRequest{}).Do(m.UI.Printer)
	if err != nil {
		utils.LogError("bed-level.defineLevelingPoints()", "Do(ConnectionRequest)", err)
		return
	}

	utils.Logger.Info(c.Current.PrinterProfile)

	profile, err := (&octoprint.PrinterProfilesRequest{Id: c.Current.PrinterProfile}).Do(m.UI.Printer)
	if err != nil {
		utils.LogError("bed-level.defineLevelingPoints()", "Do(PrinterProfilesRequest)", err)
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

func (m *bedLevelPanel) createLevelButton(placement string) *gtk.Button {
	imageFileName := fmt.Sprintf("bed-level-%s-65%%.svg", placement)
	noLabel := ""
	b := MustButtonImage(noLabel, imageFileName, func() {
		m.goHomeIfRequire()

		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{
			"G0 Z10 F2000",
			fmt.Sprintf("G0 X%f Y%f F10000", m.points[placement][0], m.points[placement][1]),
			"G0 Z0 F400",
		}

		if err := cmd.Do(m.UI.Printer); err != nil {
			utils.LogError("bed-level.createLevelButton()", "Do(CommandRequest)", err)
			return
		}
	})

	return b
}

func (m *bedLevelPanel) goHomeIfRequire() {
	if m.homed {
		return
	}

	cmd := &octoprint.CommandRequest{}
	cmd.Commands = []string{
		"G28",
	}

	if err := cmd.Do(m.UI.Printer); err != nil {
		utils.LogError("bed-level.goHomeIfRequire()", "Do(CommandRequest)", err)
		return
	}

	m.homed = true
}

func (m *bedLevelPanel) createAutoLevelButton(gcode string) *gtk.Button {
	b := MustButtonImage("Auto Level", "bed-level.svg", func() {
		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{
			gcode,
		}

		if err := cmd.Do(m.UI.Printer); err != nil {
			utils.LogError("bed-level.createAutoLevelButton()", "Do(CommandRequest)", err)
			return
		}
	})

	return b
}
