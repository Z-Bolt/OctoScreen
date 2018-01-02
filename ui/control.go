package ui

import (
	"strings"

	"github.com/gotk3/gotk3/gtk"
	octoprint "github.com/mcuadros/go-octoprint"
)

var control = []*octoprint.ControlDefinition{{
	Name:    "Motor Off",
	Command: "M18",
}, {
	Name:    "Fan On",
	Command: "M106",
}, {
	Name:    "Fan Off",
	Command: "M106 S0",
}}

type ControlPanel struct {
	CommonPanel
}

func NewControlPanel(ui *UI) *ControlPanel {
	m := &ControlPanel{CommonPanel: NewCommonPanel(ui)}
	m.initialize()
	return m
}

func (m *ControlPanel) initialize() {
	m.Initialize()

	for _, c := range m.getControl() {
		b := m.createControlButton(c)
		m.AddButton(b)
	}
}

func (m *ControlPanel) getControl() []*octoprint.ControlDefinition {
	control := control

	Logger.Info("Retrieving custom commands")
	r, err := (&octoprint.CustomCommandsRequest{}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return control
	}

	for _, c := range r.Controls {
		control = append(control, c.Children...)
	}

	return control
}

func (m *ControlPanel) createControlButton(c *octoprint.ControlDefinition) gtk.IWidget {
	icon := strings.ToLower(strings.Replace(c.Name, " ", "-", -1))
	return MustButtonImage(c.Name, icon+".svg", func() {
		r := &octoprint.CommandRequest{
			Commands: c.Commands,
		}

		if len(c.Command) != 0 {
			r.Commands = []string{c.Command}
		}

		Logger.Infof("Executing command %q", c.Name)
		if err := r.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})
}
