
package ui

import (
	"fmt"
	"time"

	"pifke.org/wpasupplicant"
	"github.com/gotk3/gotk3/gtk"

	// "github.com/Z-Bolt/OctoScreen/interfaces"
	// "github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)


var keyBoardChars = []byte{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '_', '+', '-',
	'=', '.', ',', '|', ':', ';', '/', '~', '`', '[', ']', '{', '}',
	'±', '§', '\\', ' ',
}

type connectToNetworkPanel struct {
	CommonPanel
	pass				*gtk.Entry
	cursorPosition		int
	SSID				string
	SSIDLabel			*gtk.Label
}

var connectToNetworkPanelInstance *connectToNetworkPanel

func GetConnectToNetworkPanelInstance(
	ui					*UI,
	SSID				string,
) *connectToNetworkPanel {
	if connectToNetworkPanelInstance == nil {
		instance := &connectToNetworkPanel {
			CommonPanel:		CreateCommonPanel("ConnectToNetworkPanel", ui),
			cursorPosition:		0,
		}
		instance.initialize()
		connectToNetworkPanelInstance = instance
		connectToNetworkPanelInstance.setSSID(SSID)
	}

	return connectToNetworkPanelInstance
}

func (this *connectToNetworkPanel) initialize() {
	layoutBox := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	layoutBox.SetHExpand(true)

	layoutBox.Add(this.createTopBar())
	layoutBox.Add(this.createKeyboardWindow())
	layoutBox.Add(this.createActionBar())
	this.Grid().Add(layoutBox)
}

func (this *connectToNetworkPanel) setSSID(SSID string) {
	this.SSID = SSID
	str := fmt.Sprintf("Enter password for \"%s\": ", utils.StrEllipsisLen(this.SSID, 18))
	this.SSIDLabel.SetText(str)
}

func (this *connectToNetworkPanel) createTopBar() *gtk.Box {
	this.pass, _ = gtk.EntryNew()
	this.pass.SetProperty("height-request", this.Scaled(40))
	this.pass.SetProperty("width-request", this.Scaled(150))
	this.pass.SetHExpand(true)

	topBar := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	topBar.SetHExpand(true)
	topBar.SetMarginStart(25)
	topBar.SetHAlign(gtk.ALIGN_START)
	this.SSIDLabel = utils.MustLabel(fmt.Sprintf("Pass for %s: ", this.SSID))
	topBar.Add(this.SSIDLabel)
	topBar.Add(this.pass)

	image := utils.MustImageFromFileWithSize("backspace.svg", this.Scaled(40), this.Scaled(40))
	backspaceButton := utils.MustButton(image, func() {
		if this.cursorPosition == 0 {
			return
		}

		this.pass.DeleteText(this.cursorPosition - 1, this.cursorPosition)
		this.cursorPosition--
	})

	topBar.Add(backspaceButton)

	return topBar
}

func (this *connectToNetworkPanel) createActionBar() *gtk.Box {
	image := utils.MustImageFromFileWithSize("back.svg", this.Scaled(40), this.Scaled(40))
	backspaceButton := utils.MustButton(image, func() {
		this.UI.GoToPreviousPanel()
	})

	backspaceButton.SetHAlign(gtk.ALIGN_END)

	actionBar := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	actionBar.SetHAlign(gtk.ALIGN_END)
	actionBar.Add(this.createConnectToNetworkButton())
	actionBar.Add(backspaceButton)

	return actionBar
}

func (this *connectToNetworkPanel) createKeyboardWindow() *gtk.ScrolledWindow {
	keysGrid := utils.MustGrid()
	keysGrid.SetRowHomogeneous(true)
	keysGrid.SetColumnHomogeneous(true)

	keyboardWindow, _ := gtk.ScrolledWindowNew(nil, nil)
	keyboardWindow.SetProperty("overlay-scrolling", false)
	keyboardWindow.SetVExpand(true)
	keyboardWindow.Add(keysGrid)

	row := this.Scaled(3)

	for i, k := range keyBoardChars {
		buttonHander := &keyButtonHander{
			char: k,
			connectToNetworkPanel: this,
		}
		button := utils.MustButtonText(string(k), buttonHander.clicked)
		ctx, _ := button.GetStyleContext()
		ctx.AddClass("keyboard")
		button.SetProperty("height-request", this.Scaled(40))
		keysGrid.Attach(button, i % row, i / row, 1, 1)
	}

	return keyboardWindow
}

func (this *connectToNetworkPanel) createConnectToNetworkButton() *gtk.Button {
	var button *gtk.Button

	button = utils.MustButtonText("Connect", func() {
		button.SetSensitive(false)
		time.Sleep(time.Second * 1)
		psk, _ := this.pass.GetText()
		wpa, _ := wpasupplicant.Unixgram("wlan0")

		if wpa != nil {
			wpa.RemoveAllNetworks()
			wpa.AddNetwork()
			wpa.SetNetwork(0, "ssid", this.SSID)
			wpa.SetNetwork(0, "psk", psk)

			go wpa.EnableNetwork(0)
			time.Sleep(time.Second * 1)
			go wpa.SaveConfig()
		}

		time.Sleep(time.Second * 1)
		this.UI.GoToPreviousPanel()
	})

	ctx, _ := button.GetStyleContext()
	ctx.AddClass("color3")

	button.SetProperty("width-request", this.Scaled(150))

	return button
}



type keyButtonHander struct {
	char					byte
	connectToNetworkPanel	*connectToNetworkPanel
}

func (this *keyButtonHander) clicked() {
	this.connectToNetworkPanel.pass.InsertText(
		string(this.char),
		this.connectToNetworkPanel.cursorPosition,
	)
	this.connectToNetworkPanel.cursorPosition++
}
