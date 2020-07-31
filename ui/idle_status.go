package ui

import (
	"fmt"
	"sync"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

var idleStatusPanelInstance *idleStatusPanel

type idleStatusPanel struct {
	CommonPanel
	step       *StepButton
	pb         *gtk.ProgressBar
	toolsCount int

	tool0, tool1, tool2, tool3, bed *ToolHeatup
}

func IdleStatusPanel(ui *UI) Panel {
	if idleStatusPanelInstance == nil {
		m := &idleStatusPanel{CommonPanel: NewCommonPanel(ui, nil)}
		m.panelH = 3
		m.b = NewBackgroundTask(time.Second*2, m.update)
		m.initialize()

		idleStatusPanelInstance = m
	}

	return idleStatusPanelInstance
}

func (m *idleStatusPanel) initialize() {
	defer m.Initialize()

	var menuItems []octoprint.MenuItem

	Logger.Info(m.UI.Settings)

	if m.UI.Settings == nil || len(m.UI.Settings.MenuStructure) == 0 {
		Logger.Info("Loading default menu")
		menuItems = getDefaultMenu()
	} else {
		Logger.Info("Loading octo menu")
		menuItems = m.UI.Settings.MenuStructure
	}

	buttons := MustGrid()
	buttons.SetRowHomogeneous(true)
	buttons.SetColumnHomogeneous(true)
	m.Grid().Attach(buttons, 3, 0, 2, 2)

	m.arrangeMenuItems(buttons, menuItems, 2)

	m.Grid().Attach(MustButtonImageStyle("Print", "print.svg", "color2", m.showFiles), 3, 2, 2, 1)

	m.showTools()
}

func (m *idleStatusPanel) showFiles() {
	m.UI.Add(FilesPanel(m.UI, m))
}

func (m *idleStatusPanel) update() {
	m.updateTemperature()
}

func (m *idleStatusPanel) showTools() {
	toolsCount := m.defineToolsCount()

	m.tool0 = ToolHeatupNew(0, m.UI.Printer)
	m.tool1 = ToolHeatupNew(1, m.UI.Printer)
	m.tool2 = ToolHeatupNew(2, m.UI.Printer)
	m.tool3 = ToolHeatupNew(3, m.UI.Printer)
	m.bed = ToolHeatupNew(-1, m.UI.Printer)

	switch toolsCount {
	case 1:
		g := MustGrid()
		g.SetRowHomogeneous(true)
		g.SetColumnHomogeneous(true)
		m.Grid().Attach(g, 1, 0, 2, 3)
		g.Attach(m.tool0,  1, 0, 2, 1)
		g.Attach(m.bed,    1, 1, 2, 1)

	case 2:
		m.Grid().Attach(m.tool0, 1, 0, 2, 1)
		m.Grid().Attach(m.tool1, 1, 1, 2, 1)
		m.Grid().Attach(m.bed,   1, 2, 2, 1)

	case 3:
		m.Grid().Attach(m.tool0, 1, 0, 1, 1)
		m.Grid().Attach(m.tool1, 2, 0, 1, 1)
		m.Grid().Attach(m.tool2, 1, 1, 2, 1)
		m.Grid().Attach(m.bed,   1, 2, 2, 1)

	case 4:
		m.Grid().Attach(m.tool0, 1, 0, 1, 1)
		m.Grid().Attach(m.tool1, 2, 0, 1, 1)
		m.Grid().Attach(m.tool2, 1, 1, 1, 1)
		m.Grid().Attach(m.tool3, 2, 1, 1, 1)
		m.Grid().Attach(m.bed,   1, 2, 2, 1)
	}
}

func (m *idleStatusPanel) updateTemperature() {
	s, err := (&octoprint.StateRequest{Exclude: []string{"sd"}}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return
	}

	for tool, s := range s.Temperature.Current {
		switch tool {
		case "bed":
			m.bed.SetTemperatures(s.Actual, s.Target)
		case "tool0":
			m.tool0.SetTemperatures(s.Actual, s.Target)
		case "tool1":
			m.tool1.SetTemperatures(s.Actual, s.Target)
		case "tool2":
			m.tool2.SetTemperatures(s.Actual, s.Target)
		case "tool3":
			m.tool3.SetTemperatures(s.Actual, s.Target)
		}
	}
}

func (m *idleStatusPanel) defineToolsCount() int {
	c, err := (&octoprint.ConnectionRequest{}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return 0
	}

	profile, err := (&octoprint.PrinterProfilesRequest{Id: c.Current.PrinterProfile}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return 0
	}

	if profile.Extruder.SharedNozzle {
		return 1
	}

	return profile.Extruder.Count
}

type ToolHeatup struct {
	isHeating bool
	tool      string
	*gtk.Button
	sync.RWMutex
	printer *octoprint.Client
}

func ToolHeatupNew(num int, printer *octoprint.Client) *ToolHeatup {
	var (
		image string
		tool  string
	)

	if num < 0 {
		image = "bed.svg"
		tool = "bed"
	} else {
		image = fmt.Sprintf("extruder-%d.svg", num+1)
		tool = fmt.Sprintf("tool%d", num)
	}

	t := &ToolHeatup{
		Button:  MustButtonImage("", image, nil),
		tool:    tool,
		printer: printer,
	}

	t.Connect("clicked", t.clicked)

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

func (t *ToolHeatup) SetTemperatures(actual float64, target float64) {
	text := fmt.Sprintf("%.0f°C / %.0f°C", actual, target)
	t.SetLabel(text)
	t.updateStatus(target > 0)
}

func (t *ToolHeatup) getProfileTemperature() float64 {
	temperature := 0.0

	s, err := (&octoprint.SettingsRequest{}).Do(t.printer)
	if err != nil {
		Logger.Error(err)
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
	} else {
		cmd := &octoprint.ToolTargetRequest{Targets: map[string]float64{t.tool: target}}
		err = cmd.Do(t.printer)
	}

	if err != nil {
		Logger.Error(err)
	}
}
