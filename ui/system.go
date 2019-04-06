package ui

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/dustin/go-humanize"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

var systemPanelInstance *systemPanel

type systemPanel struct {
	CommonPanel

	list *gtk.Box
}

func SystemPanel(ui *UI, parent Panel) *systemPanel {
	if systemPanelInstance == nil {
		m := &systemPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.initialize()
		systemPanelInstance = m
	}

	return systemPanelInstance
}

func (m *systemPanel) initialize() {
	box := MustBox(gtk.ORIENTATION_VERTICAL, 0)
	box.Add(m.createInfoBox())
	box.Add(m.createActionBar())
	m.Grid().Add(box)
}

func (m *systemPanel) createActionBar() gtk.IWidget {
	bar := MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	bar.SetHAlign(gtk.ALIGN_END)
	bar.SetHExpand(true)
	bar.SetMarginTop(5)
	bar.SetMarginBottom(5)
	bar.SetMarginEnd(5)

	if b := m.createRestartButton(); b != nil {
		bar.Add(b)
	}

	if c := m.createShutdownButton(); c != nil {
		bar.Add(c)
	}


	bar.Add(MustButton(MustImageFromFileWithSize("back.svg", 60, 60), m.UI.GoHistory))

	return bar
}

func (m *systemPanel) createRestartButton() gtk.IWidget {
	r, err := (&octoprint.SystemCommandsRequest{}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return nil
	}

	var cmd *octoprint.CommandDefinition
	for _, c := range r.Core {
		if c.Action == "reboot" {
			cmd = c
		}
	}

	if cmd == nil {
		return nil
	}

	return m.doCreateButtonFromCommand(cmd)
}

func (m *systemPanel) createShutdownButton() gtk.IWidget {
	r, err := (&octoprint.SystemCommandsRequest{}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return nil
	}

	var cmd *octoprint.CommandDefinition
	for _, c := range r.Core {
		if c.Action == "shutdown" {
			cmd = c
		}
	}

	if cmd == nil {
		return nil
	}

	return m.doCreateButtonFromCommand(cmd)
}


func (m *systemPanel) doCreateButtonFromCommand(cmd *octoprint.CommandDefinition) gtk.IWidget {
	do := func() {
		r := &octoprint.SystemExecuteCommandRequest{
			Source: octoprint.Core,
			Action: cmd.Action,
		}

		if err := r.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	}

	cb := do
	if len(cmd.Confirm) != 0 {
		cb = MustConfirmDialog(m.UI.w, cmd.Confirm, do)
	}

	return MustButton(MustImageFromFileWithSize(cmd.Action+".svg", 60, 60), cb)
}

func (m *systemPanel) createInfoBox() gtk.IWidget {
	main := MustBox(gtk.ORIENTATION_HORIZONTAL, 10)
	main.SetHExpand(true)
	main.SetHAlign(gtk.ALIGN_CENTER)
	main.SetVExpand(true)
	// main.Add(MustImageFromFileWithSize("octoprint-logo.png", 140, 140))

	info := MustBox(gtk.ORIENTATION_VERTICAL, 0)
	info.SetVExpand(true)
	info.SetVAlign(gtk.ALIGN_CENTER)
	m.addOctoPrintTFT(info)

	title := MustLabel("<b>Versions Information</b>")
	title.SetMarginTop(15)
	title.SetMarginBottom(5)
	info.Add(title)

	m.addOctoPrint(info)
	m.addOctoPi(info)
	m.addSystemInfo(info)

	main.Add(info)

	return main
}

func (m *systemPanel) addOctoPrintTFT(box *gtk.Box) {
	title := MustLabel("<b>OctoPrint-TFT Version</b>")
	title.SetMarginBottom(5)

	info := MustBox(gtk.ORIENTATION_VERTICAL, 0)
	box.Add(info)

	info.Add(title)
	info.Add(MustLabel("<b>%s (%s)</b>", Version, Build))
}

func (m *systemPanel) addOctoPi(box *gtk.Box) {
	v, err := ioutil.ReadFile("/etc/octopi_version")
	if err != nil {
		Logger.Error(err)
		return
	}

	box.Add(MustLabel("OctoPi Version: <b>%s</b>", bytes.Trim(v, "\n")))
}

func (m *systemPanel) addOctoPrint(box *gtk.Box) {
	r, err := (&octoprint.VersionRequest{}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return
	}

	box.Add(MustLabel("OctoPrint Version: <b>%s (%s)</b>", r.Server, r.API))
}

func (m *systemPanel) addSystemInfo(box *gtk.Box) {
	info := MustBox(gtk.ORIENTATION_VERTICAL, 0)
	box.Add(info)

	title := MustLabel("<b>System Information</b>")
	title.SetMarginBottom(5)
	title.SetMarginTop(15)
	info.Add(title)

	v, _ := mem.VirtualMemory()
	info.Add(MustLabel(fmt.Sprintf(
		"Memory Total / Free: <b>%s / %s</b>",
		humanize.Bytes(v.Total), humanize.Bytes(v.Free),
	)))

	l, _ := load.Avg()
	info.Add(MustLabel(fmt.Sprintf(
		"Load Average: <b>%.2f, %.2f, %.2f</b>",
		l.Load1, l.Load5, l.Load15,
	)))
}
