package ui

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/dustin/go-humanize"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

type SystemPanel struct {
	CommonPanel

	list *gtk.Box
}

func NewSystemPanel(ui *UI) *SystemPanel {
	m := &SystemPanel{CommonPanel: NewCommonPanel(ui)}
	m.initialize()
	return m
}

func (m *SystemPanel) initialize() {
	box := MustBox(gtk.ORIENTATION_VERTICAL, 0)
	box.Add(m.createInfoBox())
	box.Add(m.createActionBar())
	m.grid.Add(box)
}

func (m *SystemPanel) createActionBar() gtk.IWidget {
	bar := MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	bar.SetHAlign(gtk.ALIGN_END)
	bar.SetHExpand(true)
	bar.SetMarginTop(5)
	bar.SetMarginBottom(5)
	bar.SetMarginEnd(5)
	bar.Add(MustButton(MustImageFromFileWithSize("power.svg", 40, 40), m.UI.ShowDefaultPanel))
	bar.Add(MustButton(MustImageFromFileWithSize("back.svg", 40, 40), m.UI.ShowDefaultPanel))

	return bar
}

func (m *SystemPanel) createInfoBox() gtk.IWidget {
	main := MustBox(gtk.ORIENTATION_HORIZONTAL, 20)
	main.SetHExpand(true)
	main.SetHAlign(gtk.ALIGN_CENTER)
	main.SetVExpand(true)
	main.Add(MustImageFromFileWithSize("octoprint-logo.png", 140, 140))

	info := MustBox(gtk.ORIENTATION_VERTICAL, 0)
	info.SetVExpand(true)
	info.SetVAlign(gtk.ALIGN_CENTER)

	title := MustLabel("<b>Versions Information</b>")
	title.SetMarginBottom(5)
	info.Add(title)

	m.addOctoPrintTFT(info)
	m.addOctoPrint(info)
	m.addOctoPi(info)
	m.addKernel(info)
	m.addSystemInfo(info)

	main.Add(info)

	return main
}

func (m *SystemPanel) addOctoPrintTFT(box *gtk.Box) {
	box.Add(MustLabel("OctoPrint-TFT Version: <b>%s (%s)</b>", Version, Commit[:7]))
}

func (m *SystemPanel) addOctoPi(box *gtk.Box) {
	v, err := ioutil.ReadFile("/etc/octopi_version")
	if err != nil {
		Logger.Error(err)
		return
	}

	box.Add(MustLabel("OctoPi Version: <b>%s</b>", bytes.Trim(v, "\n")))
}

func (m *SystemPanel) addOctoPrint(box *gtk.Box) {
	r, err := (&octoprint.VersionRequest{}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return
	}

	box.Add(MustLabel("OctoPrint Version: <b>%s (%s)</b>", r.Server, r.API))
}

func (m *SystemPanel) addKernel(box *gtk.Box) {
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		log.Fatal(err)
	}

	box.Add(MustLabel("Kernel Version: <b>%s</b>", out))
}

func (m *SystemPanel) addSystemInfo(box *gtk.Box) {
	info := MustBox(gtk.ORIENTATION_VERTICAL, 0)
	box.Add(info)

	title := MustLabel("<b>System Information</b>")
	title.SetMarginBottom(5)
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
