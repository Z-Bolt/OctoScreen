package ui

import (
	"fmt"
	// "os"
	// "strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-collections/collections/stack"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)


type UI struct {
	sync.Mutex

	PanelHistory				*stack.Stack
	Client						*octoprintApis.Client
	ConnectionState				dataModels.ConnectionState
	Settings					*dataModels.OctoScreenSettingsResponse

	UiState						UiState
	WaitingForUserToContinue	bool

	OctoPrintPluginIsAvailable	bool
	NotificationsBox			*uiWidgets.NotificationsBox
	grid						*gtk.Grid
	window						*gtk.Window
	time						time.Time

	width						int
	height						int
	scaleFactor					int
	connectionAttempts			int
}

func CreateUi() *UI {
	logger.TraceEnter("ui.CreateUi()")

	octoScreenConfig := utils.GetOctoScreenConfigInstance()
	endpoint := octoScreenConfig.OctoPrintConfig.Server.Host
	key := octoScreenConfig.OctoPrintConfig.API.Key
	width := octoScreenConfig.Width
	height := octoScreenConfig.Height

	if width == 0 {
		panic("the window's width was not specified")
	}

	if height == 0 {
		panic("the window's height was not specified")
	}

	instance := &UI {
		PanelHistory:				stack.New(),
		Client:						octoprintApis.NewClient(endpoint, key),
		NotificationsBox:			uiWidgets.NewNotificationsBox(),
		OctoPrintPluginIsAvailable:	false,
		Settings:					nil,
		UiState:					Uninitialized,
		WaitingForUserToContinue:	false,
		window:						utils.MustWindow(gtk.WINDOW_TOPLEVEL),
		time:						time.Now(),
		width:						width,
		height:						height,
	}

	instance.initialize1()

	logger.TraceLeave("ui.CreateUi()")

	return instance
}

// TODO: rename initialize1() and initialize2()
func (this *UI) initialize1() {
	logger.TraceEnter("ui.initialize1()")

	this.window.Connect("configure-event", func(win *gtk.Window) {
		allocatedWidth:= win.GetAllocatedWidth()
		allocatedHeight:= win.GetAllocatedHeight()
		sizeWidth, sizeHeight := win.GetSize()

		if (allocatedWidth > this.width || allocatedHeight > this.height) ||
			(sizeWidth > this.width || sizeHeight > this.height) {
			logger.Errorf(
				"Window resize went past max size.  allocatedWidth:%d allocatedHeight:%d sizeWidth:%d sizeHeight:%d",
				allocatedWidth,
				allocatedHeight,
				sizeWidth,
				sizeHeight,
			)
			logger.Errorf(
				"Window resize went past max size.  Target width and height: %dx%d",
				this.width,
				this.height,
			)
		}
	})

	switch {
		case this.width > 480:
			this.scaleFactor = 2

		case this.width > 1000:
			this.scaleFactor = 3

		default:
			this.scaleFactor = 1
	}

	this.initialize2()
	GoToConnectionPanel(this)

	logger.TraceLeave("ui.initialize1()")
}

func (this *UI) initialize2() {
	logger.TraceEnter("ui.initialize2()")

	defer this.window.ShowAll()
	this.loadStyle()

	this.window.SetTitle(utils.WindowName)
	this.window.SetDefaultSize(this.width, this.height)
	this.window.SetResizable(false)

	this.window.Connect("destroy", func() {
		logger.Debug("window destroy callback was called, now executing MainQuit()")
		gtk.MainQuit()
	})

	overlay := utils.MustOverlay()
	this.window.Add(overlay)

	this.grid = utils.MustGrid()
	overlay.Add(this.grid)

	GetOctoPrintResponseManagerInstance(this)

	logger.TraceLeave("ui.initialize2()")
}

func (this *UI) loadStyle() {
	logger.TraceEnter("ui.loadStyle()")

	cssProvider := utils.MustCssProviderFromFile(utils.CssFileName)

	screenDefault, err := gdk.ScreenGetDefault()
	if err != nil {
		logger.LogError("ui.loadStyle()", "ScreenGetDefault()", err)
		logger.TraceLeave("ui.loadStyle()")
		return
	}

	gtk.AddProviderForScreen(screenDefault, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_USER)

	logger.TraceLeave("ui.loadStyle()")
}

func (this *UI) Update() {
	logger.TraceEnter("ui.Update()")
	connectionManager := utils.GetConnectionManagerInstance(this.Client)
	if connectionManager.IsConnectedToOctoPrint == true {
		if this.Settings == nil {
			this.loadSettings()
		}
	}

	logger.TraceLeave("ui.Update()")
}


func (this *UI) loadSettings() {
	logger.TraceEnter("ui.loadSettings()")

	if this.Settings != nil {
		logger.Error("ui.loadSettings() - this.Settings has already been set")
		logger.TraceLeave("ui.loadSettings()")
		return
	}

	settingsResponse, err := (&octoprintApis.OctoScreenSettingsRequest{}).Do(this.Client)
	if err != nil {
		text := err.Error()
		if strings.Contains(strings.ToLower(text), "unexpected status code: 404") {
			// The call to GetSettings is also used to determine whether or not the
			// OctoScreen plug-in is available.  If calling GetSettings returns
			// a 404, the plug-in isn't available.
			logger.Info("The OctoScreen plug-in is not available")
		} else {
			// If we get back any other kind of error, something bad happened, so log an error.
			logger.LogError("ui.loadSettings()", "Do(GetSettingsRequest)", err)
		}

		this.OctoPrintPluginIsAvailable = false
		// Use default settings
		this.Settings = &dataModels.OctoScreenSettingsResponse {
			FilamentInLength: 100,
			FilamentOutLength: 100,
			ToolChanger: false,
			XAxisInverted: false,
			YAxisInverted: false,
			ZAxisInverted: false,
			MenuStructure: nil,
		}
	} else {
		logger.Info("The call to GetSettings succeeded and the OctoPrint plug-in is available")
		this.OctoPrintPluginIsAvailable = true

		if !this.validateMenuItems(settingsResponse.MenuStructure, "", true) {
			settingsResponse.MenuStructure = nil
		}

		this.Settings = settingsResponse
	}

	logger.TraceLeave("ui.loadSettings()")
}

func (this *UI) validateMenuItems(menuItems []dataModels.MenuItem, name string, isRoot bool) bool {
	logger.TraceEnter("ui.validateMenuItems()")

	if menuItems == nil {
		logger.TraceLeave("ui.validateMenuItems()")
		return true
	}

	maxCount := 11
	if isRoot {
		maxCount = 4
	}

	menuItemsLength := len(menuItems)
	if menuItemsLength > maxCount {
		message := ""
		description := ""
		if isRoot {
			message = fmt.Sprintf("Error!  The custom menu structure can only have %d items\n    at the root level (the idle panel).", maxCount)
			description = fmt.Sprintf("\n    When the MenuStructure was parsed, %d items were found.", menuItemsLength)
		} else {
			message = fmt.Sprintf("Error!  A panel can only have a maximum of %d items.", maxCount)
			description = fmt.Sprintf("\n    When the MenuStructure for '%s' was parsed,\n    %d items were found.", name, menuItemsLength)
		}

		fatalErrorWindow := CreateFatalErrorWindow(
			message,
			description,
		)
		fatalErrorWindow.ShowAll()

		logger.TraceLeave("ui.validateMenuItems()")
		return false
	}

	for i := 0; i < len(menuItems); i++ {
		menuItem := menuItems[i]
		if menuItem.Panel == "menu" {
			if !this.validateMenuItems(menuItem.Items, menuItem.Name, false) {
				logger.TraceLeave("ui.validateMenuItems()")
				return false
			}
		}
	}

	logger.TraceLeave("ui.validateMenuItems()")
	return true
}

func (this *UI) GoToPanel(panel interfaces.IPanel) {
	logger.TraceEnter("ui.GoToPanel()")

	panelName := panel.Name()
	logger.Debugf("ui.GoToPanel() - panel name is %s", panelName)

	this.SetUiToPanel(panel)
	this.PanelHistory.Push(panel)

	logger.TraceLeave("ui.GoToPanel()")
}

func (this *UI) GoToPreviousPanel() {
	logger.TraceEnter("ui.GoToPreviousPanel()")

	stackLength := this.PanelHistory.Len()
	if stackLength < 2 {
		logger.Error("ui.GoToPreviousPanel() - stack does not contain current panel and parent panel")
		logger.TraceLeave("ui.GoToPreviousPanel()")
		return
	}

	if stackLength < 1 {
		logger.Error("ui.GoToPreviousPanel() - GoToPreviousPanel() was called but the stack is empty")
		logger.TraceLeave("ui.GoToPreviousPanel()")
		return
	}

	currentPanel := this.PanelHistory.Pop().(interfaces.IPanel)
	this.RemovePanelFromUi(currentPanel)

	parentPanel := this.PanelHistory.Peek().(interfaces.IPanel)
	this.SetUiToPanel(parentPanel)

	logger.TraceLeave("ui.GoToPreviousPanel()")
}

func (this *UI) GetCurrentPanel() interfaces.IPanel {
	logger.TraceEnter("ui.GetCurrentPanel()")

	currentPanel := interfaces.IPanel(nil)

	stackLength := this.PanelHistory.Len()
	if stackLength > 0 {
		currentPanel = this.PanelHistory.Peek().(interfaces.IPanel)
	} else {
		logger.Error("ui.GetCurrentPanel() was called, but PanelHistory is empty")
	}

	logger.TraceLeave("ui.GetCurrentPanel()")

	return currentPanel
}

func (this *UI) SetUiToPanel(panel interfaces.IPanel) {
	logger.TraceEnter("ui.SetUiToPanel()")

	logger.Infof("Setting panel to %q", panel.Name())

	stackLength := this.PanelHistory.Len()
	if stackLength > 0 {
		currentPanel := this.GetCurrentPanel()
		this.RemovePanelFromUi(currentPanel)
	}

	panel.PreShow()
	panel.Show()
	this.grid.Attach(panel.Grid(), 0, 0, 1, 1)
	this.grid.ShowAll()

	logger.TraceLeave("ui.SetUiToPanel()")
}

func (this *UI) RemovePanelFromUi(panel interfaces.IPanel) {
	logger.TraceEnter("ui.RemovePanelFromUi()")

	defer panel.Hide()
	this.grid.Remove(panel.Grid())

	logger.TraceLeave("ui.RemovePanelFromUi()")
}

func (this *UI) errToUser(err error) string {
	logger.TraceEnter("ui.errToUser()")

	text := strings.ToLower(err.Error())
	if strings.Contains(text, "connection refused") {
		logger.TraceLeave("ui.errToUser() - connection refused")
		return "Unable to connect to OctoPrint, check if it is running."
	} else if strings.Contains(text, "request canceled") {
		logger.TraceLeave("ui.errToUser() - request canceled")
		return "Loading..."
	} else if strings.Contains(text, "connection broken") {
		logger.TraceLeave("ui.errToUser() - connection broken")
		return "Loading..."
	}

	msg := fmt.Sprintf("ui.errToUser() - unexpected error: %s", text)
	logger.TraceLeave(msg)
	return fmt.Sprintf("Unexpected Error: %s", text)
}
