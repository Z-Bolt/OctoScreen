package ui

import (
	"fmt"
	"net"

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

	bar.Add(MustButton(MustImageFromFileWithSize("back.svg", 40, 40), m.UI.GoHistory))

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

	return MustButton(MustImageFromFileWithSize(cmd.Action+".svg", 40, 40), cb)
}

func (m *systemPanel) createInfoBox() gtk.IWidget {
	main := MustBox(gtk.ORIENTATION_HORIZONTAL, 10)
	main.SetHExpand(true)
	main.SetHAlign(gtk.ALIGN_CENTER)
	main.SetVExpand(true)

	img := MustImageFromFileWithSize("logo-white.svg", 140, 112)
	img.SetMarginTop(35)
	main.Add(img)

	info := MustBox(gtk.ORIENTATION_VERTICAL, 0)
	info.SetVExpand(true)
	info.SetVAlign(gtk.ALIGN_CENTER)

	m.addNetwork(info)
	m.addOctoPrint(info)
	m.addSystemInfo(info)

	main.Add(info)

	return main
}

// func (m *systemPanel) addOctoPrintTFT(box *gtk.Box) {
// 	title := MustLabel("<b>OctoPrint-TFT Version</b>")
// 	title.SetMarginTop(15)
// 	title.SetMarginBottom(5)

// 	info := MustBox(gtk.ORIENTATION_VERTICAL, 0)
// 	box.Add(info)

// 	info.Add(title)
// 	info.Add(MustLabel("<b>%s (%s)</b>", Version, Build))
// }

func (m *systemPanel) addNetwork(box *gtk.Box) {
	title := MustLabel("<b>Network Information</b>")
	title.SetMarginTop(40)
	title.SetMarginBottom(5)

	box.Add(title)
	addrs, _ := net.InterfaceAddrs()

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				box.Add(MustLabel("IP Address <b>%s</b>", ipnet.IP.String()))
			}
		}
	}
}

func (m *systemPanel) addOctoPrint(box *gtk.Box) {
	title := MustLabel("<b>Versions Information</b>")
	title.SetMarginTop(15)
	title.SetMarginBottom(5)
	box.Add(title)

	r, err := (&octoprint.VersionRequest{}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return
	}

	box.Add(MustLabel("UI Version: <b>%s (%s)</b>", Version, Build))
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
