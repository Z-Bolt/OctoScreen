package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

var printStatusPanelInstance *printStatusPanel

type printStatusPanel struct {
	CommonPanel
	step *StepButton
	pb   *gtk.ProgressBar

	bed, tool0, tool1, tool2, tool3 *gtk.Button
	file, time, timeLeft            *LabelWithImage
	complete, pause, stop, menu     *gtk.Button
}

func PrintStatusPanel(ui *UI) Panel {
	if printStatusPanelInstance == nil {
		m := &printStatusPanel{CommonPanel: NewCommonPanel(ui, nil)}
		m.panelH = 3
		m.b = NewBackgroundTask(time.Second*2, m.update)
		m.initialize()

		printStatusPanelInstance = m
	}

	return printStatusPanelInstance
}

func (m *printStatusPanel) initialize() {
	defer m.Initialize()

	m.Grid().Attach(m.createInfoBox(), 3, 0, 2, 1)
	m.Grid().Attach(m.createProgressBar(), 3, 1, 2, 1)
	m.Grid().Attach(m.createPauseButton(), 2, 2, 1, 1)
	m.Grid().Attach(m.createStopButton(), 3, 2, 1, 1)
	m.Grid().Attach(m.createMenuButton(), 4, 2, 1, 1)
	m.Grid().Attach(m.createCompleteButton(), 2, 2, 3, 1)

	m.showTools()
}

func (m *printStatusPanel) showTools() {
	toolsCount := m.defineToolsCount()

	m.tool0 = m.createToolButton(0)
	m.tool1 = m.createToolButton(1)
	m.tool2 = m.createToolButton(2)
	m.tool3 = m.createToolButton(3)

	m.bed = m.createBedButton()

	switch toolsCount {
	case 1:
		m.Grid().Attach(m.tool0, 1, 0, 2, 1)
		m.Grid().Attach(m.bed, 1, 1, 2, 1)

	case 2:
		m.Grid().Attach(m.tool0, 1, 0, 1, 1)
		m.Grid().Attach(m.tool1, 2, 0, 1, 1)
		m.Grid().Attach(m.bed, 1, 1, 2, 1)
	case 3:
		m.Grid().Attach(m.tool0, 1, 0, 1, 1)
		m.Grid().Attach(m.tool1, 2, 0, 1, 1)
		m.Grid().Attach(m.tool2, 1, 1, 1, 1)
		m.Grid().Attach(m.bed, 2, 1, 1, 1)
	case 4:
		m.Grid().Attach(m.tool0, 1, 0, 1, 1)
		m.Grid().Attach(m.tool1, 2, 0, 1, 1)
		m.Grid().Attach(m.tool2, 1, 1, 1, 1)
		m.Grid().Attach(m.tool3, 2, 1, 1, 1)
		m.Grid().Attach(m.bed, 1, 2, 1, 1)
	}

}

func (m *printStatusPanel) createCompleteButton() *gtk.Button {
	m.complete = MustButtonImageStyle("Complete", "complete.svg", "color3", func() {
		m.UI.Add(IdleStatusPanel(m.UI))
	})
	return m.complete
}

func (m *printStatusPanel) createProgressBar() *gtk.ProgressBar {
	m.pb = MustProgressBar()
	m.pb.SetShowText(true)
	m.pb.SetMarginTop(10)
	m.pb.SetMarginEnd(m.Scaled(20))
	m.pb.SetVAlign(gtk.ALIGN_CENTER)
	m.pb.SetVExpand(true)

	ctx, _ := m.pb.GetStyleContext()
	ctx.AddClass("printing-progress-bar")

	return m.pb
}

func (m *printStatusPanel) createInfoBox() *gtk.Box {

	m.file = MustLabelWithImage("file.svg", "")
	ctx, _ := m.file.GetStyleContext()
	ctx.AddClass("printing-status-label")

	m.time = MustLabelWithImage("speed-step.svg", "")
	ctx, _ = m.time.GetStyleContext()
	ctx.AddClass("printing-status-label")

	m.timeLeft = MustLabelWithImage("speed-step.svg", "")
	ctx, _ = m.timeLeft.GetStyleContext()
	ctx.AddClass("printing-status-label")

	info := MustBox(gtk.ORIENTATION_VERTICAL, 5)
	info.SetHAlign(gtk.ALIGN_START)
	info.SetHExpand(true)
	info.SetVExpand(true)
	info.SetVAlign(gtk.ALIGN_CENTER)
	info.Add(m.file)
	info.Add(m.time)
	info.Add(m.timeLeft)

	return info
}

func (m *printStatusPanel) createToolButton(num int) *gtk.Button {
	name := fmt.Sprintf("extruder-%d.svg", num+1)
	b := MustButtonImage("", name, func() {})

	ctx, _ := b.GetStyleContext()
	ctx.AddClass("printing-state")
	return b
}

func (m *printStatusPanel) createBedButton() *gtk.Button {
	b := MustButtonImage("", "bed.svg", func() {})

	ctx, _ := b.GetStyleContext()
	ctx.AddClass("printing-state")
	return b
}

func (m *printStatusPanel) createPauseButton() gtk.IWidget {
	m.pause = MustButtonImageStyle("Pause", "pause.svg", "color1", func() {
		defer m.updateTemperature()

		Logger.Warning("Pausing/Resuming job")
		cmd := &octoprint.PauseRequest{Action: octoprint.Toggle}
		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})

	return m.pause
}

func (m *printStatusPanel) createStopButton() gtk.IWidget {
	m.stop = MustButtonImageStyle("Stop", "stop.svg", "color2",
		ConfirmStopDialog(m.UI.w, "Are you sure you want to stop current print?", m),
	)
	return m.stop
}

func (m *printStatusPanel) createMenuButton() gtk.IWidget {
	m.menu = MustButtonImageStyle("Control", "control.svg", "color3", func() {
		m.UI.Add(PrintMenuPanel(m.UI, m))
	})
	return m.menu
}

func (m *printStatusPanel) update() {
	m.updateTemperature()
	m.updateJob()
}

func (m *printStatusPanel) updateTemperature() {
	s, err := (&octoprint.StateRequest{Exclude: []string{"sd"}}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return
	}

	m.doUpdateState(&s.State)

	for tool, s := range s.Temperature.Current {
		text := fmt.Sprintf("%.0f°C / %.0f°C", s.Actual, s.Target)
		switch tool {
		case "bed":
			m.bed.SetLabel(text)
		case "tool0":
			m.tool0.SetLabel(text)
		case "tool1":
			m.tool1.SetLabel(text)
		case "tool2":
			m.tool2.SetLabel(text)
		case "tool3":
			m.tool3.SetLabel(text)
		}
	}
}

func (m *printStatusPanel) doUpdateState(s *octoprint.PrinterState) {
	switch {
	case s.Flags.Printing:
		m.pause.SetSensitive(true)
		m.stop.SetSensitive(true)

		m.pause.Show()
		m.stop.Show()
		m.menu.Show()
		m.back.Show()
		m.complete.Hide()

	case s.Flags.Paused:
		m.pause.SetLabel("Resume")
		m.pause.SetImage(MustImageFromFile("resume.svg"))
		m.pause.SetSensitive(true)
		m.stop.SetSensitive(true)

		m.pause.Show()
		m.stop.Show()
		m.menu.Show()
		m.back.Show()
		m.complete.Hide()
		return
	case s.Flags.Ready:
		m.pause.SetSensitive(false)
		m.stop.SetSensitive(false)

		m.pause.Hide()
		m.stop.Hide()
		m.menu.Hide()
		m.back.Hide()
		m.complete.Show()
	default:
		m.pause.SetSensitive(false)
		m.stop.SetSensitive(false)
	}

	m.pause.SetLabel("Pause")
	m.pause.SetImage(MustImageFromFile("pause.svg"))
}

func (m *printStatusPanel) updateJob() {

	s, err := (&octoprint.JobRequest{}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return
	}

	file := "<i>not-set</i>"
	if s.Job.File.Name != "" {
		file = s.Job.File.Name
		file = strings.Replace(file, ".gcode", "", -1)
		file = strEllipsisLen(file, 20)
	}

	m.file.Label.SetLabel(file)
	m.pb.SetFraction(s.Progress.Completion / 100)

	var timeSpent, timeLeft string
	switch s.Progress.Completion {
	case 100:
		timeSpent = fmt.Sprintf("Completed in %s", time.Duration(int64(s.Job.LastPrintTime)*1e9))
		timeLeft = ""
	case 0:
		timeSpent = "Warming up ..."
		timeLeft = ""
	default:
		Logger.Info(s.Progress.PrintTime)
		e := time.Duration(int64(s.Progress.PrintTime) * 1e9)
		r := time.Duration(int64(s.Progress.PrintTimeLeft) * 1e9)
		timeSpent = fmt.Sprintf("Time: %s", e)
		timeLeft = fmt.Sprintf("Left: %s", r)
	}

	m.time.Label.SetLabel(timeSpent)
	m.timeLeft.Label.SetLabel(timeLeft)
}

func (m *printStatusPanel) defineToolsCount() int {
	c, err := (&octoprint.ConnectionRequest{}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return 0
	}

	profile, err := (&octoprint.PrinterProfilesRequest{Id: c.Current.PrinterProfile}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return 0
	}

	if profile.Extruder.SharedNozzle {
		return 1
	}

	return profile.Extruder.Count
}

func ConfirmStopDialog(parent *gtk.Window, msg string, ma *printStatusPanel) func() {
	return func() {
		win := gtk.MessageDialogNewWithMarkup(
			parent,
			gtk.DIALOG_MODAL,
			gtk.MESSAGE_INFO,
			gtk.BUTTONS_YES_NO,
			"",
		)

		win.SetMarkup(CleanHTML(msg))
		defer win.Destroy()

		box, _ := win.GetContentArea()
		box.SetMarginStart(15)
		box.SetMarginEnd(15)
		box.SetMarginTop(15)
		box.SetMarginBottom(15)

		ctx, _ := win.GetStyleContext()
		ctx.AddClass("dialog")

		ergebnis := win.Run()

		if ergebnis == int(gtk.RESPONSE_YES) {
			Logger.Warning("Stopping job")
			if err := (&octoprint.CancelRequest{}).Do(ma.UI.Printer); err != nil {
				Logger.Error(err)
				return
			}
		}
	}
}
