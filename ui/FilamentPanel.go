package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/utils"
)

var filamentPanelInstance *filamentPanel

type filamentPanel struct {
	CommonPanel

	amountStepButton     *StepButton
	box                  *gtk.Box
	labels               map[string]*LabelWithImage
	previous             *octoprint.TemperatureState
}

func FilamentPanel(ui *UI, parent Panel) Panel {
	if filamentPanelInstance == nil {
		m := &filamentPanel{CommonPanel: NewCommonPanel(ui, parent),
			labels: map[string]*LabelWithImage{},
		}
		m.b = NewBackgroundTask(time.Second * 5, m.updateTemperatures)
		m.initialize()
		filamentPanelInstance = m
	}

	return filamentPanelInstance
}

func (m *filamentPanel) initialize() {
	defer m.Initialize()

	toolheadCount := utils.GetToolheadCount(m.UI.Printer)
	if toolheadCount > 1 {
		m.createToolheadButtons()
	} else {
		extrudeButton := m.createExtrudeButton("Extrude", "extruder-extrude.svg", 1)
		m.Grid().Attach(extrudeButton, 0, 0, 1, 1)

		m.amountStepButton = MustStepButton(
			"move-step.svg",
			Step{" 20mm",  20},
			Step{" 50mm",  50},
			Step{"100mm", 100},
			Step{"  1mm",   1},
			Step{"  5mm",   5},
			Step{" 10mm",  10},
		)
		m.Grid().Attach(m.amountStepButton, 1, 0, 1, 1)

		flowRateButton := m.createFlowRateButton()
		m.Grid().Attach(flowRateButton, 2, 0, 1, 1)

		retractButton := m.createExtrudeButton("Retract", "extruder-retract.svg", -1) 
		m.Grid().Attach(retractButton, 3, 0, 1, 1)
	}


	m.Grid().Attach(m.createLoadButton(),   0, 1, 1, 1)

	m.box = MustBox(gtk.ORIENTATION_VERTICAL, 5)
	m.box.SetVAlign(gtk.ALIGN_CENTER)
	m.box.SetHAlign(gtk.ALIGN_CENTER)
	m.Grid().Attach(m.box, 1, 1, 2, 2)

	m.Grid().Attach(m.createUnloadButton(), 3, 1, 1, 1)


	m.Grid().Attach(MustButtonImageStyle("Temperature", "heat-up.svg", "color4", m.showTemperature), 0, 2, 1, 1)
}

func (m *filamentPanel) updateTemperatures() {
	s, err := (&octoprint.ToolStateRequest{
		History: true,
		Limit:   1,
	}).Do(m.UI.Printer)

	if err != nil {
		utils.LogError("filament.updateTemperatures()", "Do(ToolStateRequest)", err)
		return
	}

	m.loadTemperatureState(s)
}

func (m *filamentPanel) loadTemperatureState(s *octoprint.TemperatureState) {
	for tool, current := range s.Current {
		if _, ok := m.labels[tool]; !ok {
			m.addNewTool(tool)
		}

		m.loadTemperatureData(tool, &current)
	}

	m.previous = s
}

func (m *filamentPanel) addNewTool(tool string) {
	m.labels[tool] = MustLabelWithImage("toolhead.svg", "")
	m.box.Add(m.labels[tool])
	utils.Logger.Infof("New tool detected %s", tool)
}

func (m *filamentPanel) loadTemperatureData(tool string, d *octoprint.TemperatureData) {
	displayNameForTool := utils.GetDisplayNameForTool(tool)
	text := fmt.Sprintf("%.1f°C / %.1f°C", d.Actual, d.Target)
	toolheadCount := utils.GetToolheadCount(m.UI.Printer)
	if toolheadCount > 1 {
		text = (strings.Title(displayNameForTool) + ": " + text)
	}

	if m.previous != nil && d.Target > 0 {
		if p, ok := m.previous.Current[tool]; ok {
			text = fmt.Sprintf("%s (%.1f°C)", text, d.Actual - p.Actual)
		}
	}

	m.labels[tool].Label.SetText(text)
	m.labels[tool].ShowAll()
}

func (m *filamentPanel) createExtrudeButton(label, image string, dir int) gtk.IWidget {
	return MustPressedButton(label, image, func() {
		cmd := &octoprint.ToolExtrudeRequest{}
		cmd.Amount = m.amountStepButton.Value().(int) * dir

		utils.Logger.Infof("Sending extrude request, with amount %d", cmd.Amount)
		if err := cmd.Do(m.UI.Printer); err != nil {
			utils.LogError("filament.createExtrudeButton()", "Do(ToolExtrudeRequest)", err)
			return
		}
	}, 200)
}

func (m *filamentPanel) createFlowRateButton() *StepButton {
	b := MustStepButton(
		"speed-normal.svg",
		Step{"Normal (100%)", 100},
		Step{"Fast (125%)", 125},
		Step{"Slow (75%)", 75},
	)

	b.Callback = func() {
		cmd := &octoprint.ToolFlowrateRequest{}
		cmd.Factor = b.Value().(int)

		utils.Logger.Infof("Changing flowrate to %d%%", cmd.Factor)
		if err := cmd.Do(m.UI.Printer); err != nil {
			utils.LogError("filament.createFlowRateButton()", "Do(ToolFlowrateRequest)", err)
			return
		}
	}

	return b
}

func (m *filamentPanel) createLoadButton() gtk.IWidget {
	length := 750.0
	if m.UI.Settings != nil {
		length = m.UI.Settings.FilamentInLength
	}

	return MustButtonImageStyle("Load", "filament-spool-load.svg", "", func() {
		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{
			"G91",
			fmt.Sprintf("G0 E%.1f F5000", length * 0.80),
			fmt.Sprintf("G0 E%.1f F500", length * 0.20),
			"G90",
		}

		utils.Logger.Info("Sending filament load request")
		if err := cmd.Do(m.UI.Printer); err != nil {
			utils.LogError("filament.createLoadButton()", "Do(CommandRequest)", err)
			return
		}
	})
}

func (m *filamentPanel) createUnloadButton() gtk.IWidget {
	length := 800.0
	if m.UI.Settings != nil {
		length = m.UI.Settings.FilamentOutLength
	}

	return MustButtonImageStyle("Unload", "filament-spool-unload.svg", "", func() {
		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{
			"G91",
			fmt.Sprintf("G0 E-%.1f F5000", length),
			"G90",
		}

		utils.Logger.Info("Sending filament unload request")
		if err := cmd.Do(m.UI.Printer); err != nil {
			utils.LogError("filament.createUnloadButton()", "Do(CommandRequest)", err)
			return
		}
	})
}

func (m *filamentPanel) showTemperature() {
	m.UI.Add(TemperaturePanel(m.UI, m))
}

func (m *filamentPanel) createToolheadButtons() {
	toolheadCount := utils.GetToolheadCount(m.UI.Printer)
	toolheadButtons := CreateChangeToolheadButtonsAndAttachToGrid(toolheadCount, m.Grid())
	m.setToolheadButtonClickHandlers(toolheadButtons)
}

func (m *filamentPanel) setToolheadButtonClickHandlers(toolheadButtons []*gtk.Button) {
	for index, toolheadButton := range toolheadButtons {
		m.setToolheadButtonClickHandler(toolheadButton, index)
	}
}

func (m *filamentPanel) setToolheadButtonClickHandler(toolheadButton *gtk.Button, toolheadIndex int) {
	toolheadButton.Connect("clicked", func() {
		utils.Logger.Infof("Changing tool to tool%d", toolheadIndex)

		gcode := fmt.Sprintf("T%d", toolheadIndex)
		m.command(gcode)
	})
}
