package ui

import (
	"fmt"
	// "math"
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

	PrintWasCanceled bool

	// Tools
	tool0Button				*uiWidgets.ToolPrintingButton
	tool1Button				*uiWidgets.ToolPrintingButton
	tool2Button				*uiWidgets.ToolPrintingButton
	tool3Button				*uiWidgets.ToolPrintingButton
	tool4Button				*uiWidgets.ToolPrintingButton
	bedButton  				*uiWidgets.ToolPrintingButton

	// Statistics/Info
	fileLabelWithImage		*utils.LabelWithImage
	timeLabelWithImage		*utils.LabelWithImage
	timeLeftLabelWithImage	*utils.LabelWithImage
	// layerLabelWithImage	*utils.LabelWithImage
	// The info for the current / total layers is not available
	// See https://community.octoprint.org/t/layer-number-and-total-layers-from-api/8005/4
	// and https://docs.octoprint.org/en/master/api/datamodel.html#sec-api-datamodel-jobs-job
	// Darn.
	
	// Progress
	progressBar				*gtk.ProgressBar

	// Toolbar buttons
	pauseButton				*gtk.Button
	cancelButton			*gtk.Button
	controlButton			*gtk.Button
	completedButton			*gtk.Button

	backgroundTask			*utils.BackgroundTask
}

var printStatusPanelInstance *printStatusPanel

func getPrintStatusPanelInstance(ui *UI) *printStatusPanel {
	if printStatusPanelInstance == nil {
		printStatusPanelInstance = &printStatusPanel {
			CommonPanel: CreateTopLevelCommonPanel("PrintStatusPanel", ui),
			PrintWasCanceled: false,
		}
		printStatusPanelInstance.initialize()
		printStatusPanelInstance.createBackgroundTask()
	}

	return printStatusPanelInstance
}

func GoToPrintStatusPanel(ui *UI) {
	instance := getPrintStatusPanelInstance(ui)
	instance.progressBar.SetText("0%")
	ui.GoToPanel(instance)
}

func (this *printStatusPanel) initialize() {
	defer this.Initialize()

	this.createToolButtons()
	this.Grid().Attach(this.createInfoBox(),        2, 0, 2, 1)
	this.Grid().Attach(this.createProgressBar(),    2, 1, 2, 1)
	this.createToolBarButtons()
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
	this.tool4Button = uiWidgets.CreateToolPrintingButton( 5)
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

		case 5:
			this.Grid().Attach(this.tool0Button, 0, 0, 1, 1)
			this.Grid().Attach(this.tool1Button, 1, 0, 1, 1)
			this.Grid().Attach(this.tool2Button, 0, 1, 1, 1)
			this.Grid().Attach(this.tool3Button, 1, 1, 1, 1)
			// this.Grid().Attach(this.tool4Button, 0, 2, 1, 1)
			// this.Grid().Attach(this.bedButton,   1, 2, 1, 1)
			// ...there's not enough sceen realestate for the 5th toolhead button,
			// ...so use the same layout as the "4 toolhead" version:
			this.Grid().Attach(this.bedButton,   0, 2, 1, 1)
	}
}

func (this *printStatusPanel) createInfoBox() *gtk.Box {
	this.fileLabelWithImage = utils.MustLabelWithImage("file-gcode.svg", "")
	ctx, _ := this.fileLabelWithImage.GetStyleContext()
	ctx.AddClass("printing-status-label")

	this.timeLabelWithImage = utils.MustLabelWithImage("time.svg", "Print time:")
	ctx, _ = this.timeLabelWithImage.GetStyleContext()
	ctx.AddClass("printing-status-label")

	this.timeLeftLabelWithImage = utils.MustLabelWithImage("time.svg", "Print time left:")
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

func (this *printStatusPanel) createToolBarButtons() {
	this.pauseButton = utils.MustButtonImageUsingFilePath(
		"Pause",
		"pause.svg",
		this.handlePauseClicked,
	)
	this.Grid().Attach(this.pauseButton,    1, 2, 1, 1)

	this.cancelButton = utils.MustButtonImageStyle(
		"Cancel",
		"stop.svg",
		"color-warning-sign-yellow",
		this.handleCancelClicked,
	)
	this.Grid().Attach(this.cancelButton,   2, 2, 1, 1)

	this.controlButton = utils.MustButtonImageStyle(
		"Control",
		"printing-control.svg",
		"color3",
		this.handleControlClicked,
	)
	this.Grid().Attach(this.controlButton,  3, 2, 1, 1)

	this.completedButton = utils.MustButtonImageStyle(
		"Completed",
		"complete.svg",
		"color3",
		this.handleCompleteClicked,
	)
	this.Grid().Attach(this.completedButton, 1, 2, 3, 1)
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

	octoPrintResponseManager := GetOctoPrintResponseManagerInstance(this.UI)
	if octoPrintResponseManager.IsConnected() != true {
		// If not connected, do nothing and leave.
		logger.Debugf("PrintStatusPanel.update() - not connected, now exiting")
		logger.TraceLeave("PrintStatusPanel.update()")
		return
	}

	jobResponse, err := (&octoprintApis.JobRequest{}).Do(this.UI.Client)
	if err != nil {
		logger.LogError("PrintStatusPanel.update()", "Do(JobRequest)", err)
		logger.TraceLeave("PrintStatusPanel.update()")
		return
	}

	logger.Debugf("PrintStatusPanel.update() - jobResponse.State is %s", jobResponse.State)

	this.updateStates(jobResponse)
	this.updateToolTemperatures(&octoPrintResponseManager.FullStateResponse.Temperature)
	this.updateInfoBox(jobResponse)
	this.updateProgress(jobResponse)
	this.updateToolBarButtons(jobResponse)

	logger.TraceLeave("PrintStatusPanel.update()")
}

func (this *printStatusPanel) updateStates(jobResponse *dataModels.JobResponse) {
	if jobResponse.State == "Cancelling" {
		this.PrintWasCanceled = true
	}

	if jobResponse.State == "Printing" || jobResponse.State == "Operational" {
		if jobResponse.Progress.PrintTimeLeft <= 0.0 {
			// Special case for handling the buttons
			this.pauseButton.SetSensitive(false)
			this.pauseButton.Hide()
	
			this.cancelButton.SetSensitive(false)
			this.cancelButton.Hide()
	
			this.controlButton.SetSensitive(false)
			this.controlButton.Hide()

			this.completedButton.Show()
			this.completedButton.SetSensitive(true)

			this.progressBar.Hide();

			this.timeLeftLabelWithImage.Label.SetLabel("Print time left: 00:00:00")
		}
	}
}

func (this *printStatusPanel) updateToolTemperatures(temperature *dataModels.TemperatureStateResponse) {
	logger.TraceEnter("PrintStatusPanel.updateToolTemperatures()")

	for tool, currentTemperatureData := range temperature.CurrentTemperatureData {
		text := utils.GetTemperatureDataString(currentTemperatureData)
		switch tool {
			case "bed":
				logger.Debug("Updating the UI's bed temp")
				this.bedButton.SetLabel(text)

			case "tool0":
				logger.Debug("Updating the UI's tool0 temp")
				this.tool0Button.SetLabel(text)

			case "tool1":
				logger.Debug("Updating the UI's tool1 temp")
				this.tool1Button.SetLabel(text)

			case "tool2":
				logger.Debug("Updating the UI's tool2 temp")
				this.tool2Button.SetLabel(text)

			case "tool3":
				logger.Debug("Updating the UI's tool3 temp")
				this.tool3Button.SetLabel(text)

			case "tool4":
				logger.Debug("Updating the UI's tool4 temp")
				this.tool4Button.SetLabel(text)

			default:
				logger.Errorf("PrintStatusPanel.updateToolTemperatures() - GetOctoPrintResponseManagerInstance() returned an unknown tool: %q", tool)
		}
	}

	logger.TraceLeave("PrintStatusPanel.updateToolTemperatures()")
}

func (this *printStatusPanel) updateInfoBox(jobResponse *dataModels.JobResponse) {
	logger.TraceEnter("PrintStatusPanel.updateInfoBox()")

	if jobResponse.State != "Printing" {
		logger.TraceLeave("PrintStatusPanel.updateInfoBox()")
		return;
	}

	jobFileName := "<i>not-set</i>"
	if jobResponse.Job.File.Name != "" {
		jobFileName = jobResponse.Job.File.Name
		jobFileName = strings.Replace(jobFileName, ".gcode", "", -1)
		jobFileName = strings.Replace(jobFileName, ".gco", "", -1)
		jobFileName = utils.TruncateString(jobFileName, 20)
	}

	this.fileLabelWithImage.Label.SetLabel(jobFileName)

	var timeSpent string
	var timeLeft string
	switch jobResponse.Progress.Completion {
		case 0:
			timeSpent = "Warming up ..."
			timeLeft = ""
			break

		case 100:
			timeSpent = fmt.Sprintf("Completed in %s", time.Duration(int64(jobResponse.Job.LastPrintTime) * 1e9))
			timeLeft = ""
			break

		default:
			logger.Info(jobResponse.Progress.PrintTime)
			printTimeDuration := time.Duration(int64(jobResponse.Progress.PrintTime) * 1e9)
			printTimeLeftDuration := time.Duration(int64(jobResponse.Progress.PrintTimeLeft) * 1e9)
			timeSpent = fmt.Sprintf("Print time: %s", formattedDuration(printTimeDuration))
			timeLeft = fmt.Sprintf("Print time left: %s", formattedDuration(printTimeLeftDuration))
			break
	}

	this.timeLabelWithImage.Label.SetLabel(timeSpent)
	this.timeLeftLabelWithImage.Label.SetLabel(timeLeft)

	logger.TraceLeave("PrintStatusPanel.updateInfoBox()")
}

func (this *printStatusPanel) updateProgress(jobResponse *dataModels.JobResponse) {
	logger.TraceEnter("PrintStatusPanel.updateProgress()")
	logger.Debugf("PrintStatusPanel.updateProgress() - jobResponse.State is: '%s'", jobResponse.State)

	if jobResponse.State != "Printing" {
		this.progressBar.SetText(jobResponse.State)
		logger.TraceLeave("PrintStatusPanel.updateProgress()")
		return;
	}

	// Use the following for the percentage of time taken
	if jobResponse.Progress.PrintTime <= 0.0 {
		// logger.Debugf("Error! jobResponse.Progress.PrintTime is: %f", jobResponse.Progress.PrintTime)
		// this.progressBar.SetText("0%")
		this.progressBar.SetText("Warming up")
		logger.TraceLeave("PrintStatusPanel.updateProgress()")
		return
	}

	progresBarFraction := jobResponse.Progress.PrintTime / (jobResponse.Progress.PrintTime + jobResponse.Progress.PrintTimeLeft)
	this.progressBar.SetFraction(progresBarFraction)
	progresssPercentage := int64(progresBarFraction * 100.0)
	this.progressBar.SetText(fmt.Sprintf("%d%%", progresssPercentage))

	logger.TraceLeave("PrintStatusPanel.updateProgress()")
}

func (this *printStatusPanel) updateToolBarButtons(jobResponse *dataModels.JobResponse) {
	logger.TraceEnter("PrintStatusPanel.updateToolBarButtons()")
	logger.Debugf("PrintStatusPanel.updateToolBarButtons() - jobResponse.State is: '%s'", jobResponse.State)

	switch jobResponse.State {
		case "Printing":
			this.pauseButton.SetLabel("Pause")
			pauseImage := utils.MustImageFromFile("pause.svg")
			this.pauseButton.SetImage(pauseImage)
			this.pauseButton.SetSensitive(true)
			this.pauseButton.Show()

			this.cancelButton.SetSensitive(true)
			this.cancelButton.Show()

			this.controlButton.SetSensitive(true)
			this.controlButton.Show()

			this.completedButton.SetSensitive(false)
			this.completedButton.Hide()
			break;

		case "Pausing":
			this.pauseButton.SetSensitive(false)
			this.pauseButton.Show()

			this.cancelButton.SetSensitive(false)
			this.cancelButton.Show()

			this.controlButton.SetSensitive(true)
			this.controlButton.Show()

			this.completedButton.SetSensitive(false)
			this.completedButton.Hide()
			break;

		case "Paused":
			this.pauseButton.SetLabel("Resume")
			pauseImage := utils.MustImageFromFile("resume.svg")
			this.pauseButton.SetImage(pauseImage)
			this.pauseButton.SetSensitive(true)
			this.pauseButton.Show()

			this.cancelButton.SetSensitive(true)
			this.cancelButton.Show()

			this.controlButton.SetSensitive(true)
			this.controlButton.Show()

			this.completedButton.SetSensitive(false)
			this.completedButton.Hide()
			break;

		case "Cancelling":
			this.pauseButton.SetSensitive(false)
			this.pauseButton.Show()

			this.cancelButton.SetSensitive(false)
			this.cancelButton.Show()

			this.controlButton.SetSensitive(true)
			this.controlButton.Show()

			this.completedButton.SetSensitive(false)
			this.completedButton.Hide()
			break;

		case "Finishing":
			break;

		case "Operational":
			break;

		default:
			logLevel := logger.LogLevel()
			if logLevel == "debug" {
				logger.Debugf("PrintStatusPanel.updateToolBarButtons() - unknown jobResponse.State: '%s'", jobResponse.State)
				logger.Panicf("PrintStatusPanel.updateToolBarButtons() - unknown jobResponse.State: '%s'", jobResponse.State)
			}
	}

	logger.TraceLeave("PrintStatusPanel.updateToolBarButtons()")
}

func (this *printStatusPanel) handlePauseClicked() {
	logger.TraceEnter("PrintStatusPanel.handlePauseClicked()")

	// TODO: is this needed?
	// defer this.updateTemperature()

	cmd := &octoprintApis.PauseRequest{Action: dataModels.Toggle}
	err := cmd.Do(this.UI.Client)
	if err != nil {
		logger.LogError("PrintStatusPanel.handlePauseClicked()", "Do(PauseRequest)", err)
		logger.TraceLeave("PrintStatusPanel.handlePauseClicked()")
		return
	}

	logger.TraceLeave("PrintStatusPanel.handlePauseClicked()")
}

func (this *printStatusPanel) handleCancelClicked() {
	userResponse := this.confirmCancelDialogBox(
		this.UI.window,
		"Are you sure you want to cancel the current print?",
	)

	if userResponse == int(gtk.RESPONSE_YES) {
		this.cancelPrintJob()
	}
}

func (this *printStatusPanel) handleControlClicked() {
	this.UI.GoToPanel(GetPrintMenuPanelInstance(this.UI))
}

func (this *printStatusPanel) handleCompleteClicked() {
	this.UI.WaitingForUserToContinue = false
}

func (this *printStatusPanel) confirmCancelDialogBox(
	parentWindow		*gtk.Window,
	message				string,
) int {
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

	return userResponse
}

func (this *printStatusPanel) cancelPrintJob() {
	logger.TraceEnter("PrintStatusPanel.cancelPrintJob()")

	err := (&octoprintApis.CancelRequest{}).Do(this.UI.Client)
	if err == nil {
		this.pauseButton.SetSensitive(false)
		this.cancelButton.SetSensitive(false)
		this.controlButton.SetSensitive(false)
	} else {
		logger.LogError("PrintStatusPanel.cancelPrintJob()", "Do(CancelRequest)", err)
	}

	logger.TraceLeave("PrintStatusPanel.cancelPrintJob()")
}

func formattedDuration(duration time.Duration) string {
	hours := duration / time.Hour
	duration -= hours * time.Hour

	minutes := duration / time.Minute
	duration -= minutes * time.Minute

	seconds := duration / time.Second

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
