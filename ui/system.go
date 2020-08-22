package ui

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/Z-Bolt/OctoScreen/utils"
)

var systemPanelInstance *systemPanel = nil

type systemPanel struct {
	CommonPanel

	//list *gtk.Box
}

func SystemPanel(ui *UI, parent Panel) *systemPanel {
	if systemPanelInstance == nil {
		m := &systemPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.initialize()
		systemPanelInstance = m
	} else {
		systemPanelInstance.p = parent
	}

	return systemPanelInstance
}

func (m *systemPanel) initialize() {
	defer m.Initialize()

	// First row
	m.Grid().Attach(m.createOctoPrintInfo(),        0, 0, 1, 1)
	m.Grid().Attach(m.createOctoScreenInfo(),       1, 0, 2, 1)
	m.Grid().Attach(m.createOctoScreenPluginInfo(), 3, 0, 1, 1)

	// Second row
	m.Grid().Attach(m.createSystemInfo(),           0, 1, 4, 1)

	// Third row
	if b := m.createCommandButton("Octo Restart", "restart", "color-warning-sign-yellow"); b != nil {
		m.Grid().Attach(b, 2, 2, 1, 1)
	}

	if b := m.createCommandButton("Sys Restart", "reboot", "color-warning-sign-yellow"); b != nil {
		m.Grid().Attach(b, 1, 2, 1, 1)
	}

	if b := m.createCommandButton("Shutdown", "shutdown", "color-warning-sign-yellow"); b != nil {
		m.Grid().Attach(b, 0, 2, 1, 1)
	}
}

func (m *systemPanel) createOctoPrintInfo() *gtk.Box {
	r, err := (&octoprint.VersionRequest{}).Do(m.UI.Printer)
	if err != nil {
		utils.LogError("system.createOctoPrintInfo()", "Do(VersionRequest)", err)
		return nil
	}

	infoBox := MustBox(gtk.ORIENTATION_VERTICAL, 0)
	infoBox.SetHExpand(true)
	infoBox.SetHAlign(gtk.ALIGN_CENTER)
	infoBox.SetVExpand(true)
	infoBox.SetVAlign(gtk.ALIGN_CENTER)

	logoWidth := m.Scaled(52)
	logoHeight := int(float64(logoWidth) * 1.25)
	logoImage := MustImageFromFileWithSize("logo-octoprint.png", logoWidth, logoHeight)

	infoBox.Add(logoImage)

	infoBox.Add(MustLabel(""))
	infoBox.Add(MustLabel("OctoPrint"))

	// Can't display on two lines - it's too tall and causes a window resize.
	// infoBox.Add(MustLabel("<b>%s</b>", r.Server))
	// infoBox.Add(MustLabel("<b>(API %s)</b>", r.API))

	// Just display on one line and hope that the version string back from
	// OctoPrint isn't too large.
	infoBox.Add(MustLabel("<b>%s (API %s)</b>", r.Server, r.API))

	return infoBox
}

func (m *systemPanel) createOctoScreenInfo() *gtk.Box {
	infoBox := MustBox(gtk.ORIENTATION_VERTICAL, 0)
	infoBox.SetHExpand(true)
	infoBox.SetHAlign(gtk.ALIGN_CENTER)
	infoBox.SetVExpand(true)
	infoBox.SetVAlign(gtk.ALIGN_CENTER)

	logoImage := MustImageFromFile("octoScreen-isometric.png")

	infoBox.Add(logoImage)
	infoBox.Add(MustLabel("OctoScreen"))
	infoBox.Add(MustLabel("<b>%s</b>", OctoScreenVersion))

	return infoBox
}

func (m *systemPanel) createOctoScreenPluginInfo() *gtk.Box {
	infoBox := MustBox(gtk.ORIENTATION_VERTICAL, 0)
	infoBox.SetHExpand(true)
	infoBox.SetHAlign(gtk.ALIGN_CENTER)
	infoBox.SetVExpand(true)
	infoBox.SetVAlign(gtk.ALIGN_CENTER)

	logoImage := MustImageFromFile("puzzle-piece.png")
	infoBox.Add(logoImage)

	infoBox.Add(MustLabel(""))
	infoBox.Add(MustLabel("OctoScreen plugin"))

	if m.UI.OctoPrintPlugin {
		getPluginManagerInfoResponse, err := (&octoprint.GetPluginManagerInfoRequest{}).Do(m.UI.Printer)
		if err != nil {
			utils.LogError("system.createOctoScreenPluginInfo()", "Do(GetPluginManagerInfoRequest)", err)
			return infoBox
		}

		found := false
		for i := 0; i < len(getPluginManagerInfoResponse.Plugins) && !found; i++ {
			plugin := getPluginManagerInfoResponse.Plugins[i]
			if plugin.Key == "zbolt_octoscreen" {
				found = true
				infoBox.Add(MustLabel("<b>%s</b>", plugin.Version))
			}
		}

		if !found {
			// OK, the plugin is there, we just can't get the info from a GET request.
			// Default to displaying, "Present"
			infoBox.Add(MustLabel("<b>%s</b>", "Present"))
		}
	} else {
		infoBox.Add(MustLabel("<b>%s</b>", "Not installed"))
	}

	return infoBox
}

func (m *systemPanel) createSystemInfo() *gtk.Box {
	infoBox := MustBox(gtk.ORIENTATION_VERTICAL, 0)
	infoBox.SetVExpand(true)
	infoBox.SetVAlign(gtk.ALIGN_CENTER)

	title := MustLabel("<b>System Information</b>")
	title.SetMarginBottom(5)
	title.SetMarginTop(15)
	infoBox.Add(title)

	v, _ := mem.VirtualMemory()
	infoBox.Add(MustLabel(fmt.Sprintf(
		"Memory Total / Free: <b>%s / %s</b>",
		humanize.Bytes(v.Total), humanize.Bytes(v.Free),
	)))

	l, _ := load.Avg()
	infoBox.Add(MustLabel(fmt.Sprintf(
		"Load Average: <b>%.2f, %.2f, %.2f</b>",
		l.Load1, l.Load5, l.Load15,
	)))

	return infoBox
}

func (m *systemPanel) createCommandButton(name string, action string, style string) gtk.IWidget {
	systemCommandsResponse, err := (&octoprint.SystemCommandsRequest{}).Do(m.UI.Printer)
	if err != nil {
		utils.LogError("system.createCommandButton()", "Do(SystemCommandsRequest)", err)
		return nil
	}

	var cmd *octoprint.CommandDefinition
	var cb func()

	for _, commandDefinition := range systemCommandsResponse.Core {
		if commandDefinition.Action == action {
			cmd = commandDefinition
		}
	}

	if cmd != nil {
		do := func() {
			systemExecuteCommandRequest := &octoprint.SystemExecuteCommandRequest{
				Source: octoprint.Core,
				Action: cmd.Action,
			}

			if err := systemExecuteCommandRequest.Do(m.UI.Printer); err != nil {
				utils.LogError("system.createCommandButton()", "Do(SystemExecuteCommandRequest)", err)
				return
			}
		}

		cb = do
		if len(cmd.Confirm) != 0 {
			cb = MustConfirmDialog(m.UI.w, cmd.Confirm, do)
		}
	}

	button := MustButtonImageStyle(name, action + ".svg", style, cb)

	if cmd == nil {
		button.SetSensitive(false)
	}

	return button
}
