package ui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

var toolchangerPanelInstance *toolchangerPanel

type pointCoordinates struct {
	x float64
	y float64
	z float64
}

type toolchangerPanel struct {
	CommonPanel
	zCalibrationMode bool
	cPoint           pointCoordinates
	zOffset          float64
	labZOffsetLabel  *gtk.Label
}

func ToolchangerPanel(ui *UI, parent Panel) Panel {
	if toolchangerPanelInstance == nil {
		m := &toolchangerPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.panelH = 3
		m.cPoint = pointCoordinates{x: 20, y: 20, z: 0}
		// m.b = NewBackgroundTask(time.Second, m.updateTemperatures)
		m.initialize()

		toolchangerPanelInstance = m
	}

	return toolchangerPanelInstance
}

func (m *toolchangerPanel) initialize() {
	defer m.Initialize()
	m.Grid().Attach(m.createChangeToolButton(0), 1, 0, 1, 1)
	m.Grid().Attach(m.createChangeToolButton(1), 2, 0, 1, 1)
	m.Grid().Attach(m.createChangeToolButton(2), 3, 0, 1, 1)
	m.Grid().Attach(m.createChangeToolButton(3), 4, 0, 1, 1)

	m.Grid().Attach(m.createHomeButton(), 1, 1, 1, 1)
	m.Grid().Attach(m.createIncreaseOffsetButton(), 2, 1, 1, 1)
	m.Grid().Attach(m.createZOffsetLabel(), 3, 1, 1, 1)
	m.Grid().Attach(m.createDecreaseOffsetButton(), 4, 1, 1, 1)

	m.Grid().Attach(m.createMagnetOnButton(), 1, 2, 1, 1)
	m.Grid().Attach(m.createMagnetOffButton(), 2, 2, 1, 1)
	m.Grid().Attach(m.createZCalibrationModeButton(), 3, 2, 1, 1)

}

func (m *toolchangerPanel) createZCalibrationModeButton() gtk.IWidget {
	b := MustStepButton("z-calibration.svg", Step{"Z Offset", false}, Step{"Z Offset", true})
	ctx, _ := b.GetStyleContext()
	ctx.AddClass("color2")

	b.Callback = func() {
		m.zCalibrationMode = b.Value().(bool)
		if m.zCalibrationMode == true {
			ctx.AddClass("active")
			cmd := &octoprint.CommandRequest{}
			cmd.Commands = []string{
				"G28",
				fmt.Sprintf("G0 X%f Y%f", m.cPoint.x, m.cPoint.y),
				fmt.Sprintf("G0 Z%f", m.cPoint.z),
			}

			if err := cmd.Do(m.UI.Printer); err != nil {
				Logger.Error(err)
			}

			m.updateZOffset(0)
		} else {
			ctx.RemoveClass("active")
			m.labZOffsetLabel.SetText("Press \"Z Offset\" button to switch to Z-Offset calibration mode.")
		}
	}

	return b
}

func (m *toolchangerPanel) createHomeButton() gtk.IWidget {
	return MustButtonImageStyle("Home XY", "home.svg", "color3", func() {
		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{
			"G28 X",
			"G28 Y",
		}
		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
		}
	})
}

func (m *toolchangerPanel) createIncreaseOffsetButton() gtk.IWidget {
	return MustButtonImage("Bed Down", "z-offset-increase.svg", func() {
		if !m.zCalibrationMode {
			return
		}
		m.updateZOffset(m.zOffset + 0.02)
	})
}

func (m *toolchangerPanel) createDecreaseOffsetButton() gtk.IWidget {
	return MustButtonImage("Bed Up", "z-offset-decrease.svg", func() {
		if !m.zCalibrationMode {
			return
		}
		m.updateZOffset(m.zOffset - 0.02)
	})
}

func (m *toolchangerPanel) updateZOffset(v float64) {
	m.zOffset = v

	m.labZOffsetLabel.SetText(fmt.Sprintf("Z-Offset: %.2f", m.zOffset))

	cmd := &octoprint.CommandRequest{}
	cmd.Commands = []string{
		fmt.Sprintf("SET_GCODE_OFFSET Z=%f", m.zOffset),
		"G0 Z0",
	}
	if err := cmd.Do(m.UI.Printer); err != nil {
		Logger.Error(err)
	}
}

func (m *toolchangerPanel) createChangeToolButton(num int) gtk.IWidget {
	style := fmt.Sprintf("color%d", num+1)
	name := fmt.Sprintf("Tool%d", num+1)
	gcode := fmt.Sprintf("T%d", num)
	return MustButtonImageStyle(name, "extruder.svg", style, func() {
		cmd := &octoprint.CommandRequest{}

		if m.zCalibrationMode {
			cmd.Commands = []string{
				fmt.Sprintf("G0 Z%f", 5.0),
				gcode,
				fmt.Sprintf("G0 X%f Y%f", m.cPoint.x, m.cPoint.y),
			}
		} else {
			cmd.Commands = []string{gcode}
		}

		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}

		if m.zCalibrationMode {
			m.updateZOffset(0)
		}
	})
}

func (m *toolchangerPanel) createMagnetOnButton() gtk.IWidget {
	return MustButtonImageStyle("Magnet On", "magnet-on.svg", "color4", func() {
		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{"SET_PIN PIN=sol VALUE=1"}

		Logger.Info("Turn on magnet")
		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})
}

func (m *toolchangerPanel) createMagnetOffButton() gtk.IWidget {
	return MustButtonImageStyle("Magnet Off", "magnet-off.svg", "color3", func() {
		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{"SET_PIN PIN=sol VALUE=0"}

		Logger.Info("Turn off magnet")
		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})
}

func (m *toolchangerPanel) createZOffsetLabel() gtk.IWidget {
	m.labZOffsetLabel = MustLabel("Press \"Z Offset\" button to switch to Z-Offset calibration mode.")
	m.labZOffsetLabel.SetVAlign(gtk.ALIGN_CENTER)
	m.labZOffsetLabel.SetHAlign(gtk.ALIGN_CENTER)
	m.labZOffsetLabel.SetVExpand(true)
	m.labZOffsetLabel.SetHExpand(true)
	m.labZOffsetLabel.SetLineWrap(true)
	return m.labZOffsetLabel
}
