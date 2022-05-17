package ui

import (
	"fmt"
	"net"
	// "os"
	// "strconv"
	// "time"

	"pifke.org/wpasupplicant"
	"github.com/gotk3/gotk3/gtk"

	// "github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/logger"
	// "github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)


type networkPanel struct {
	CommonPanel
	listBox					*gtk.Box
	netStatus				*gtk.Label
	wifiStatus				*gtk.Label
	overrideForDebugging	bool
}

var networkPanelInstance *networkPanel

func GetNetworkPanelInstance(
	ui				*UI,
) *networkPanel {
	if networkPanelInstance == nil {
		instance := &networkPanel {
			CommonPanel: CreateCommonPanel("NetworkPanel", ui),
		}
		instance.initialize()

		networkPanelInstance = instance
	}

	return networkPanelInstance
}

func (this *networkPanel) initialize() {
	this.Grid().Attach(this.createNetworkListWindow(), 0, 0, 4, 1)
	this.Grid().Attach(this.createLeftBar(), 4, 0, 3, 1)

	// TODO: make sure overrideForDebugging is set to false before checking in.
	this.overrideForDebugging = false;
}

func (this *networkPanel) update() {
	logger.TraceEnter("NetworkPanel.update()")

	utils.EmptyTheContainer(&this.listBox.Container)
	this.setNetStatusText()
	this.setNetworkItems()
	this.listBox.ShowAll()

	logger.TraceLeave("NetworkPanel.update()")
}

func (this *networkPanel) setNetStatusText() {
	netStatus := ""
	addresses, _ := net.InterfaceAddrs()

	if this.overrideForDebugging {
		netStatus += fmt.Sprintf("IP Address: %s\n", "111.222.333.444")
	} else {
		for _, address := range addresses {
			if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					netStatus += fmt.Sprintf("IP Address: %s\n", ipNet.IP.String())
				}
			}
		}
	}

	this.netStatus.SetText(netStatus)
}

func (this *networkPanel) setNetworkItems() {
	var wpa wpasupplicant.Conn

	if this.overrideForDebugging {
		wpa = nil
	} else {
		wpa, _ = wpasupplicant.Unixgram("wlan0")
	}

	this.setWiFiStatusText(wpa)
	this.setNetworkListItems(wpa)
}

func (this *networkPanel) setWiFiStatusText(wpa wpasupplicant.Conn) {
	wifiStatus := ""

	if this.overrideForDebugging {
		wifiStatus = fmt.Sprintf(
			"WiFi Information\nSSID: %s\nIP Address: %s\n",
			"FooNet 2.4G",
			"111.222.333.444",
		)
	} else {
		if wpa != nil {
			status, _ := wpa.Status()

			if status.WPAState() == "COMPLETED" {
				wifiStatus = fmt.Sprintf(
					"WiFi Information\nSSID: %s\nIP Address: %s\n",
					status.SSID(),
					status.IPAddr(),
				)
			} else {
				wifiStatus = fmt.Sprintf("WiFi status is: %s\n", status.WPAState())
			}
		}
	}

	this.wifiStatus.SetText(wifiStatus)
}

func (this *networkPanel) setNetworkListItems(wpa wpasupplicant.Conn) {
	if this.overrideForDebugging {
		ssids := []string {
			"Vodafone-750C",
			"KabelBox-4E80",
			"Vodafone Hotspot",
			"Vodafone Homespot",
			"kabelBox-4E80",
			"Telekom_FON",
			"Test7",
			"Test8",
			"Test9-qwertyuiop-asdfghjkl-zxcvbnm-1234567890",
		}

		for i := 0; i < len(ssids); i++ {
			this.addNetwork(this.listBox, ssids[i])
		}
	} else {
		if wpa != nil {
			result, _ := wpa.ScanResults()
			for _, bss := range result {
				this.addNetwork(this.listBox, bss.SSID())
			}

			err := wpa.Scan()
			if err != nil {
				logger.LogError("NetworkPanel.update()", "Scan()", err)
			}
		} else {
			label := utils.MustLabel("\n\nWiFi management is not available\non this hardware")
			this.listBox.Add(label)
		}
	}
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
		this.UI.GoToPanel(GetConnectToNetworkPanelInstance(this.UI, ssid))
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
		this.UI.GoToPreviousPanel()
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
