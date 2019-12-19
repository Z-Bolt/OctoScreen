package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

var filamentMultitoolPanelInstance *filamentMultitoolPanel

type filamentMultitoolPanel struct {
	CommonPanel

	box      *gtk.Box
	labels   map[string]*LabelWithImage
	previous *octoprint.TemperatureState
}

func FilamentMultitoolPanel(ui *UI, parent Panel) Panel {
	if filamentMultitoolPanelInstance == nil {
		m := &filamentMultitoolPanel{CommonPanel: NewCommonPanel(ui, parent),
			labels: map[string]*LabelWithImage{},
		}
		m.panelH = 3
		m.b = NewBackgroundTask(time.Second*5, m.updateTemperatures)
		m.initialize()
		filamentMultitoolPanelInstance = m
	}

	return filamentMultitoolPanelInstance
}

func (m *filamentMultitoolPanel) initialize() {
	defer m.Initialize()

	m.Grid().Attach(m.createChangeToolButton(0), 1, 0, 1, 1)
	m.Grid().Attach(m.createChangeToolButton(1), 2, 0, 1, 1)
	m.Grid().Attach(m.createChangeToolButton(2), 3, 0, 1, 1)
	m.Grid().Attach(m.createChangeToolButton(3), 4, 0, 1, 1)

	m.Grid().Attach(m.createLoadButton(), 1, 1, 1, 1)
	m.Grid().Attach(m.createUnloadButton(), 4, 1, 1, 1)

	m.Grid().Attach(MustButtonImageStyle("Temperature", "heat-up.svg", "color4", m.showTemperature), 1, 2, 1, 1)

	m.box = MustBox(gtk.ORIENTATION_VERTICAL, 5)
	m.box.SetVAlign(gtk.ALIGN_CENTER)
	m.box.SetHAlign(gtk.ALIGN_CENTER)

	m.Grid().Attach(m.box, 2, 1, 2, 2)

}

func (m *filamentMultitoolPanel) updateTemperatures() {
	s, err := (&octoprint.ToolStateRequest{
		History: true,
		Limit:   1,
	}).Do(m.UI.Printer)

	if err != nil {
		Logger.Error(err)
		return
	}

	m.loadTemperatureState(s)
}

func (m *filamentMultitoolPanel) loadTemperatureState(s *octoprint.TemperatureState) {
	for tool, current := range s.Current {
		if _, ok := m.labels[tool]; !ok {
			m.addNewTool(tool)
		}

		m.loadTemperatureData(tool, &current)
	}

	m.previous = s
}

func (m *filamentMultitoolPanel) addNewTool(tool string) {
	m.labels[tool] = MustLabelWithImage("extruder.svg", "")
	m.box.Add(m.labels[tool])

	Logger.Infof("New tool detected %s", tool)
}

func (m *filamentMultitoolPanel) loadTemperatureData(tool string, d *octoprint.TemperatureData) {
	text := fmt.Sprintf("%s: %.1f°C / %.1f°C", strings.Title(tool), d.Actual, d.Target)

	if m.previous != nil && d.Target > 0 {
		if p, ok := m.previous.Current[tool]; ok {
			text = fmt.Sprintf("%s (%.1f°C)", text, d.Actual-p.Actual)
		}
	}

	m.labels[tool].Label.SetText(text)
	m.labels[tool].ShowAll()
}

func (m *filamentMultitoolPanel) createLoadButton() gtk.IWidget {
	length := 750.0

	if m.UI.Settings != nil {
		length = m.UI.Settings.FilamentInLength
	}

	return MustButtonImageStyle("Load", "extrude.svg", "color3", func() {
		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{
			"G91",
			fmt.Sprintf("G0 E%.1f F5000", length*0.80),
			fmt.Sprintf("G0 E%.1f F500", length*0.20),
			"G90",
		}

		Logger.Info("Sending filament load request")
		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})
}

func (m *filamentMultitoolPanel) createUnloadButton() gtk.IWidget {
	length := 800.0

	if m.UI.Settings != nil {
		length = m.UI.Settings.FilamentOutLength
	}

	return MustButtonImageStyle("Unload", "extrude.svg", "color2", func() {
		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{
			"G91",
			fmt.Sprintf("G0 E-%.1f F5000", length),
			"G90",
		}

		Logger.Info("Sending filament unload request")
		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})
}

func (m *filamentMultitoolPanel) createChangeToolButton(num int) gtk.IWidget {
	style := fmt.Sprintf("color%d", num+1)
	name := fmt.Sprintf("Tool%d", num+1)
	gcode := fmt.Sprintf("T%d", num)
	return MustButtonImageStyle(name, "extruder.svg", style, func() {
		m.command(gcode)
	})
}

func (m *filamentMultitoolPanel) showTemperature() {
	m.UI.Add(TemperaturePanel(m.UI, m))
}
