package ui

import (
	"fmt"
	"math"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

var nozzleCalibrationPanelInstance *nozzleCalibrationPanel

type pointCoordinates struct {
	x float64
	y float64
	z float64
}

type nozzleCalibrationPanel struct {
	CommonPanel
	zCalibrationMode bool
	activeTool       int
	cPoint           pointCoordinates
	zOffset          float64
	labZOffsetLabel  *gtk.Label
}

func NozzleCalibrationPanel(ui *UI, parent Panel) Panel {
	if nozzleCalibrationPanelInstance == nil {
		m := &nozzleCalibrationPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.panelH = 3
		m.cPoint = pointCoordinates{x: 20, y: 20, z: 0}
		m.initialize()

		nozzleCalibrationPanelInstance = m
	}

	return nozzleCalibrationPanelInstance
}

func (m *nozzleCalibrationPanel) initialize() {
	defer m.Initialize()
	m.Grid().Attach(m.createChangeToolButton(0), 1, 0, 1, 1)
	m.Grid().Attach(m.createChangeToolButton(1), 2, 0, 1, 1)
	m.Grid().Attach(m.createChangeToolButton(2), 3, 0, 1, 1)
	m.Grid().Attach(m.createChangeToolButton(3), 4, 0, 1, 1)

	m.Grid().Attach(m.createIncreaseOffsetButton(), 1, 1, 1, 1)
	m.Grid().Attach(m.createZOffsetLabel(), 2, 1, 2, 1)
	m.Grid().Attach(m.createDecreaseOffsetButton(), 4, 1, 1, 1)

	m.Grid().Attach(m.createZCalibrationModeButton(), 1, 2, 1, 1)
	m.Grid().Attach(m.createAutoZCalibrationButton(), 2, 2, 2, 1)

}

func (m *nozzleCalibrationPanel) createZCalibrationModeButton() gtk.IWidget {
	b := MustStepButton("z-calibration.svg", Step{"Start Manual\nCalibration", false}, Step{"Stop Manual\nCalibration", true})
	ctx, _ := b.GetStyleContext()
	ctx.AddClass("color2")

	b.Callback = func() {
		m.zCalibrationMode = b.Value().(bool)
		if m.zCalibrationMode == true {
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

func (m *nozzleCalibrationPanel) createAutoZCalibrationButton() gtk.IWidget {
	return MustButtonImageStyle("Auto Z Calibration", "z-calibration.svg", "color3", func() {
		if m.zCalibrationMode {
			return
		}

		cmd := &octoprint.RunZOffsetCalibrationRequest{}
		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
		}
	})
}

func (m *nozzleCalibrationPanel) createIncreaseOffsetButton() gtk.IWidget {
	return MustButtonImage("Bed Down", "z-offset-increase.svg", func() {
		if !m.zCalibrationMode {
			return
		}
		m.updateZOffset(m.zOffset + 0.02)
	})
}

func (m *nozzleCalibrationPanel) createDecreaseOffsetButton() gtk.IWidget {
	return MustButtonImage("Bed Up", "z-offset-decrease.svg", func() {
		if !m.zCalibrationMode {
			return
		}
		m.updateZOffset(m.zOffset - 0.02)
	})
}

func (m *nozzleCalibrationPanel) updateZOffset(v float64) {
	m.zOffset = toFixed(v, 4)

	m.labZOffsetLabel.SetText(fmt.Sprintf("Z-Offset: %.2f", m.zOffset))

	cmd := &octoprint.CommandRequest{}
	cmd.Commands = []string{
		fmt.Sprintf("SET_GCODE_OFFSET Z=%f", m.zOffset),
		"G0 Z0 F100",
	}
	if err := cmd.Do(m.UI.Printer); err != nil {
		Logger.Error(err)
	}

	cmd2 := &octoprint.SetZOffsetRequest{Value: m.zOffset, Tool: m.activeTool}
	if err := cmd2.Do(m.UI.Printer); err != nil {
		Logger.Error(err)
	}
}

func (m *nozzleCalibrationPanel) createChangeToolButton(num int) gtk.IWidget {
	style := fmt.Sprintf("color%d", num+1)
	name := fmt.Sprintf("Tool%d", num+1)
	gcode := fmt.Sprintf("T%d", num)
	return MustButtonImageStyle(name, "extruder.svg", style, func() {
		if m.zCalibrationMode {
			m.activeTool = num
			m.command(fmt.Sprintf("G0 Z%f", 5.0))
			m.command(gcode)
			time.Sleep(time.Second * 1)
			m.command(fmt.Sprintf("G0 X%f Y%f F10000", m.cPoint.x, m.cPoint.y))

			cmd := &octoprint.GetZOffsetRequest{Tool: m.activeTool}
			response, err := cmd.Do(m.UI.Printer)

			if err != nil {
				Logger.Error(err)
				return
			}

			m.updateZOffset(response.Offset)

		} else {
			m.command(gcode)
		}
	})
}

func (m *nozzleCalibrationPanel) createZOffsetLabel() gtk.IWidget {
	m.labZOffsetLabel = MustLabel("---")
	m.labZOffsetLabel.SetVAlign(gtk.ALIGN_CENTER)
	m.labZOffsetLabel.SetHAlign(gtk.ALIGN_CENTER)
	m.labZOffsetLabel.SetVExpand(true)
	m.labZOffsetLabel.SetHExpand(true)
	m.labZOffsetLabel.SetLineWrap(true)
	return m.labZOffsetLabel
}

func (m *nozzleCalibrationPanel) command(gcode string) error {
	cmd := &octoprint.CommandRequest{}
	cmd.Commands = []string{gcode}
	return cmd.Do(m.UI.Printer)
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
