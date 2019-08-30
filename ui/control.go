package ui

import (
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
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

var controlPanelInstance *controlPanel

type controlPanel struct {
	CommonPanel
}

func ControlPanel(ui *UI, parent Panel) Panel {
	if controlPanelInstance == nil {
		m := &controlPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.initialize()
		controlPanelInstance = m
	}

	return controlPanelInstance
}

func (m *controlPanel) initialize() {
	defer m.Initialize()

	for _, c := range control {
		icon := strings.ToLower(strings.Replace(c.Name, " ", "-", -1))
		b := m.createControlButton(c, icon)
		m.AddButton(b)
	}

	for _, c := range m.getCustomControl() {
		b := m.createControlButton(c, "custom-script")
		m.AddButton(b)
	}

	for _, c := range m.getCommands() {
		b := m.createCommandButton(c)
		m.AddButton(b)
	}
}

func (m *controlPanel) getCustomControl() []*octoprint.ControlDefinition {

	control := []*octoprint.ControlDefinition{}

	Logger.Info("Retrieving custom controls")
	r, err := (&octoprint.CustomCommandsRequest{}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return control
	}

	for _, c := range r.Controls {
		for _, cc := range c.Children {
			if cc.Command != "" || cc.Script != "" {
				control = append(control, cc)
			}
		}
	}

	return control
}

func (m *controlPanel) createControlButton(c *octoprint.ControlDefinition, icon string) gtk.IWidget {

	do := func() {
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
	}

	cb := do
	if len(c.Confirm) != 0 {
		cb = MustConfirmDialog(m.UI.w, c.Confirm, do)
	}

	return MustButtonImage(strEllipsisLen(c.Name, 25), icon+".svg", cb)
}

func (m *controlPanel) createCommandButton(c *octoprint.CommandDefinition) gtk.IWidget {
	do := func() {
		r := &octoprint.SystemExecuteCommandRequest{
			Source: octoprint.Custom,
			Action: c.Action,
		}

		if err := r.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	}

	cb := do
	if len(c.Confirm) != 0 {
		cb = MustConfirmDialog(m.UI.w, c.Confirm, do)
	}

	return MustButtonImage(strEllipsisLen(c.Name, 25), c.Action+".svg", cb)
}

func (m *controlPanel) getCommands() []*octoprint.CommandDefinition {
	Logger.Info("Retrieving custom commands")
	r, err := (&octoprint.SystemCommandsRequest{}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return nil
	}

	return r.Custom
}
