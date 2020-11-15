package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)

var printStatusPanelInstance *printStatusPanel

type printStatusPanel struct {
	CommonPanel

	progressBar *gtk.ProgressBar

	tool0Button *uiWidgets.ToolButton
	tool1Button *uiWidgets.ToolButton
	tool2Button *uiWidgets.ToolButton
	tool3Button *uiWidgets.ToolButton
	bedButton   *uiWidgets.ToolButton

	fileLabel      *utils.LabelWithImage
	timeLabel      *utils.LabelWithImage
	timeLeftLabel  *utils.LabelWithImage
	completeButton *gtk.Button
	pauseButton    *gtk.Button
	stopButton     *gtk.Button
	menuButton     *gtk.Button
}

func PrintStatusPanel(ui *UI) *printStatusPanel {
	if printStatusPanelInstance == nil {
		instance := &printStatusPanel{
			CommonPanel: NewTopLevelCommonPanel(ui, nil),
		}

		// TODO: revisit... some set the background task and the initialize
		// and others initialize and then set the background task
		instance.backgroundTask = utils.CreateBackgroundTask(time.Second * 2, instance.update)
		instance.initialize()
		printStatusPanelInstance = instance
	}

	return printStatusPanelInstance
}

func (this *printStatusPanel) initialize() {
	defer this.Initialize()

	this.Grid().Attach(this.createInfoBox(),        2, 0, 2, 1)

	this.Grid().Attach(this.createProgressBar(),    2, 1, 2, 1)

	this.Grid().Attach(this.createPauseButton(),    1, 2, 1, 1)
	this.Grid().Attach(this.createCancelButton(),   2, 2, 1, 1)
	this.Grid().Attach(this.createControlButton(),  3, 2, 1, 1)

	this.Grid().Attach(this.createCompleteButton(), 1, 2, 3, 1)

	this.showTools()
}

func (this *printStatusPanel) showTools() {
	// Note: The creation and initialization of the tool buttons in IdleStatusPanel and
	// PrintStatusPanel look similar, but there are subtle differences between the two
	// and they can't be reused.
	toolheadCount := utils.GetToolheadCount(this.UI.Client)
	if toolheadCount == 1 {
		this.tool0Button = uiWidgets.CreateToolButton(0, this.UI.Client)
	} else {
		this.tool0Button = uiWidgets.CreateToolButton(1, this.UI.Client)
	}
	this.tool1Button = uiWidgets.CreateToolButton(2, this.UI.Client)
	this.tool2Button = uiWidgets.CreateToolButton(3, this.UI.Client)
	this.tool3Button = uiWidgets.CreateToolButton(4, this.UI.Client)
	this.bedButton = uiWidgets.CreateToolButton(-1, this.UI.Client)

	switch toolheadCount {
		case 1:
			this.Grid().Attach(this.tool0Button, 0, 0, 2, 1)
			this.Grid().Attach(this.bedButton,   0, 1, 2, 1)

		case 2:
			this.Grid().Attach(this.tool0Button, 0, 0, 1, 1)
			this.Grid().Attach(this.tool1Button, 1, 0, 1, 1)
			this.Grid().Attach(this.bedButton,   0, 1, 2, 1)

		case 3:
			this.Grid().Attach(this.tool0Button, 0, 0, 1, 1)
			this.Grid().Attach(this.tool1Button, 1, 0, 1, 1)
			this.Grid().Attach(this.tool2Button, 0, 1, 1, 1)
			this.Grid().Attach(this.bedButton,   1, 1, 1, 1)

		case 4:
			this.Grid().Attach(this.tool0Button, 0, 0, 1, 1)
			this.Grid().Attach(this.tool1Button, 1, 0, 1, 1)
			this.Grid().Attach(this.tool2Button, 0, 1, 1, 1)
			this.Grid().Attach(this.tool3Button, 1, 1, 1, 1)
			this.Grid().Attach(this.bedButton,   0, 2, 1, 1)
	}
}

func (this *printStatusPanel) createCompleteButton() *gtk.Button {
	this.completeButton = utils.MustButtonImageStyle("Complete", "complete.svg", "color3", func() {
		this.UI.GoToPanel(IdleStatusPanel(this.UI))
	})

	return this.completeButton
}

func (this *printStatusPanel) createProgressBar() *gtk.ProgressBar {
	this.progressBar = utils.MustProgressBar()
	this.progressBar.SetShowText(true)
	this.progressBar.SetMarginTop(10)
	this.progressBar.SetMarginEnd(this.Scaled(20))
	this.progressBar.SetVAlign(gtk.ALIGN_CENTER)
	this.progressBar.SetVExpand(true)

	ctx, _ := this.progressBar.GetStyleContext()
	ctx.AddClass("printing-progress-bar")

	return this.progressBar
}

func (this *printStatusPanel) createInfoBox() *gtk.Box {
	this.fileLabel = utils.MustLabelWithImage("file-stl.svg", "")
	ctx, _ := this.fileLabel.GetStyleContext()
	ctx.AddClass("printing-status-label")

	this.timeLabel = utils.MustLabelWithImage("time.svg", "")
	ctx, _ = this.timeLabel.GetStyleContext()
	ctx.AddClass("printing-status-label")

	this.timeLeftLabel = utils.MustLabelWithImage("time.svg", "")
	ctx, _ = this.timeLeftLabel.GetStyleContext()
	ctx.AddClass("printing-status-label")

	infoBox := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	infoBox.SetHAlign(gtk.ALIGN_START)
	infoBox.SetHExpand(true)
	infoBox.SetVExpand(true)
	infoBox.SetVAlign(gtk.ALIGN_CENTER)
	infoBox.Add(this.fileLabel)
	infoBox.Add(this.timeLabel)
	infoBox.Add(this.timeLeftLabel)

	return infoBox
}

func (this *printStatusPanel) createPauseButton() gtk.IWidget {
	this.pauseButton = utils.MustButtonImageStyle("Pause", "pause.svg", "color-warning-sign-yellow", func() {
		defer this.updateTemperature()

		utils.Logger.Info("Pausing/Resuming job")
		cmd := &octoprint.PauseRequest{Action: octoprint.Toggle}
		err := cmd.Do(this.UI.Client)
		utils.Logger.Info("Pausing/Resuming job 2, Do() was just called")

		if err != nil {
			utils.LogError("print_status.createPauseButton()", "Do(PauseRequest)", err)
			return
		}

		utils.Logger.Info("Pausing/Resuming job 2c")
	})

	return this.pauseButton
}

func (this *printStatusPanel) createCancelButton() gtk.IWidget {
	this.stopButton = utils.MustButtonImageStyle(
		"Cancel",
		"stop.svg",
		"color-warning-sign-yellow",
		confirmStopDialogBox(this.UI.window, "Are you sure you want to stop the current print?", this),
	)

	return this.stopButton
}

func (this *printStatusPanel) createControlButton() gtk.IWidget {
	this.menuButton = utils.MustButtonImageStyle(
		"Control",
		"printing-control.svg",
		"color3",
		func() {
			this.UI.GoToPanel(PrintMenuPanel(this.UI, this))
		},
	)
	return this.menuButton
}

func (this *printStatusPanel) update() {
	//Logger.Printf("Now in print_status.update()")

	this.updateTemperature()
	this.updateJob()

	//Logger.Printf("Now leaving print_status.update()")
}

func (this *printStatusPanel) updateTemperature() {
	//Logger.Printf("Now in print_status.updateTemperature()")

	fullStateResponse, err := (&octoprint.FullStateRequest{Exclude: []string{"sd"}}).Do(this.UI.Client)
	if err != nil {
		utils.LogError("print_status.updateTemperature()", "Do(StateRequest)", err)
		return
	}

	this.doUpdateState(&fullStateResponse.State)

	for tool, currentTemperatureData := range fullStateResponse.Temperature.CurrentTemperatureData {
		text := utils.GetTemperatureDataString(currentTemperatureData)
		switch tool {
			case "bed":
				this.bedButton.SetLabel(text)

			case "tool0":
				this.tool0Button.SetLabel(text)

			case "tool1":
				this.tool1Button.SetLabel(text)

			case "tool2":
				this.tool2Button.SetLabel(text)

			case "tool3":
				this.tool3Button.SetLabel(text)
		}
	}

	//Logger.Printf("Now leaving print_status.updateTemperature()")
}

func (this *printStatusPanel) doUpdateState(printerState *octoprint.PrinterState) {
	switch {
		case printerState.Flags.Printing:
			this.pauseButton.SetSensitive(true)
			this.stopButton.SetSensitive(true)

			this.pauseButton.Show()
			this.stopButton.Show()
			if this.menuButton != nil {
				this.menuButton.Show()
			}
			this.backButton.Show()
			this.completeButton.Hide()

		case printerState.Flags.Paused:
			this.pauseButton.SetLabel("Resume")
			resumeImage := utils.MustImageFromFile("resume.svg")
			this.pauseButton.SetImage(resumeImage)
			this.pauseButton.SetSensitive(true)
			this.stopButton.SetSensitive(true)

			this.pauseButton.Show()
			this.stopButton.Show()
			if this.menuButton != nil {
				this.menuButton.Show()
			}
			this.backButton.Show()
			this.completeButton.Hide()
			return

		case printerState.Flags.Ready:
			this.pauseButton.SetSensitive(false)
			this.stopButton.SetSensitive(false)

			this.pauseButton.Hide()
			this.stopButton.Hide()
			if this.menuButton != nil {
				this.menuButton.Hide()
			}
			this.backButton.Hide()
			this.completeButton.Show()

		default:
			logLevel := utils.LowerCaseLogLevel()
			if logLevel == "debug" {
				utils.Logger.Fatalf("PrintStatusPanel.doUpdateState() - unknown printerState.Flags")
			}

			this.pauseButton.SetSensitive(false)
			this.stopButton.SetSensitive(false)
	}

	this.pauseButton.SetLabel("Pause")
	pauseImage := utils.MustImageFromFile("pause.svg")
	this.pauseButton.SetImage(pauseImage)
}

func (this *printStatusPanel) updateJob() {
	//Logger.Printf("Now in print_status.updateJob()")

	jobResponse, err := (&octoprint.JobRequest{}).Do(this.UI.Client)
	if err != nil {
		utils.LogError("print_status.updateJob()", "Do(JobRequest)", err)
		return
	}

	jobFileName := "<i>not-set</i>"
	if jobResponse.Job.File.Name != "" {
		jobFileName = jobResponse.Job.File.Name
		jobFileName = strings.Replace(jobFileName, ".gcode", "", -1)
		jobFileName = utils.StrEllipsisLen(jobFileName, 20)
	}

	this.fileLabel.Label.SetLabel(jobFileName)
	this.progressBar.SetFraction(jobResponse.Progress.Completion / 100)

	var timeSpent, timeLeft string
	switch jobResponse.Progress.Completion {
		case 100:
			timeSpent = fmt.Sprintf("Completed in %s", time.Duration(int64(jobResponse.Job.LastPrintTime) * 1e9))
			timeLeft = ""

		case 0:
			timeSpent = "Warming up ..."
			timeLeft = ""

		default:
			utils.Logger.Info(jobResponse.Progress.PrintTime)
			printTime := time.Duration(int64(jobResponse.Progress.PrintTime) * 1e9)
			printTimeLeft := time.Duration(int64(jobResponse.Progress.PrintTimeLeft) * 1e9)
			timeSpent = fmt.Sprintf("Time: %s", printTime)
			timeLeft = fmt.Sprintf("Left: %s", printTimeLeft)
	}

	this.timeLabel.Label.SetLabel(timeSpent)
	this.timeLeftLabel.Label.SetLabel(timeLeft)

	//Logger.Printf("Now leaving print_status.updateJob()")
}

func confirmStopDialogBox(
	parentWindow		*gtk.Window,
	message				string,
	printStatusPanel	*printStatusPanel,
) func() {
	return func() {
		dialogBox := gtk.MessageDialogNewWithMarkup(
			parentWindow,
			gtk.DIALOG_MODAL,
			gtk.MESSAGE_INFO,
			gtk.BUTTONS_YES_NO,
			"",
		)

		dialogBox.SetMarkup(utils.CleanHTML(message))
		defer dialogBox.Destroy()

		box, _ := dialogBox.GetContentArea()
		box.SetMarginStart(15)
		box.SetMarginEnd(15)
		box.SetMarginTop(15)
		box.SetMarginBottom(15)

		ctx, _ := dialogBox.GetStyleContext()
		ctx.AddClass("dialog")

		userResponse := dialogBox.Run()
		if userResponse == int(gtk.RESPONSE_YES) {
			utils.Logger.Warning("Stopping job")
			if err := (&octoprint.CancelRequest{}).Do(printStatusPanel.UI.Client); err != nil {
				utils.LogError("print_status.confirmStopDialogBox()", "Do(CancelRequest)", err)
				return
			}
		}
	}
}
