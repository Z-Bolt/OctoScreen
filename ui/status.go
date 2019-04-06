package ui

import (
	"fmt"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

var statusPanelInstance *statusPanel

type statusPanel struct {
	CommonPanel
	step *StepButton
	pb   *gtk.ProgressBar

	bed, tool0                 *LabelWithImage
	file, left, finish         *LabelWithImage
	//status			   *LabelWithImage
	print, pause, stop *gtk.Button
	taskRunning bool
	printerStatus uint8
}

func StatusPanel(ui *UI, parent Panel) Panel {
	if statusPanelInstance == nil {
		m := &statusPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.b = NewBackgroundTask(time.Second*1, m.update)
		m.initialize()
		m.taskRunning = false
		m.printerStatus = 0
		statusPanelInstance = m
	}

	return statusPanelInstance
}

func (m *statusPanel) initialize() {
	defer m.Initialize()

	m.Grid().Attach(m.createMainBox(), 1, 0, 4, 2)


}

func (m *statusPanel) createProgressBar() *gtk.ProgressBar {
	m.pb = MustProgressBar()
	m.pb.SetShowText(true)
	//m.pb.SetMarginTop(12)
	//m.pb.SetMarginStart(5)
	//m.pb.SetMarginEnd(5)
	//m.pb.SetMarginBottom(20)
	m.pb.SetName("PrintProg")
	return m.pb
}

func (m *statusPanel) createMainBox() *gtk.Box {

	box := MustBox(gtk.ORIENTATION_VERTICAL, 5)
	box.SetVAlign(gtk.ALIGN_START)
	box.SetVExpand(true)

	grid := MustGrid()
	grid.SetHExpand(true)
	grid.Add(m.createInfoBox())
	grid.SetVAlign(gtk.ALIGN_START)
	grid.SetMarginTop(20)
	
	box.Add(grid)

	pb_box := MustBox(gtk.ORIENTATION_VERTICAL, 5)
	pb_box.SetVExpand(true)
	pb_box.SetHExpand(true)
	//pb_box.Add(MustButton(MustImageFromFileWithSize("back.svg", 60, 60), m.UI.GoHistory))
	pb_box.Add(m.createProgressBar())	
	box.Add(pb_box)

	butt := MustBox(gtk.ORIENTATION_HORIZONTAL, 0)
	butt.SetHAlign(gtk.ALIGN_END)
	butt.SetVAlign(gtk.ALIGN_END)
	butt.SetVExpand(true)
	butt.SetMarginTop(0)
	butt.SetMarginEnd(0)
	butt.Add(m.createPrintButton())
	butt.Add(m.createPauseButton())
	butt.Add(m.createStopButton())
	butt.Add(MustButton(MustImageFromFileWithSize("back.svg", 60, 60), m.UI.GoHistory))
	box.Add(butt)
	return box
}

func (m *statusPanel) createInfoBox() *gtk.Box {
	m.file = MustLabelWithImage("file.svg", "")
	m.file.SetName("NameLabel")
	m.left = MustLabelWithImage("speed-step.svg", "")
	m.left.SetName("TimeLabel")
	m.finish = MustLabelWithImage("finish.svg", "")
	m.finish.SetName("TimeLabel")
	m.bed = MustLabelWithImage("bed.svg", "")
	m.bed.SetName("TempLabel")
	m.tool0 = MustLabelWithImage("extruder.svg", "")
	m.tool0.SetName("TempLabel")
//	m.status = MustLabelWithImage("file.svg", "")

	info := MustBox(gtk.ORIENTATION_VERTICAL, 5)
	info.SetHAlign(gtk.ALIGN_START)
	info.SetHExpand(true)
	info.Add(m.file)
	info.Add(m.left)
	info.Add(m.finish)
	info.Add(m.tool0)
	info.Add(m.bed)
//	info.Add(m.status)
	info.SetMarginStart(20)

	return info
}

func (m *statusPanel) createPrintButton() gtk.IWidget {
	m.print = MustButton(MustImageFromFileWithSize("status.svg", 60, 60), func() {
		defer m.updateTemperature() // modified dark

		Logger.Warning("Starting a new job")
		if err := (&octoprint.StartRequest{}).Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})

	return m.print
}

func (m *statusPanel) createPauseButton() gtk.IWidget {
	m.pause = MustButton(MustImageFromFileWithSize("pause.svg", 60, 60), func() {
		defer m.updateTemperature()

		Logger.Warning("Pausing/Resuming job")
		if err := (&octoprint.PauseRequest{Action: octoprint.Toggle}).Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})

	return m.pause
}

func (m *statusPanel) createStopButton() gtk.IWidget {
	m.stop = MustButton(MustImageFromFileWithSize("stop.svg", 60, 60), 
			 ConfirmStopDialog(m.UI.w, "Are you sure you want to stop current print?", m), 
			 )
	return m.stop
}

func (m *statusPanel) update() {
	if m.taskRunning == false {
		m.taskRunning = true
		m.updateTemperature()
		m.updateJob()
		m.taskRunning = false
	}
}

func (m *statusPanel) updateTemperature() {
	s, err := (&octoprint.StateRequest{Exclude: []string{"sd"}}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return
	}
	m.doUpdateState(&s.State)

	for tool, s := range s.Temperature.Current {
		text := fmt.Sprintf("%.0f°C ⇒ %.0f°C ", s.Actual, s.Target)
		switch tool {
		case "bed":
			m.bed.Label.SetLabel(text)
		case "tool0":
			m.tool0.Label.SetLabel(text)
		}
	}
}

func (m *statusPanel) doUpdateState(s *octoprint.PrinterState) {

	currentPrinterStatus := btou(s.Flags.Printing) * 4 + btou(s.Flags.Paused) * 2 + btou(s.Flags.Ready) // printer status value
//	text := fmt.Sprintf("Status: %d", currentPrinterStatus)
//	m.status.Label.SetLabel(text)
	
	if currentPrinterStatus != m.printerStatus {
		m.printerStatus = currentPrinterStatus
		switch currentPrinterStatus{
		case 4: // printing
			m.print.SetSensitive(false)
			m.pause.SetImage(MustImageFromFileWithSize("pause.svg", 60, 60))
			m.pause.SetSensitive(true)
			m.stop.SetSensitive(true)
			return
		case 3: // paused
			m.print.SetSensitive(false)
			m.pause.SetImage(MustImageFromFileWithSize("resume.svg", 60, 60))
			m.pause.SetSensitive(true)
			m.stop.SetSensitive(true)
			return
		case 1: // ready
			m.print.SetSensitive(true)
			m.pause.SetSensitive(false)
			m.stop.SetSensitive(false)
			return
		default:
			m.print.SetSensitive(false)
			m.pause.SetSensitive(false)
			m.stop.SetSensitive(false)
			return
		}
	}
}

func (m *statusPanel) updateJob() {
	s, err := (&octoprint.JobRequest{}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return

	}
	file := "<i>File not set</i>"
	if s.Job.File.Name != "" {
		file = filenameEllipsis_long(s.Job.File.Name)
	}

	m.file.Label.SetLabel(fmt.Sprintf("%s", file))
	m.pb.SetFraction(s.Progress.Completion / 100)

	if m.UI.State.IsOperational() {
		m.left.Label.SetLabel("Printer is ready")
		m.finish.Label.SetLabel("-")
		return
	}

	var text string
	var finishText string
    finishText = fmt.Sprintf("-")
	switch s.Progress.Completion {
	case 100:
		text = fmt.Sprintf("Completed in %s", time.Duration(int64(s.Job.LastPrintTime)*1e9))
	case 0:
		text = "Warming up ..."
	default:
		e := time.Duration(int64(s.Progress.PrintTime) * 1e9)
		l := time.Duration(int64(s.Progress.PrintTimeLeft) * 1e9)
		f := time.Now().Local().Add(time.Duration(int64(s.Progress.PrintTimeLeft)) * time.Second)
		text = fmt.Sprintf("Elapsed: %s / Left: %s", e, l)
		finishText = fmt.Sprintf("Finish time: %s", f.Format("15:04 02-Jan-06"))
		if l == 0 {
			text = fmt.Sprintf("Elapsed: %s", e)
		}
	}
	m.left.Label.SetLabel(text)
	m.finish.Label.SetLabel(finishText)
}

func filenameEllipsis_long(name string) string {
	if len(name) > 35 {
		return name[:32] + "…"
	}

	return name
}

func filenameEllipsis(name string) string {
	if len(name) > 31 {
		return name[:28] + "…"
	}

	return name
}

func filenameEllipsis_short(name string) string {
	if len(name) > 27 {
		return name[:24] + "…"
	}

	return name
}

func btou(b bool) uint8 {
        if b {
                return 1
        }
        return 0
}

func ConfirmStopDialog(parent *gtk.Window, msg string, ma *statusPanel) func() {
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

