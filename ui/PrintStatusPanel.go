package ui

import (
	"fmt"
	// "os"
	// "strconv"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gtk"

	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)


type printStatusPanel struct {
	CommonPanel

	progressBar			*gtk.ProgressBar

	tool0Button			*uiWidgets.ToolPrintingButton
	tool1Button			*uiWidgets.ToolPrintingButton
	tool2Button			*uiWidgets.ToolPrintingButton
	tool3Button			*uiWidgets.ToolPrintingButton
	bedButton  			*uiWidgets.ToolPrintingButton

	fileLabel			*utils.LabelWithImage
	timeLabel			*utils.LabelWithImage
	timeLeftLabel		*utils.LabelWithImage
	completeButton		*gtk.Button
	pauseButton			*gtk.Button
	stopButton			*gtk.Button
	menuButton			*gtk.Button

	backgroundTask		*utils.BackgroundTask
}

var printStatusPanelInstance *printStatusPanel

func GetPrintStatusPanelInstance(ui *UI) *printStatusPanel {
	if printStatusPanelInstance == nil {
		printStatusPanelInstance = &printStatusPanel{
			CommonPanel: CreateTopLevelCommonPanel("PrintStatusPanel", ui),
		}
		printStatusPanelInstance.initialize()
		printStatusPanelInstance.createBackgroundTask()
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
	hotendCount := utils.GetHotendCount(this.UI.Client)
	if hotendCount == 1 {
		this.tool0Button = uiWidgets.CreateToolPrintingButton(0)
	} else {
		this.tool0Button = uiWidgets.CreateToolPrintingButton(1)
	}
	this.tool1Button = uiWidgets.CreateToolPrintingButton( 2)
	this.tool2Button = uiWidgets.CreateToolPrintingButton( 3)
	this.tool3Button = uiWidgets.CreateToolPrintingButton( 4)
	this.bedButton   = uiWidgets.CreateToolPrintingButton(-1)

	switch hotendCount {
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
		this.UI.GoToPanel(GetIdleStatusPanelInstance(this.UI))
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
	this.fileLabel = utils.MustLabelWithImage("file-gcode.svg", "")
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

		logger.Info("Pausing/Resuming job")
		cmd := &octoprintApis.PauseRequest{Action: dataModels.Toggle}
		err := cmd.Do(this.UI.Client)
		logger.Info("Pausing/Resuming job 2, Do() was just called")

		if err != nil {
			logger.LogError("PrintStatusPanel.createPauseButton()", "Do(PauseRequest)", err)
			return
		}

		logger.Info("Pausing/Resuming job 2c")
	})

	return this.pauseButton
}

func (this *printStatusPanel) createCancelButton() gtk.IWidget {
	this.stopButton = utils.MustButtonImageStyle(
		"Cancel",
		"stop.svg",
		"color-warning-sign-yellow",
		confirmStopDialogBox(this.UI.window, "Are you sure you want to cancel the current print?", this),
	)

	return this.stopButton
}

func (this *printStatusPanel) createControlButton() gtk.IWidget {
	this.menuButton = utils.MustButtonImageStyle(
		"Control",
		"printing-control.svg",
		"color3",
		func() {
			this.UI.GoToPanel(GetPrintMenuPanelInstance(this.UI))
		},
	)
	return this.menuButton
}

func (this *printStatusPanel) createBackgroundTask() {
	logger.TraceEnter("PrintStatusPanel.createBackgroundTask()")

	// Default timeout of 1 second.
	duration := utils.GetExperimentalFrequency(1, "EXPERIMENTAL_PRINT_UPDATE_FREQUENCY")
	this.backgroundTask = utils.CreateBackgroundTask(duration, this.update)
	// Update the UI every second, but the data is only updated once every 10 seconds.
	// See OctoPrintResponseManager.update(). 
	this.backgroundTask.Start()

	logger.TraceLeave("PrintStatusPanel.createBackgroundTask()")
}

func (this *printStatusPanel) update() {
	logger.TraceEnter("PrintStatusPanel.update()")

	this.updateTemperature()
	this.updateJob()

	logger.TraceLeave("PrintStatusPanel.update()")
}

func (this *printStatusPanel) updateTemperature() {
	logger.TraceEnter("PrintStatusPanel.updateTemperature()")

	octoPrintResponseManager := GetOctoPrintResponseManagerInstance(this.UI)
	if octoPrintResponseManager.IsConnected() != true {
		// If not connected, do nothing and leave.
		logger.TraceLeave("PrintStatusPanel.updateTemperature() (not connected)")
		return
	}

	this.doUpdateState(&octoPrintResponseManager.FullStateResponse.State)

	for tool, currentTemperatureData := range octoPrintResponseManager.FullStateResponse.Temperature.CurrentTemperatureData {
		text := utils.GetTemperatureDataString(currentTemperatureData)
		switch tool {
			case "bed":
				logger.Debug("Updating the UI's bed temp")
				// this.bedButton.SetTemperatures(currentTemperatureData)
				this.bedButton.SetLabel(text)

			case "tool0":
				logger.Debug("Updating the UI's tool0 temp")
				// this.tool0Button.SetTemperatures(currentTemperatureData)
				this.tool0Button.SetLabel(text)

			case "tool1":
				logger.Debug("Updating the UI's tool1 temp")
				// this.tool1Button.SetTemperatures(currentTemperatureData)
				this.tool1Button.SetLabel(text)

			case "tool2":
				logger.Debug("Updating the UI's tool2 temp")
				// this.tool2Button.SetTemperatures(currentTemperatureData)
				this.tool2Button.SetLabel(text)

			case "tool3":
				logger.Debug("Updating the UI's tool3 temp")
				// this.tool3Button.SetTemperatures(currentTemperatureData)
				this.tool3Button.SetLabel(text)

			default:
				logger.Errorf("PrintStatusPanel.updateTemperature() - GetOctoPrintResponseManagerInstance() returned an unknown tool: %q", tool)
		}
	}

	logger.TraceLeave("PrintStatusPanel.updateTemperature()")
}

func (this *printStatusPanel) doUpdateState(printerState *dataModels.PrinterState) {
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
			logLevel := logger.LogLevel()
			if logLevel == "debug" {
				logger.Fatalf("PrintStatusPanel.doUpdateState() - unknown printerState.Flags")
			}

			this.pauseButton.SetSensitive(false)
			this.stopButton.SetSensitive(false)
	}

	this.pauseButton.SetLabel("Pause")
	pauseImage := utils.MustImageFromFile("pause.svg")
	this.pauseButton.SetImage(pauseImage)
}

func (this *printStatusPanel) updateJob() {
	logger.TraceEnter("PrintStatusPanel.updateJob()")

	jobResponse, err := (&octoprintApis.JobRequest{}).Do(this.UI.Client)
	if err != nil {
		logger.LogError("PrintStatusPanel.updateJob()", "Do(JobRequest)", err)
		logger.TraceLeave("PrintStatusPanel.updateJob()")
		return
	}

	jobFileName := "<i>not-set</i>"
	if jobResponse.Job.File.Name != "" {
		jobFileName = jobResponse.Job.File.Name
		jobFileName = strings.Replace(jobFileName, ".gcode", "", -1)
		jobFileName = strings.Replace(jobFileName, ".gco", "", -1)
		jobFileName = utils.TruncateString(jobFileName, 20)
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
			logger.Info(jobResponse.Progress.PrintTime)
			printTime := time.Duration(int64(jobResponse.Progress.PrintTime) * 1e9)
			printTimeLeft := time.Duration(int64(jobResponse.Progress.PrintTimeLeft) * 1e9)
			timeSpent = fmt.Sprintf("Time: %s", printTime)
			timeLeft = fmt.Sprintf("Left: %s", printTimeLeft)
	}

	this.timeLabel.Label.SetLabel(timeSpent)
	this.timeLeftLabel.Label.SetLabel(timeLeft)

	logger.TraceLeave("PrintStatusPanel.updateJob()")
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
			logger.Warn("Stopping job")
			if err := (&octoprintApis.CancelRequest{}).Do(printStatusPanel.UI.Client); err != nil {
				logger.LogError("PrintStatusPanel.confirmStopDialogBox()", "Do(CancelRequest)", err)
				return
			}
		}
	}
}
