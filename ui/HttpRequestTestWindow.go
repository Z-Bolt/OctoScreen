package ui

import (
	"fmt"
	// "io/ioutil"
	// "os"
	// "os/user"
	// "path/filepath"
	// "strconv"
	// "strings"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/gotk3/gotk3/gtk"

	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/ui"
	"github.com/Z-Bolt/OctoScreen/utils"
)


type HttpRequestTestWindow struct {
	Window						*gtk.Window
	Label						*gtk.Label
	Client						*octoprintApis.Client
	BackgroundTask				*utils.BackgroundTask
	UpdateCount					int
}

func CreateHttpRequestTestWindow(
	endpoint string,
	key string,
	width int,
	height int,
) *HttpRequestTestWindow {
	logger.TraceEnter("CreateHttpRequestTestWindow()")

	instance := &HttpRequestTestWindow {
		Window: nil,
		Label: nil,
		Client: octoprintApis.NewClient(endpoint, key),
		BackgroundTask: nil,
		UpdateCount: 0,
	}

	instance.BackgroundTask = utils.CreateBackgroundTask(time.Second * 10, instance.updateTestWindow)

	instance.createWindow(width, height)
	instance.addControls()
	defer instance.Window.ShowAll()

	instance.sdNotify(daemon.SdNotifyReady)

	logger.TraceLeave("CreateHttpRequestTestWindow()")
	return instance
}


func (this *HttpRequestTestWindow) createWindow(
	width int,
	height int,
) {
	window, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		logger.LogFatalError("createWindow()", "WindowNew()", err)
	}

	this.Window = window
	this.Window.SetTitle("HTTP Request Test")
	this.Window.SetDefaultSize(width, height)


	this.Window.Connect("show", this.BackgroundTask.Start)

	this.Window.Connect("destroy", func() {
		logger.Debug("destroy() callback was called")
		gtk.MainQuit()
	})
}


func (this *HttpRequestTestWindow) addControls() {
	// Create a new label widget to show in the window.
	label, err := gtk.LabelNew("")
	if err != nil {
		logger.LogFatalError("CreateHttpRequestTestWindow()", "LabelNew()", err)
	}

	this.Label = label
	this.Label.SetHAlign(gtk.ALIGN_START)
	this.Label.SetVAlign(gtk.ALIGN_START)
	this.Window.Add(this.Label)
}


func (this *HttpRequestTestWindow) updateTestWindow() {
	logger.TraceEnter("updateTestWindow()")

	this.UpdateCount++
	strUpdateCount := fmt.Sprintf("UpdateCount: %d", this.UpdateCount)
	logger.Debug(strUpdateCount)
	this.Label.SetLabel(strUpdateCount)


	//this.checkNotification()
	this.verifyConnection()

	logger.TraceLeave("updateTestWindow()")
}


func (this *HttpRequestTestWindow) verifyConnection() {
	logger.TraceEnter("verifyConnection()")

	this.sdNotify(daemon.SdNotifyWatchdog)

	connectionResponse, err := (&octoprintApis.ConnectionRequest{}).Do(this.Client)
	if err != nil {
		logger.LogError("verifyConnection()", "ConnectionRequest.Do()", err)
	} else {
		strCurrentState := string(connectionResponse.Current.State)
		logger.Debugf("verifyConnection() succeeded")
		logger.Debugf("connectionResponse.Current.State is %q", strCurrentState)
	}

	logger.TraceLeave("verifyConnection()")
}


func (this *HttpRequestTestWindow) sdNotify(state string) {
	logger.TraceEnter("sdNotify()")

	_, err := daemon.SdNotify(false, state)
	if err != nil {
		logger.LogError("sdNotify()", "daemon.SdNotify()", err)
		logger.TraceLeave("sdNotify()")
		return
	}

	logger.TraceLeave("sdNotify()")
}
