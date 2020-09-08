package ui

import (
	"fmt"
	"net"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"pifke.org/wpasupplicant"
	"github.com/Z-Bolt/OctoScreen/interfaces"
	// "github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)

var networkPanelInstance *networkPanel

type networkPanel struct {
	CommonPanel
	listBox				*gtk.Box
	netStatus			*gtk.Label
	wifiStatus			*gtk.Label
}

func NetworkPanel(
	ui				*UI,
	parentPanel		interfaces.IPanel,
) *networkPanel {
	if networkPanelInstance == nil {
		instance := &networkPanel {
			CommonPanel: NewCommonPanel(ui, parentPanel),
		}
		instance.initialize()
		instance.backgroundTask = utils.CreateBackgroundTask(time.Second * 3, instance.update)
		networkPanelInstance = instance
	} else {
		networkPanelInstance.parentPanel = parentPanel
	}

	return networkPanelInstance
}

func (this *networkPanel) update() {
	utils.EmptyTheContainer(&this.listBox.Container)

	netStatus := ""
	addresses, _ := net.InterfaceAddrs()

	for _, address := range addresses {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				netStatus += fmt.Sprintf("IP Address: %s\n", ipNet.IP.String())
			}
		}
	}

	this.netStatus.SetText(netStatus)

	wpa, _ := wpasupplicant.Unixgram("wlan0")
	if wpa != nil {
		status, _ := wpa.Status()
		wifiStatus := ""

		if status.WPAState() == "COMPLETED" {
			wifiStatus = fmt.Sprintf("Wifi Information\nSSID: %s\nIP Address: %s\n", status.SSID(), status.IPAddr())
		} else {
			wifiStatus = fmt.Sprintf("Wifi status is: %s\n", status.WPAState())
		}

		this.wifiStatus.SetText(wifiStatus)

		result, _ := wpa.ScanResults()
		for _, bss := range result {
			this.addNetwork(this.listBox, bss.SSID())
		}

		err := wpa.Scan()
		if err != nil {
			utils.LogError("NetworkPanel.update()", "Scan()", err)
		}
	} else {
		label := utils.MustLabel("\n\nWifi management is not available\non this hardware")
		this.listBox.Add(label)
	}

	this.listBox.ShowAll()
}

func (this *networkPanel) initialize() {
	this.Grid().Attach(this.createNetworkListWindow(), 0, 0, 4, 1)
	this.Grid().Attach(this.createLeftBar(), 4, 0, 3, 1)
}

func (this *networkPanel) createNetworkListWindow() gtk.IWidget {
	this.listBox = utils.MustBox(gtk.ORIENTATION_VERTICAL, 0)
	this.listBox.SetVExpand(true)

	networkListWindow, _ := gtk.ScrolledWindowNew(nil, nil)
	networkListWindow.SetProperty("overlay-scrolling", false)
	networkListWindow.Add(this.listBox)

	return networkListWindow
}

func (this *networkPanel) addNetwork(box *gtk.Box, ssid string) {
	frame, _ := gtk.FrameNew("")

	clicked := func() {
		this.UI.Add(ConnectionPanel(this.UI, this, ssid))
	}

	image := utils.MustImageFromFileWithSize("network.svg", this.Scaled(25), this.Scaled(25))
	button := utils.MustButton(image, clicked)

	name := utils.MustButtonText(utils.StrEllipsisLen(ssid, 18), clicked)
	name.SetHExpand(true)

	network := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 0)
	network.SetHExpand(true)

	network.Add(name)
	network.Add(button)

	frame.Add(network)
	box.Add(frame)
}

func (this *networkPanel) createLeftBar() gtk.IWidget {
	leftBar := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	leftBar.SetHExpand(true)

	leftBar.Add(this.createInfoBar())
	leftBar.Add(this.createActionBar())

	return leftBar
}

func (this *networkPanel) createActionBar() gtk.IWidget {
	layout := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	layout.SetHAlign(gtk.ALIGN_END)
	layout.SetHExpand(true)

	// NOTE: If a message is logged that the image (SVG) can't be loaded, try installing librsvg.
	backImage := utils.MustImageFromFileWithSize("back.svg", this.Scaled(40), this.Scaled(40))
	backButton := utils.MustButton(backImage, func() {
		this.UI.GoHistory()
	})

	layout.Add(backButton)

	return layout
}

func (this *networkPanel) createInfoBar() gtk.IWidget {
	infoBar := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	infoBar.SetVExpand(true)
	infoBar.SetHExpand(true)
	infoBar.SetVAlign(gtk.ALIGN_START)
	infoBar.SetHAlign(gtk.ALIGN_START)

	infoBar.SetMarginEnd(25)
	infoBar.SetMarginStart(25)

	t1 := utils.MustLabel("<b>Network Information</b>")
	t1.SetHAlign(gtk.ALIGN_START)
	t1.SetMarginTop(25)

	infoBar.Add(t1)

	this.netStatus = utils.MustLabel("")
	this.netStatus.SetHAlign(gtk.ALIGN_START)
	infoBar.Add(this.netStatus)

	this.wifiStatus = utils.MustLabel("")
	this.wifiStatus.SetHAlign(gtk.ALIGN_START)
	this.wifiStatus.SetMarginTop(25)
	infoBar.Add(this.wifiStatus)

	return infoBar
}
