package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

var extrudePanelInstance *extrudePanel

type extrudePanel struct {
	CommonPanel

	amount *StepButton
	tool   *StepButton

	box      *gtk.Box
	labels   map[string]*LabelWithImage
	previous *octoprint.TemperatureState
}

func ExtrudePanel(ui *UI, parent Panel) Panel {
	if extrudePanelInstance == nil {
		m := &extrudePanel{CommonPanel: NewCommonPanel(ui, parent),
			labels: map[string]*LabelWithImage{},
		}
		m.panelH = 2
		m.b = NewBackgroundTask(time.Second*5, m.updateTemperatures)
		m.initialize()
		extrudePanelInstance = m
	}

	return extrudePanelInstance
}

func (m *extrudePanel) initialize() {
	defer m.Initialize()

	m.Grid().Attach(m.createExtrudeButton("Extrude", "extrude.svg", 1), 1, 0, 1, 1)
	m.Grid().Attach(m.createExtrudeButton("Retract", "retract.svg", -1), 4, 0, 1, 1)

	m.box = MustBox(gtk.ORIENTATION_VERTICAL, 5)
	m.box.SetVAlign(gtk.ALIGN_CENTER)
	m.box.SetHAlign(gtk.ALIGN_CENTER)
	m.Grid().Attach(m.box, 2, 0, 2, 1)

	m.amount = MustStepButton("move-step.svg", Step{"1mm", 1}, Step{"5mm", 5}, Step{"10mm", 10})
	m.Grid().Attach(m.amount, 2, 1, 1, 1)

	m.Grid().Attach(m.createToolButton(), 1, 1, 1, 1)
	m.Grid().Attach(m.createFlowrateButton(), 3, 1, 1, 1)
}

func (m *extrudePanel) updateTemperatures() {
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

func (m *extrudePanel) loadTemperatureState(s *octoprint.TemperatureState) {
	for tool, current := range s.Current {
		if _, ok := m.labels[tool]; !ok {
			m.addNewTool(tool)
		}

		m.loadTemperatureData(tool, &current)
	}

	m.previous = s
}

func (m *extrudePanel) addNewTool(tool string) {
	m.labels[tool] = MustLabelWithImage("extruder.svg", "")
	m.box.Add(m.labels[tool])
	m.tool.AddStep(Step{strings.Title(tool), tool})

	Logger.Infof("New tool detected %s", tool)
}

func (m *extrudePanel) loadTemperatureData(tool string, d *octoprint.TemperatureData) {
	text := fmt.Sprintf("%s: %.1f°C / %.1f°C", strings.Title(tool), d.Actual, d.Target)

	if m.previous != nil && d.Target > 0 {
		if p, ok := m.previous.Current[tool]; ok {
			text = fmt.Sprintf("%s (%.1f°C)", text, d.Actual-p.Actual)
		}
	}

	m.labels[tool].Label.SetText(text)
	m.labels[tool].ShowAll()
}

func (m *extrudePanel) createToolButton() *StepButton {
	m.tool = MustStepButton("extruder.svg")
	m.tool.Callback = func() {
		cmd := &octoprint.ToolSelectRequest{}
		cmd.Tool = m.tool.Value().(string)

		Logger.Infof("Changing tool to %s", cmd.Tool)
		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	}

	return m.tool
}

func (m *extrudePanel) createFlowrateButton() *StepButton {
	b := MustStepButton("speed-step.svg", Step{"Slow", 75}, Step{"Normal", 100}, Step{"High", 125})
	b.Callback = func() {
		cmd := &octoprint.ToolFlowrateRequest{}
		cmd.Factor = b.Value().(int)

		Logger.Infof("Changing flowrate to %d%%", cmd.Factor)
		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	}

	return b
}

func (m *extrudePanel) createLoadButton() gtk.IWidget {
	length := 750.0

	if m.UI.Settings != nil {
		length = m.UI.Settings.FilamentInLength
	}

	return MustButtonImage("Load", "extrude.svg", func() {
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

func (m *extrudePanel) createUnloadButton() gtk.IWidget {

	length := 800.0

	if m.UI.Settings != nil {
		length = m.UI.Settings.FilamentOutLength
	}

	return MustButtonImage("Unload", "retract.svg", func() {
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

func (m *extrudePanel) createExtrudeButton(label, image string, dir int) gtk.IWidget {

	return MustPressedButton(label, image, func() {
		cmd := &octoprint.ToolExtrudeRequest{}
		cmd.Amount = m.amount.Value().(int) * dir

		Logger.Infof("Sending extrude request, with amount %d", cmd.Amount)
		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	}, 200)
}
