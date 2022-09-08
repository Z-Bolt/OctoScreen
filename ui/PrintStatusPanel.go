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

	tool0Button				*uiWidgets.ToolPrintingButton
	tool1Button				*uiWidgets.ToolPrintingButton
	tool2Button				*uiWidgets.ToolPrintingButton
	tool3Button				*uiWidgets.ToolPrintingButton
	bedButton  				*uiWidgets.ToolPrintingButton

	fileLabelWithImage		*utils.LabelWithImage
	timeLabelWithImage		*utils.LabelWithImage
	timeLeftLabelWithImage	*utils.LabelWithImage
	// layerLabelWithImage	*utils.LabelWithImage
	// The info for the current / total layers is not available
	// See https://community.octoprint.org/t/layer-number-and-total-layers-from-api/8005/4
	// and https://docs.octoprint.org/en/master/api/datamodel.html#sec-api-datamodel-jobs-job
	// darn.
	
	progressBar				*gtk.ProgressBar

	pauseButton				*gtk.Button
	cancelButton			*gtk.Button
	controlButton			*gtk.Button
	completeButton			*gtk.Button

	backgroundTask			*utils.BackgroundTask
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

	this.createToolButtons()
}


func (this *printStatusPanel) createInfoBox() *gtk.Box {
	this.fileLabelWithImage = utils.MustLabelWithImage("file-gcode.svg", "")
	ctx, _ := this.fileLabelWithImage.GetStyleContext()
	ctx.AddClass("printing-status-label")

	this.timeLabelWithImage = utils.MustLabelWithImage("time.svg", "Print Time:")
	ctx, _ = this.timeLabelWithImage.GetStyleContext()
	ctx.AddClass("printing-status-label")

	this.timeLeftLabelWithImage = utils.MustLabelWithImage("time.svg", "Print Time Left:")
	ctx, _ = this.timeLeftLabelWithImage.GetStyleContext()
	ctx.AddClass("printing-status-label")

	// this.layerLabelWithImage = utils.MustLabelWithImage("time.svg", "")
	// ctx, _ = this.layerLabelWithImage.GetStyleContext()
	// ctx.AddClass("printing-status-label")

	infoBox := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	infoBox.SetHAlign(gtk.ALIGN_START)
	infoBox.SetHExpand(true)
	infoBox.SetVExpand(true)
	infoBox.SetVAlign(gtk.ALIGN_CENTER)
	infoBox.Add(this.fileLabelWithImage)
	infoBox.Add(this.timeLabelWithImage)
	infoBox.Add(this.timeLeftLabelWithImage)
	// infoBox.Add(this.layerLabelWithImage)

	return infoBox
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

func (this *printStatusPanel) createPauseButton() gtk.IWidget {
	/*
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
	*/

	this.pauseButton = utils.MustButtonImage(
		"Pause",
		"pause.svg",
		this.handlePauseClicked,
	)

	return this.pauseButton
}

func (this *printStatusPanel) createCancelButton() gtk.IWidget {
	this.cancelButton = utils.MustButtonImageStyle(
		"Cancel",
		"stop.svg",
		"color-warning-sign-yellow",
		this.handleCancelClicked,
	)

	return this.cancelButton
}

func (this *printStatusPanel) createControlButton() gtk.IWidget {
	this.controlButton = utils.MustButtonImageStyle(
		"Control",
		"printing-control.svg",
		"color3",
		this.handleControlClicked,
	)

	return this.controlButton
}

func (this *printStatusPanel) createCompleteButton() *gtk.Button {
	this.completeButton = utils.MustButtonImageStyle(
		"Complete",
		"complete.svg",
		"color3",
		this.handleCompleteClicked,
	)

	return this.completeButton
}











func (this *printStatusPanel) createToolButtons() {
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
			this.cancelButton.SetSensitive(true)

			this.pauseButton.Show()
			this.cancelButton.Show()
			if this.controlButton != nil {
				this.controlButton.Show()
			}
			this.backButton.Show()
			this.completeButton.Hide()

		case printerState.Flags.Paused:
			this.pauseButton.SetLabel("Resume")
			resumeImage := utils.MustImageFromFile("resume.svg")
			this.pauseButton.SetImage(resumeImage)
			this.pauseButton.SetSensitive(true)
			this.cancelButton.SetSensitive(true)

			this.pauseButton.Show()
			this.cancelButton.Show()
			if this.controlButton != nil {
				this.controlButton.Show()
			}
			this.backButton.Show()
			this.completeButton.Hide()
			return

		case printerState.Flags.Ready:
			this.pauseButton.SetSensitive(false)
			this.cancelButton.SetSensitive(false)

			this.pauseButton.Hide()
			this.cancelButton.Hide()
			if this.controlButton != nil {
				this.controlButton.Hide()
			}
			this.backButton.Hide()
			this.completeButton.Show()

		default:
			logLevel := logger.LogLevel()
			if logLevel == "debug" {
				logger.Fatalf("PrintStatusPanel.doUpdateState() - unknown printerState.Flags")
			}

			this.pauseButton.SetSensitive(false)
			this.cancelButton.SetSensitive(false)
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

	this.fileLabelWithImage.Label.SetLabel(jobFileName)
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

	this.timeLabelWithImage.Label.SetLabel(timeSpent)
	this.timeLeftLabelWithImage.Label.SetLabel(timeLeft)

	logger.TraceLeave("PrintStatusPanel.updateJob()")
}




func (this *printStatusPanel) handlePauseClicked() {
	logger.TraceEnter("PrintStatusPanel.handlePauseClicked()")

	// TODO: is this needed?
	// defer this.updateTemperature()

	cmd := &octoprintApis.PauseRequest{Action: dataModels.Toggle}
	err := cmd.Do(this.UI.Client)
	if err != nil {
		logger.LogError("PrintStatusPanel.handlePauseClicked()", "Do(PauseRequest)", err)
		return
	}

	label, _ := this.pauseButton.GetLabel()
	if label == "Pause" {
		this.pauseButton.SetLabel("Resume")
		resumeImage := utils.MustImageFromFile("resume.svg")
		this.pauseButton.SetImage(resumeImage)
	} else {
		this.pauseButton.SetLabel("Pause")
		pauseImage := utils.MustImageFromFile("pause.svg")
		this.pauseButton.SetImage(pauseImage)
	}

	// this.pauseButton.SetSensitive(true)
	// this.pauseButton.Show()

	logger.TraceLeave("PrintStatusPanel.handlePauseClicked()")
}

func (this *printStatusPanel) handleCancelClicked() {
	this.confirmCancelDialogBox(this.UI.window, "Are you sure you want to cancel the current print?", this)
}

func (this *printStatusPanel) handleControlClicked() {
	// TODO:
	this.UI.GoToPanel(GetPrintMenuPanelInstance(this.UI))
	// call this.UI.GoToPanel()
	// or all this.UI.Goback() ?
}

func (this *printStatusPanel) handleCompleteClicked() {
	// TODO:
	this.UI.GoToPanel(GetIdleStatusPanelInstance(this.UI))
	// call this.UI.GoToPanel()
	// or call this.UI.Goback() ?
}





func (this *printStatusPanel) confirmCancelDialogBox(
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
			err := (&octoprintApis.CancelRequest{}).Do(printStatusPanel.UI.Client)
			if err == nil {
				// TODO: remove
				logger.Warn("err was nil")


				// pauseButton			*gtk.Button
				// stopButton			*gtk.Button
				// controlButton		*gtk.Button
				printStatusPanel.pauseButton.SetSensitive(false)
				printStatusPanel.cancelButton.SetSensitive(false)
				printStatusPanel.controlButton.SetSensitive(false)
			} else {
				logger.LogError("PrintStatusPanel.confirmCancelDialogBox()", "Do(CancelRequest)", err)
			}
		}
	}
}

func formattedDuration(duration time.Duration) string {
	hours := duration / time.Hour
	duration -= hours * time.Hour

	minutes := duration / time.Minute
	duration -= minutes * time.Minute

	seconds := duration / time.Second

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}


/*
TODO issues:

1. start a print
as the printer is warming up, click Cancel
-> all three buttons turn gray, but then the Cancel button becomes enabled again
-> the status is "Printing" I think...


2 is either...
2a. start a print
as the printer is warming up, click Cancel
we're eventually taken to the Idle screen
-> the hotend and bed buttons are disabed
    -OR-
2b. start a print
the printer IS warmed up, AND starts printing
click Cancel
we're eventually taken to the Idle screen
-> the hotend and bed buttons are disabed

...this might (maybe?) be due to going to the Idle panel, but the state not being "operational" yet


3. with the printer cold (or cold-ish)
start a print
-> some of the time, the print times are displayed
-> but some of the time, the print times are not displayed
...need to:
	* a) not display the clock icons until the time (text) is displayed
	* b) figure out why the time (values/text) appears some of the time and doesn't appear some of the time


4. Start a print
click the Pause button
then click the Cancel button
then confirm

-> when paused, and then canceled, the app might get into a weird state
	...will need to play around with this and dig ino this some more





	PAUSE   STOP   CONTROL

a) panel is created
	enabled enabled enabled

b) user clicks cancel
	confirm
	disabled disabled disabled

c) user clicks pause
	unknown printerState.Text: Pausing
	Pause button changes to Resume
	enabled enabled enabled

d) user clicks pause, then clicks cancel
	re-enable re-enable re-enable
	go to Idle panel

e) user clicks pause, then clicks resume
	Resume button changes to Pause
	re-enable re-enable re-enable


...don't go away until state is "operational"
z) panel goes away amd switches to a different panel
	re-enable re-enable re-enable
	go to Idle panel




don't forget the happy path:
a) panel is displayed
	enabled enabled enabled
b) finish

c) what to display?
	congratulations?


*/