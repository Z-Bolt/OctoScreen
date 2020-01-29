package ui

import (
	"fmt"
	"time"

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
	psuContol *PSUControl
}

func SystemPanel(ui *UI, parent Panel) *systemPanel {
	if systemPanelInstance == nil {
		m := &systemPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.panelH = 3
		m.initialize()
		systemPanelInstance = m
	} else {
		systemPanelInstance.p = parent
	}

	return systemPanelInstance
}

func (m *systemPanel) initialize() {
	defer m.Initialize()

	m.Grid().Attach(m.createOctoPrintInfo(), 1, 0, 2, 1)
	m.Grid().Attach(m.createOctoScreenInfo(), 3, 0, 2, 1)
	m.Grid().Attach(m.createSystemInfo(), 1, 1, 3, 1)
	m.psuContol = m.addPsuButton()
	m.psuContol.update()
	m.Grid().Attach(m.psuContol.Button, 4, 1, 1, 1)

	if b := m.createCommandButton("Octo Restart", "restart", "color2"); b != nil {
		m.Grid().Attach(b, 3, 2, 1, 1)
	}

	if b := m.createCommandButton("Sys Restart", "reboot", "color3"); b != nil {
		m.Grid().Attach(b, 2, 2, 1, 1)
	}

	if b := m.createCommandButton("Shutdown", "shutdown", "color1"); b != nil {
		m.Grid().Attach(b, 1, 2, 1, 1)
	}
	m.b = NewBackgroundTask(time.Second*10, m.psuContol.update)
}

func (m *systemPanel) createOctoPrintInfo() *gtk.Box {
	r, err := (&octoprint.VersionRequest{}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return nil
	}

	info := MustBox(gtk.ORIENTATION_VERTICAL, 0)

	info.SetHExpand(true)
	info.SetHAlign(gtk.ALIGN_CENTER)
	info.SetVExpand(true)
	info.SetVAlign(gtk.ALIGN_CENTER)
	logoWidth := m.Scaled(52)
	img := MustImageFromFileWithSize("logo-octoprint.png", logoWidth, int(float64(logoWidth)*1.25))
	info.Add(img)

	info.Add(MustLabel("\nOctoPrint Version"))
	info.Add(MustLabel("<b>%s (%s)</b>", r.Server, r.API))
	return info
}

func (m *systemPanel) createOctoScreenInfo() *gtk.Box {
	info := MustBox(gtk.ORIENTATION_VERTICAL, 0)

	info.SetHExpand(true)
	info.SetHAlign(gtk.ALIGN_CENTER)
	info.SetVExpand(true)
	info.SetVAlign(gtk.ALIGN_CENTER)

	logoWidth := m.Scaled(62)

	img := MustImageFromFileWithSize("logo-z-bolt.svg", logoWidth, int(float64(logoWidth)*0.8))
	info.Add(img)
	info.Add(MustLabel("OctoScreen Version"))
	info.Add(MustLabel("<b>%s (%s)</b>", Version, Build))
	return info
}

func (m *systemPanel) addPsuButton() *PSUControl {
	return PSUControlNew(m.UI, m.UI.Printer)
}

func (m *systemPanel) createSystemInfo() *gtk.Box {
	info := MustBox(gtk.ORIENTATION_VERTICAL, 0)

	info.SetVExpand(true)
	info.SetVAlign(gtk.ALIGN_CENTER)

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

	return info
}

func (m *systemPanel) createCommandButton(name string, action string, style string) gtk.IWidget {
	r, err := (&octoprint.SystemCommandsRequest{}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return nil
	}

	var cmd *octoprint.CommandDefinition
	var cb func()

	for _, c := range r.Core {
		if c.Action == action {
			cmd = c
		}
	}

	if cmd != nil {
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

		cb = do

		if len(cmd.Confirm) != 0 {
			cb = MustConfirmDialog(m.UI.w, cmd.Confirm, do)
		}
	}

	b := MustButtonImageStyle(name, action+".svg", style, cb)

	if cmd == nil {
		b.SetSensitive(false)
	}

	return b
}
