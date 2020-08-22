/*
extruder.go and extruder_multitool.go were consolidated into this single file, and handles whether
the user has a single extruder or multiple extruders (maybe it should be renamed to extruder_universal
to denote that it handles both).
*/

package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/utils"
)

var extruderPanelInstance *extruderPanel

type extruderPanel struct {
	CommonPanel
	amount   *StepButton
	box      *gtk.Box
	labels   map[string]*LabelWithImage
	previous *octoprint.TemperatureState
}

func ExtruderPanel(ui *UI, parent Panel) Panel {
	if extruderPanelInstance == nil {
		m := &extruderPanel{CommonPanel: NewCommonPanel(ui, parent),
			labels: map[string] * LabelWithImage{},
		}
		m.b = NewBackgroundTask(time.Second * 5, m.updateTemperatures)
		m.initialize()
		extruderPanelInstance = m
	}

	return extruderPanelInstance
}

func (m *extruderPanel) initialize() {
	defer m.Initialize()

	currentRow := 0
	toolheadCount := utils.GetToolheadCount(m.UI.Printer)
	if toolheadCount > 1 {
		m.createToolheadButtons()
		currentRow++
	}

	m.Grid().Attach(m.createExtrudeButton("Extrude", "extruder-extrude.svg",  1), 0, currentRow, 1, 1)
	m.Grid().Attach(m.createExtrudeButton("Retract", "extruder-retract.svg", -1), 3, currentRow, 1, 1)

	m.box = MustBox(gtk.ORIENTATION_VERTICAL, 5)
	m.box.SetVAlign(gtk.ALIGN_CENTER)
	m.box.SetHAlign(gtk.ALIGN_CENTER)
	m.Grid().Attach(m.box, 1, currentRow, 2, 1)
	currentRow++


	m.Grid().Attach(m.createTemperatureButton(), 0, currentRow, 1, 1)

	m.amount = MustStepButton(
		"move-step.svg",
		Step{" 20mm",  20},
		Step{" 50mm",  50},
		Step{"100mm", 100},
		Step{"  1mm",   1},
		Step{"  5mm",   5},
		Step{" 10mm",  10},
	)
	m.Grid().Attach(m.amount, 1, currentRow, 1, 1)

	m.Grid().Attach(m.createFlowRateButton(), 2, currentRow, 1, 1)
}

func (m *extruderPanel) updateTemperatures() {
	s, err := (&octoprint.ToolStateRequest{
		History: true,
		Limit:   1,
	}).Do(m.UI.Printer)

	if err != nil {
		utils.LogError("extruder.updateTemperatures()", "Do(ToolStateRequest)", err)
		return
	}

	m.loadTemperatureState(s)
}

func (m *extruderPanel) loadTemperatureState(s *octoprint.TemperatureState) {
	for tool, current := range s.Current {
		if _, ok := m.labels[tool]; !ok {
			m.addNewTool(tool)
		}

		m.loadTemperatureData(tool, &current)
	}

	m.previous = s
}

func (m *extruderPanel) addNewTool(tool string) {
	utils.Logger.Infof("ExtruderPanel.addNewTool() - new tool detected: %s", tool)

	m.labels[tool] = MustLabelWithImage("toolhead.svg", "")
	m.box.Add(m.labels[tool])
}

func (m *extruderPanel) loadTemperatureData(tool string, d *octoprint.TemperatureData) {
	text := fmt.Sprintf("%.1f°C / %.1f°C", d.Actual, d.Target)

	toolheadCount := utils.GetToolheadCount(m.UI.Printer)
	if toolheadCount > 1 {
		displayNameForTool := utils.GetDisplayNameForTool(tool)
		text = (fmt.Sprintf("%s: ", strings.Title(displayNameForTool)) + text)
	}

	if m.previous != nil && d.Target > 0 {
		if p, ok := m.previous.Current[tool]; ok {
			text = fmt.Sprintf("%s (%.1f°C)", text, d.Actual - p.Actual)
		}
	}

	m.labels[tool].Label.SetText(text)
	m.labels[tool].ShowAll()
}

func (m *extruderPanel) createFlowRateButton() *StepButton {
	b := MustStepButton(
		"speed-fast.svg",
		Step{"Normal (100%)", 100},
		Step{"Fast (125%)", 125},
		Step{"Slow (75%)", 75},
	)

	b.Callback = func() {
		cmd := &octoprint.ToolFlowrateRequest{}
		cmd.Factor = b.Value().(int)

		utils.Logger.Infof("extruder.createFlowRateButton() - changing flowrate to %d%%", cmd.Factor)
		if err := cmd.Do(m.UI.Printer); err != nil {
			utils.LogError("extruder.createFlowRateButton()", "Do(ToolFlowrateRequest)", err)
			return
		}
	}

	return b
}

func (m *extruderPanel) createExtrudeButton(label, image string, dir int) gtk.IWidget {
	return MustPressedButton(label, image, func() {
		cmd := &octoprint.ToolExtrudeRequest{}
		cmd.Amount = m.amount.Value().(int) * dir

		utils.Logger.Infof("extruder.createExtrudeButton() - Sending extrude request with amount %d", cmd.Amount)
		if err := cmd.Do(m.UI.Printer); err != nil {
			utils.LogError("extruder.createFlowRateButton()", "Do(ToolExtrudeRequest)", err)
			return
		}
	}, 200)
}

func (m *extruderPanel) createTemperatureButton() gtk.IWidget {
	return MustButtonImageStyle("Temperature", "heat-up.svg", "color4", func() {
		m.UI.Add(TemperaturePanel(m.UI, m))
	})
}

func (m *extruderPanel) createToolheadButtons() {
	toolheadCount := utils.GetToolheadCount(m.UI.Printer)
	toolheadButtons := CreateChangeToolheadButtonsAndAttachToGrid(toolheadCount, m.Grid())
	m.setToolheadButtonClickHandlers(toolheadButtons)
}

func (m *extruderPanel) setToolheadButtonClickHandlers(toolheadButtons []*gtk.Button) {
	for index, toolheadButton := range toolheadButtons {
		m.setToolheadButtonClickHandler(toolheadButton, index)
	}
}

func (m *extruderPanel) setToolheadButtonClickHandler(toolheadButton *gtk.Button, toolheadIndex int) {
	toolheadButton.Connect("clicked", func() {
		/*
		extruder.go and extruder_multitool.go were consolidated into this single file.  They both sent a command
		to change the tool head (even in extruder.go, which only has one tool head), however they used different
		methods to do so.

		extruder.go used octoprint.ToolSelectRequest (`json:"tool"`)

		...and extruder_multitool.go used m.command() (aka Commonpanel.command())
		octoprint.CommandRequest{}  (`json:"commands"`)

		I don't own a printer with multiple toolheads, so erring on the side of caution, both commands
		are being sent.
		*/

		toolCommand := fmt.Sprintf("tool%d", toolheadIndex)
		utils.Logger.Infof("extruder.setToolheadButtonClickHandler() - Changing tool to %s", toolCommand)

		cmd := &octoprint.ToolSelectRequest{}
		cmd.Tool = toolCommand
		if err := cmd.Do(m.UI.Printer); err != nil {
			utils.LogError("extruder.setToolheadButtonClickHandler()", "Do(ToolSelectRequest)", err)
		}

		gcode := fmt.Sprintf("T%d", toolheadIndex)
		m.command(gcode)
	})
}
