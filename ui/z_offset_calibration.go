package ui

import (
	"fmt"
	"math"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/utils"
)

var zOffsetCalibrationPanelInstance *zOffsetCalibrationPanel

type pointCoordinates struct {
	x float64
	y float64
	z float64
}

type zOffsetCalibrationPanel struct {
	CommonPanel
	zCalibrationMode bool
	activeTool       int
	cPoint           pointCoordinates
	zOffset          float64
	labZOffsetLabel  *gtk.Label
}

func ZOffsetCalibrationPanel(ui *UI, parent Panel) Panel {
	if zOffsetCalibrationPanelInstance == nil {
		m := &zOffsetCalibrationPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.cPoint = pointCoordinates{x: 20, y: 20, z: 0}
		m.initialize()

		zOffsetCalibrationPanelInstance = m
	}

	return zOffsetCalibrationPanelInstance
}

func (m *zOffsetCalibrationPanel) initialize() {
	defer m.Initialize()

	currentRow := 0
	toolheadCount := utils.GetToolheadCount(m.UI.Printer)
	if toolheadCount > 1 {
		m.createToolheadButtons()
		currentRow++
	}

	m.Grid().Attach(m.createIncreaseOffsetButton(), 0, currentRow, 1, 1)
	m.Grid().Attach(m.createZOffsetLabel(),         1, currentRow, 2, 1)
	m.Grid().Attach(m.createDecreaseOffsetButton(), 3, currentRow, 1, 1)
	currentRow++

	m.Grid().Attach(m.createZCalibrationModeButton(), 0, currentRow, 1, 1)
	m.Grid().Attach(m.createAutoZCalibrationButton(), 1, currentRow, 2, 1)
}

func (m *zOffsetCalibrationPanel) createZCalibrationModeButton() gtk.IWidget {
	b := MustStepButton("z-calibration.svg", Step{"Start Manual\nCalibration", false}, Step{"Stop Manual\nCalibration", true})
	ctx, _ := b.GetStyleContext()

	b.Callback = func() {
		m.zCalibrationMode = b.Value().(bool)
		if m.zCalibrationMode {
			ctx.AddClass("active")

			m.command("G28")
			m.command("T0")
			time.Sleep(time.Second * 1)
			m.command(fmt.Sprintf("G0 X%f Y%f F10000", m.cPoint.x, m.cPoint.y))
			m.command(fmt.Sprintf("G0 Z10 F2000"))
			m.command(fmt.Sprintf("G0 Z%f F400", m.cPoint.z))

			m.activeTool = 0
			m.updateZOffset(0)
		} else {
			ctx.RemoveClass("active")
			m.labZOffsetLabel.SetText("Press \"Z Offset\"\nbutton to start\nZ-Offset calibration")
		}
	}

	return b
}

func (m *zOffsetCalibrationPanel) createAutoZCalibrationButton() gtk.IWidget {
	return MustButtonImageStyle("Auto Z Calibration", "z-calibration.svg", "color3", func() {
		if m.zCalibrationMode {
			return
		}

		cmd := &octoprint.RunZOffsetCalibrationRequest{}
		if err := cmd.Do(m.UI.Printer); err != nil {
			utils.LogError("z_offset_calibration.createAutoZCalibrationButton()", "Do(RunZOffsetCalibrationRequest)", err)
		}
	})
}

func (m *zOffsetCalibrationPanel) createIncreaseOffsetButton() gtk.IWidget {
	return MustButtonImage("Bed Down", "z-offset-increase.svg", func() {
		if !m.zCalibrationMode {
			return
		}
		m.updateZOffset(m.zOffset + 0.02)
	})
}

func (m *zOffsetCalibrationPanel) createDecreaseOffsetButton() gtk.IWidget {
	return MustButtonImage("Bed Up", "z-offset-decrease.svg", func() {
		if !m.zCalibrationMode {
			return
		}
		m.updateZOffset(m.zOffset - 0.02)
	})
}

func (m *zOffsetCalibrationPanel) updateZOffset(v float64) {
	m.zOffset = toFixed(v, 4)

	m.labZOffsetLabel.SetText(fmt.Sprintf("Z-Offset: %.2f", m.zOffset))

	cmd := &octoprint.CommandRequest{}
	cmd.Commands = []string{
		fmt.Sprintf("SET_GCODE_OFFSET Z=%f", m.zOffset),
		"G0 Z0 F100",
	}
	if err := cmd.Do(m.UI.Printer); err != nil {
		utils.LogError("z_offset_calibration.updateZOffset()", "Do(CommandRequest)", err)
	}

	cmd2 := &octoprint.SetZOffsetRequest{Value: m.zOffset, Tool: m.activeTool}
	if err := cmd2.Do(m.UI.Printer); err != nil {
		utils.LogError("z_offset_calibration.updateZOffset()", "Do(SetZOffsetRequest)", err)
	}
}

func (m *zOffsetCalibrationPanel) createZOffsetLabel() gtk.IWidget {
	m.labZOffsetLabel = MustLabel("---")
	m.labZOffsetLabel.SetVAlign(gtk.ALIGN_CENTER)
	m.labZOffsetLabel.SetHAlign(gtk.ALIGN_CENTER)
	m.labZOffsetLabel.SetVExpand(true)
	m.labZOffsetLabel.SetHExpand(true)
	m.labZOffsetLabel.SetLineWrap(true)

	return m.labZOffsetLabel
}

func (m *zOffsetCalibrationPanel) command(gcode string) error {
	cmd := &octoprint.CommandRequest{}
	cmd.Commands = []string{gcode}

	return cmd.Do(m.UI.Printer)
}


func (m *zOffsetCalibrationPanel) createToolheadButtons() {
	toolheadCount := utils.GetToolheadCount(m.UI.Printer)
	toolheadButtons := CreateChangeToolheadButtonsAndAttachToGrid(toolheadCount, m.Grid())
	m.setToolheadButtonClickHandlers(toolheadButtons)
}

func (m *zOffsetCalibrationPanel) setToolheadButtonClickHandlers(toolheadButtons []*gtk.Button) {
	for index, toolheadButton := range toolheadButtons {
		m.setToolheadButtonClickHandler(toolheadButton, index)
	}
}

func (m *zOffsetCalibrationPanel) setToolheadButtonClickHandler(toolheadButton *gtk.Button, toolheadIndex int) {
	toolheadButton.Connect("clicked", func() {
		utils.Logger.Infof("Changing tool to tool%d", toolheadIndex)

		gcode := fmt.Sprintf("T%d", toolheadIndex)

		if m.zCalibrationMode {
			m.activeTool = toolheadIndex
			m.command(fmt.Sprintf("G0 Z%f", 5.0))
			m.command(gcode)
			time.Sleep(time.Second * 1)
			m.command(fmt.Sprintf("G0 X%f Y%f F10000", m.cPoint.x, m.cPoint.y))

			cmd := &octoprint.GetZOffsetRequest{Tool: m.activeTool}
			response, err := cmd.Do(m.UI.Printer)

			if err != nil {
				utils.LogError("z_offset_calibration.setToolheadButtonClickHandler()", "Do(GetZOffsetRequest)", err)
				return
			}

			m.updateZOffset(response.Offset)
		} else {
			m.command(gcode)
		}
	})
}



// TODO: place these function in a util file
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
