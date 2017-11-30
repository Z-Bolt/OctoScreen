package ui

import (
	"fmt"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

type FilamentPanel struct {
	CommonPanel

	amount *StepButton
	tool   *StepButton

	box      *gtk.Box
	labels   map[string]*gtk.Label
	previous *octoprint.TemperatureState
}

func NewFilamentPanel(ui *UI) *FilamentPanel {
	m := &FilamentPanel{CommonPanel: NewCommonPanel(ui),
		labels: map[string]*gtk.Label{},
	}

	m.b = NewBackgroundTask(time.Second*5, m.updateTemperatures)
	m.initialize()
	return m
}

func (m *FilamentPanel) initialize() {
	m.grid.Attach(m.createExtrudeButton("Extrude", "extrude.svg", 1), 1, 0, 1, 1)
	m.grid.Attach(m.createExtrudeButton("Retract", "retract.svg", -1), 4, 0, 1, 1)

	m.box = MustBox(gtk.ORIENTATION_VERTICAL, 5)
	m.box.SetVAlign(gtk.ALIGN_CENTER)
	m.grid.Attach(m.box, 2, 0, 2, 1)

	m.amount = MustStepButton("move-step.svg", Step{"5mm", 5}, Step{"10mm", 10}, Step{"1mm", 1})
	m.grid.Attach(m.amount, 2, 1, 1, 1)

	m.grid.Attach(m.createToolButton(), 1, 1, 1, 1)
	m.grid.Attach(m.createFlowrateButton(), 3, 1, 1, 1)
	m.grid.Attach(MustButtonImage("Return", "back.svg", m.UI.ShowDefaultPanel), 4, 1, 1, 1)

	m.grid.Connect("show", m.Show)
}

func (m *FilamentPanel) updateTemperatures() {
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

func (m *FilamentPanel) loadTemperatureState(s *octoprint.TemperatureState) {
	for tool, current := range s.Current {
		if _, ok := m.labels[tool]; !ok {
			m.addNewTool(tool)
		}

		m.loadTemperatureData(tool, &current)
	}

	m.previous = s
}

func (m *FilamentPanel) addNewTool(tool string) {
	m.labels[tool] = MustLabel("")
	m.box.Add(m.labels[tool])
	m.tool.AddStep(Step{tool, tool})

	Logger.Infof("New tool detected %s", tool)
}

func (m *FilamentPanel) loadTemperatureData(tool string, d *octoprint.TemperatureData) {
	text := fmt.Sprintf("%s: %.1f°C / %.1f°C", tool, d.Actual, d.Target)

	if m.previous != nil && d.Target > 0 {
		if p, ok := m.previous.Current[tool]; ok {
			text = fmt.Sprintf("%s (%.1f°C)", text, d.Actual-p.Actual)
		}
	}

	m.labels[tool].SetText(text)
	m.labels[tool].ShowAll()
}

func (m *FilamentPanel) createToolButton() *StepButton {
	m.tool = MustStepButton("extruct.svg")
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

func (m *FilamentPanel) createFlowrateButton() *StepButton {
	b := MustStepButton("speed-step.svg", Step{"Normal", 100}, Step{"High", 125}, Step{"Slow", 75})
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

func (m *FilamentPanel) createExtrudeButton(label, image string, dir int) gtk.IWidget {
	return MustButtonImage(label, image, func() {
		cmd := &octoprint.ToolExtrudeRequest{}
		cmd.Amount = m.amount.Value().(int) * dir

		Logger.Infof("Sending extrude request, with amount %d", cmd.Amount)
		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})
}
