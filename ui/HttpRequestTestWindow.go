package ui

import (
	"fmt"
	// "io/ioutil"
	// standardLog "log"
	// "os"
	// "os/user"
	// "path/filepath"
	// "strconv"
	// "strings"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/gotk3/gotk3/gtk"
	// "github.com/sirupsen/logrus"
	// "github.com/Z-Bolt/OctoScreen/ui"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/utils"
	// "gopkg.in/yaml.v1"
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
	utils.Logger.Debug("entering CreateHttpRequestTestWindow()")


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


	instance.sdNotify("READY=1")
	//instance.sdNotify(daemon.SdNotifyReady)


	utils.Logger.Debug("leaving CreateHttpRequestTestWindow()")

	return instance
}


func (this *HttpRequestTestWindow) createWindow(
	width int,
	height int,
) {
	window, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		utils.LogFatalError("createWindow()", "WindowNew()", err)
	}

	this.Window = window
	this.Window.SetTitle("HTTP Request Test")
	this.Window.SetDefaultSize(width, height)


	this.Window.Connect("show", this.BackgroundTask.Start)

	this.Window.Connect("destroy", func() {
		utils.Logger.Debug("destroy() callback was called")
		gtk.MainQuit()
	})
}


func (this *HttpRequestTestWindow) addControls() {
	// Create a new label widget to show in the window.
	label, err := gtk.LabelNew("")
	if err != nil {
		utils.LogFatalError("CreateHttpRequestTestWindow()", "LabelNew()", err)
	}

	this.Label = label
	this.Label.SetHAlign(gtk.ALIGN_START)
	this.Label.SetVAlign(gtk.ALIGN_START)
	this.Window.Add(this.Label)
}


func (this *HttpRequestTestWindow) updateTestWindow() {
	utils.Logger.Debug("entering updateTestWindow()")


	this.UpdateCount++
	strUpdateCount := fmt.Sprintf("UpdateCount: %d", this.UpdateCount)
	utils.Logger.Debug(strUpdateCount)
	this.Label.SetLabel(strUpdateCount)


	//this.checkNotification()
	this.verifyConnection()

	utils.Logger.Debug("leaving updateTestWindow()")
}


func (this *HttpRequestTestWindow) verifyConnection() {
	utils.Logger.Debug("entering verifyConnection()")

	this.sdNotify("WATCHDOG=1")
	//this.sdNotify(daemon.SdNotifyWatchdog)

	connectionResponse, err := (&octoprintApis.ConnectionRequest{}).Do(this.Client)
	if err != nil {
		utils.LogError("verifyConnection()", "ConnectionRequest.Do()", err)
	} else {
		strCurrentState := string(connectionResponse.Current.State)
		utils.Logger.Debugf("    verifyConnection() succeeded")
		utils.Logger.Debugf("    connectionResponse.Current.State is %q", strCurrentState)
	}

	utils.Logger.Debug("leaving verifyConnection()")
}


func (this *HttpRequestTestWindow) sdNotify(state string) {
	utils.Logger.Debug("entering sdNotify()")

	_, err := daemon.SdNotify(false, state)
	if err != nil {
		utils.Logger.Errorf("sdNotify()", "SdNotify()", err)
		utils.Logger.Debug("leaving sdNotify()")
		return
	}

	utils.Logger.Debug("leaving sdNotify()")
}
