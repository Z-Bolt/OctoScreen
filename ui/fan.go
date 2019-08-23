package ui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

var fanPanelInstance *fanPanel

type fanPanel struct {
	CommonPanel
	step *StepButton
}

func FanPanel(ui *UI, parent Panel) Panel {
	if fanPanelInstance == nil {
		m := &fanPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.panelH = 2
		m.panelW = 3
		m.initialize()
		fanPanelInstance = m
	}

	return fanPanelInstance
}

func (m *fanPanel) initialize() {
	defer m.Initialize()

	m.Grid().Attach(m.createFanButton(100), 0, 0, 1, 1)
	m.Grid().Attach(m.createFanButton(75), 1, 0, 1, 1)
	m.Grid().Attach(m.createFanButton(50), 2, 0, 1, 1)
	m.Grid().Attach(m.createFanButton(25), 3, 0, 1, 1)

	m.Grid().Attach(m.createFanButton(0), 0, 1, 1, 1)
}

func (m *fanPanel) createFanButton(speed int) gtk.IWidget {

	var (
		label string
		image string
		color string
	)

	if speed == 0 {
		label = "Fan Off"
		image = "fan-off.svg"
		color = "color2"
	} else {
		label = fmt.Sprintf("%d %%", speed)
		image = "fan.svg"
		color = "color4"
	}

	return MustButtonImageStyle(label, image, color, func() {
		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{
			fmt.Sprintf("M106 S%d", (255 * speed / 100)),
		}
		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
		}
	})
}
