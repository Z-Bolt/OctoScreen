package ui

import (
	"fmt"
	"sync"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/utils"
)

var idleStatusPanelInstance *idleStatusPanel

type idleStatusPanel struct {
	CommonPanel
	//step           *StepButton
	//pb             *gtk.ProgressBar

	tool0, tool1, tool2, tool3, bed *ToolHeatup
}

func IdleStatusPanel(ui *UI) Panel {
	if idleStatusPanelInstance == nil {
		m := &idleStatusPanel{
			CommonPanel: NewTopLevelCommonPanel(ui, nil),
		}
		m.b = NewBackgroundTask(time.Second * 2, m.update)
		m.initialize()

		idleStatusPanelInstance = m
	}

	return idleStatusPanelInstance
}

func (m *idleStatusPanel) initialize() {
	defer m.Initialize()

	utils.Logger.Info(m.UI.Settings)

	var menuItems []octoprint.MenuItem
	if m.UI.Settings == nil || len(m.UI.Settings.MenuStructure) == 0 {
		utils.Logger.Info("Loading default menu")
		m.UI.Settings.MenuStructure = getDefaultMenu()
	} else {
		utils.Logger.Info("Loading octo menu")
	}

	menuItems = m.UI.Settings.MenuStructure

	menuGrid := MustGrid()
	menuGrid.SetRowHomogeneous(true)
	menuGrid.SetColumnHomogeneous(true)
	m.Grid().Attach(menuGrid, 2, 0, 2, 2)
	m.arrangeMenuItems(menuGrid, menuItems, 2)

	printButton := MustButtonImageStyle("Print", "print.svg", "color2", m.showFiles)
	m.Grid().Attach(printButton, 2, 2, 2, 1)

	m.showTools()
}

func (m *idleStatusPanel) showFiles() {
	m.UI.Add(FilesPanel(m.UI, m))
}

func (m *idleStatusPanel) update() {
	m.updateTemperature()
}

func (m *idleStatusPanel) showTools() {
	toolheadCount := utils.GetToolheadCount(m.UI.Printer)

	if toolheadCount == 1 {
		m.tool0 = creteToolHeatupButton(0, m.UI.Printer)
	} else {
		m.tool0 = creteToolHeatupButton(1, m.UI.Printer)
	}
	m.tool1 = creteToolHeatupButton(2, m.UI.Printer)
	m.tool2 = creteToolHeatupButton(3, m.UI.Printer)
	m.tool3 = creteToolHeatupButton(4, m.UI.Printer)
	m.bed = creteToolHeatupButton(-1, m.UI.Printer)

	switch toolheadCount {
		case 1:
			g := MustGrid()
			g.SetRowHomogeneous(true)
			g.SetColumnHomogeneous(true)
			m.Grid().Attach(g, 0, 0, 2, 3)
			g.Attach(m.tool0,  0, 0, 2, 1)
			g.Attach(m.bed,    0, 1, 2, 1)

		case 2:
			m.Grid().Attach(m.tool0, 0, 0, 2, 1)
			m.Grid().Attach(m.tool1, 0, 1, 2, 1)
			m.Grid().Attach(m.bed,   0, 2, 2, 1)

		case 3:
			m.Grid().Attach(m.tool0, 0, 0, 1, 1)
			m.Grid().Attach(m.tool1, 1, 0, 1, 1)
			m.Grid().Attach(m.tool2, 0, 1, 2, 1)
			m.Grid().Attach(m.bed,   0, 2, 2, 1)

		case 4:
			m.Grid().Attach(m.tool0, 0, 0, 1, 1)
			m.Grid().Attach(m.tool1, 1, 0, 1, 1)
			m.Grid().Attach(m.tool2, 0, 1, 1, 1)
			m.Grid().Attach(m.tool3, 1, 1, 1, 1)
			m.Grid().Attach(m.bed,   0, 2, 2, 1)
	}




	// if toolheadCount == 1 {
	// 	m.tool0 = creteToolHeatupButton(0, m.UI.Printer)
	// } else {
	// 	m.tool0 = creteToolHeatupButton(1, m.UI.Printer)
	// }

	// m.tool1 = creteToolHeatupButton(2, m.UI.Printer)
	// m.tool2 = creteToolHeatupButton(3, m.UI.Printer)
	// m.tool3 = creteToolHeatupButton(4, m.UI.Printer)
	// m.bed   = creteToolHeatupButton(-1, m.UI.Printer)

	// m.Grid().Attach(m.tool0, 0, 0, 1, 1)
	// if toolheadCount >= 2 {
	// 	m.Grid().Attach(m.tool1, 1, 0, 1, 1)
	// }

	// if toolheadCount >= 3 {
	// 	m.Grid().Attach(m.tool2, 0, 1, 1, 1)
	// }

	// if toolheadCount >= 4 {
	// 	m.Grid().Attach(m.tool3, 1, 1, 1, 1)
	// }

	// m.Grid().Attach(m.bed, 0, 2, 1, 1)
}

func (m *idleStatusPanel) updateTemperature() {
	s, err := (&octoprint.StateRequest{Exclude: []string{"sd"}}).Do(m.UI.Printer)
	if err != nil {
		utils.LogError("idle_status.updateTemperature()", "Do(StateRequest)", err)
		return
	}

	for tool, s := range s.Temperature.Current {
		switch tool {
			case "bed":
				m.bed.setTemperatures(s.Actual, s.Target)

			case "tool0":
				m.tool0.setTemperatures(s.Actual, s.Target)

			case "tool1":
				m.tool1.setTemperatures(s.Actual, s.Target)

			case "tool2":
				m.tool2.setTemperatures(s.Actual, s.Target)

			case "tool3":
				m.tool3.setTemperatures(s.Actual, s.Target)
		}
	}
}

type ToolHeatup struct {
	isHeating bool
	tool      string
	*gtk.Button
	sync.RWMutex
	printer *octoprint.Client
}

func creteToolHeatupButton(num int, printer *octoprint.Client) *ToolHeatup {
	var (
		image string
		tool  string
	)

	if num < 0 {
		image = "bed.svg"
		tool = "bed"
	} else if num == 0 {
		image = "toolhead.svg"
		tool = "tool0"
	} else {
		image = fmt.Sprintf("toolhead-%d.svg", num)
		tool = fmt.Sprintf("tool%d", num - 1)
	}

	t := &ToolHeatup{
		Button:  MustButtonImage("", image, nil),
		tool:    tool,
		printer: printer,
	}

	_, err := t.Connect("clicked", t.clicked)
	if err != nil {
		utils.LogError("idle_status.creteToolHeatupButton()", "t.Connect('clicked', t.clicked)", err)
	}

	return t
}

func (t *ToolHeatup) updateStatus(heating bool) {
	ctx, _ := t.GetStyleContext()
	if heating {
		ctx.AddClass("active")
	} else {
		ctx.RemoveClass("active")
	}
	t.isHeating = heating
}

func (t *ToolHeatup) setTemperatures(actual float64, target float64) {
	text := fmt.Sprintf("%.0f°C / %.0f°C", actual, target)
	t.SetLabel(text)
	t.updateStatus(target > 0)
}

func (t *ToolHeatup) getProfileTemperature() float64 {
	temperature := 0.0

	s, err := (&octoprint.SettingsRequest{}).Do(t.printer)
	if err != nil {
		utils.LogError("idle_status.getProfileTemperature()", "Do(SettingsRequest)", err)
		return 0
	}

	if len(s.Temperature.Profiles) > 0 {
		if t.tool == "bed" {
			temperature = s.Temperature.Profiles[0].Bed
		} else {
			temperature = s.Temperature.Profiles[0].Extruder
		}
	} else {
		if t.tool == "bed" {
			temperature = 75
		} else {
			temperature = 220
		}
	}

	return temperature
}

func (t *ToolHeatup) clicked() {
	defer func() { t.updateStatus(!t.isHeating) }()

	var (
		target float64
		err    error
	)

	if t.isHeating {
		target = 0.0
	} else {
		target = t.getProfileTemperature()
	}

	if t.tool == "bed" {
		cmd := &octoprint.BedTargetRequest{Target: target}
		err = cmd.Do(t.printer)
		if err != nil {
			utils.LogError("idle_status.clicked()", "Do(BedTargetRequest)", err)
		}
	} else {
		cmd := &octoprint.ToolTargetRequest{Targets: map[string]float64{t.tool: target}}
		err = cmd.Do(t.printer)
		if err != nil {
			utils.LogError("idle_status.clicked()", "Do(ToolTargetRequest)", err)
		}
	}
}
