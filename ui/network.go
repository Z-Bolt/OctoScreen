package ui

import (
	"fmt"
	"net"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"pifke.org/wpasupplicant"
)

var networkPanelInstance *networkPanel

type networkPanel struct {
	CommonPanel
	list                  *gtk.Box
	netStatus, wifiStatus *gtk.Label
}

func NetworkPanel(ui *UI, parent Panel) Panel {
	if networkPanelInstance == nil {
		m := &networkPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.initialize()
		m.b = NewBackgroundTask(time.Second*3, m.update)
		networkPanelInstance = m
	} else {
		networkPanelInstance.p = parent
	}

	return networkPanelInstance
}

func (m *networkPanel) update() {
	EmptyContainer(&m.list.Container)

	netStatus := ""
	addrs, _ := net.InterfaceAddrs()

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				netStatus += fmt.Sprintf("IP Address: %s\n", ipnet.IP.String())
			}
		}
	}

	m.netStatus.SetText(netStatus)

	wpa, _ := wpasupplicant.Unixgram("wlan0")

	if wpa != nil {
		s, _ := wpa.Status()
		wifiStatus := "No Wifi Connection"

		if s.WPAState() == "COMPLETED" {
			wifiStatus = fmt.Sprintf("Wifi Information\nSSID: %s\nIP Address: %s\n",
				s.SSID(), s.IPAddr())
		} else {
			wifiStatus = fmt.Sprintf("Wifi status is: %s\n", s.WPAState())
		}

		m.wifiStatus.SetText(wifiStatus)

		result, _ := wpa.ScanResults()
		for _, bss := range result {
			m.addNetwork(m.list, bss.SSID())
		}

		wpa.Scan()

	} else {
		m.list.Add(MustLabel("\n\nWifi management is not available\non this hardware"))
	}

	m.list.ShowAll()
}

func (m *networkPanel) initialize() {
	m.Grid().Attach(m.createNetworkList(), 0, 0, 4, 1)
	m.Grid().Attach(m.createLeftBar(), 4, 0, 3, 1)
}

func (m *networkPanel) createNetworkList() gtk.IWidget {

	m.list = MustBox(gtk.ORIENTATION_VERTICAL, 0)
	m.list.SetVExpand(true)

	scroll, _ := gtk.ScrolledWindowNew(nil, nil)
	scroll.SetProperty("overlay-scrolling", false)
	scroll.Add(m.list)

	return scroll
}

func (m *networkPanel) addNetwork(b *gtk.Box, n string) {
	frame, _ := gtk.FrameNew("")

	clicked := func() { m.UI.Add(ConnectionPanel(m.UI, m, n)) }

	button := MustButton(MustImageFromFileWithSize("network.svg", m.Scaled(25), m.Scaled(25)), clicked)

	name := MustButtonText(strEllipsisLen(n, 18), clicked)
	name.SetHExpand(true)

	network := MustBox(gtk.ORIENTATION_HORIZONTAL, 0)
	network.SetHExpand(true)

	network.Add(name)
	network.Add(button)

	frame.Add(network)
	b.Add(frame)
}

func (m *networkPanel) createLeftBar() gtk.IWidget {

	bar := MustBox(gtk.ORIENTATION_VERTICAL, 5)
	bar.SetHExpand(true)

	bar.Add(m.createInfoBar())
	bar.Add(m.createActionBar())

	return bar
}

func (m *networkPanel) createActionBar() gtk.IWidget {
	layout := MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	layout.SetHAlign(gtk.ALIGN_END)
	layout.SetHExpand(true)

	back := MustButton(MustImageFromFileWithSize("back.svg", m.Scaled(40), m.Scaled(40)), func() {
		m.UI.GoHistory()
	})

	layout.Add(back)

	return layout
}

func (m *networkPanel) createInfoBar() gtk.IWidget {
	info := MustBox(gtk.ORIENTATION_VERTICAL, 5)
	info.SetVExpand(true)
	info.SetHExpand(true)
	info.SetVAlign(gtk.ALIGN_START)
	info.SetHAlign(gtk.ALIGN_START)

	info.SetMarginEnd(25)
	info.SetMarginStart(25)

	t1 := MustLabel("<b>Network Information</b>")
	t1.SetHAlign(gtk.ALIGN_START)
	t1.SetMarginTop(25)

	info.Add(t1)

	m.netStatus = MustLabel("")
	m.netStatus.SetHAlign(gtk.ALIGN_START)
	info.Add(m.netStatus)

	m.wifiStatus = MustLabel("")
	m.wifiStatus.SetHAlign(gtk.ALIGN_START)
	m.wifiStatus.SetMarginTop(25)
	info.Add(m.wifiStatus)

	return info
}

var connectionPanelInstance *connectionPanel
var keyBoardChars = []byte{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '_', '+', '-',
	'=', '.', ',', '|', ':', ';', '/', '~', '`', '[', ']', '{', '}',
	'±', '§', '\\',
}

type connectionPanel struct {
	CommonPanel
	pass           *gtk.Entry
	cursorPosition int
	SSID           string
	SSIDLabel      *gtk.Label
}

func ConnectionPanel(ui *UI, parent Panel, SSID string) Panel {
	if connectionPanelInstance == nil {
		m := &connectionPanel{CommonPanel: NewCommonPanel(ui, parent), cursorPosition: 0}
		m.initialize()
		connectionPanelInstance = m
	}
	connectionPanelInstance.setSSID(SSID)

	return connectionPanelInstance
}

func (m *connectionPanel) initialize() {

	layout := MustBox(gtk.ORIENTATION_VERTICAL, 5)
	layout.SetHExpand(true)

	layout.Add(m.createTopBar())
	layout.Add(m.createKeyBoard())
	layout.Add(m.createActionBar())
	m.Grid().Add(layout)
}

func (m *connectionPanel) setSSID(SSID string) {
	m.SSID = SSID
	m.SSIDLabel.SetText(fmt.Sprintf("Enter password for \"%s\": ", strEllipsisLen(m.SSID, 18)))
}

func (m *connectionPanel) createTopBar() gtk.IWidget {
	m.pass, _ = gtk.EntryNew()
	m.pass.SetProperty("height-request", m.Scaled(40))
	m.pass.SetProperty("width-request", m.Scaled(150))
	m.pass.SetHExpand(true)

	top := MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	top.SetHExpand(true)
	top.SetMarginStart(25)
	top.SetHAlign(gtk.ALIGN_START)
	m.SSIDLabel = MustLabel(fmt.Sprintf("Pass for %s: ", m.SSID))
	top.Add(m.SSIDLabel)
	top.Add(m.pass)

	backspace := MustButton(MustImageFromFileWithSize("backspace.svg", m.Scaled(40), m.Scaled(40)), func() {
		if m.cursorPosition == 0 {
			return
		}

		m.pass.DeleteText(m.cursorPosition-1, m.cursorPosition)
		m.cursorPosition -= 1
	})

	top.Add(backspace)
	return top
}

func (m *connectionPanel) createActionBar() gtk.IWidget {
	back := MustButton(MustImageFromFileWithSize("back.svg", m.Scaled(40), m.Scaled(40)), func() {
		m.UI.GoHistory()
	})

	back.SetHAlign(gtk.ALIGN_END)

	layout := MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	layout.SetHAlign(gtk.ALIGN_END)
	layout.Add(m.createConnectToNetworkButtom())
	layout.Add(back)
	return layout
}

func (m *connectionPanel) createKeyBoard() gtk.IWidget {
	keys := MustGrid()
	keys.SetRowHomogeneous(true)
	keys.SetColumnHomogeneous(true)

	keyBoard, _ := gtk.ScrolledWindowNew(nil, nil)
	keyBoard.SetProperty("overlay-scrolling", false)
	keyBoard.SetVExpand(true)
	keyBoard.Add(keys)

	row := m.Scaled(3)

	for i, k := range keyBoardChars {
		buttonHander := &keyButtonHander{char: k, p: m}
		button := MustButtonText(string(k), buttonHander.clicked)
		ctx, _ := button.GetStyleContext()
		ctx.AddClass("keyboard")
		button.SetProperty("height-request", m.Scaled(40))
		keys.Attach(button, i%row, i/row, 1, 1)
	}

	return keyBoard
}

func (m *connectionPanel) createConnectToNetworkButtom() gtk.IWidget {
	var b *gtk.Button

	b = MustButtonText("Connect", func() {
		b.SetSensitive(false)
		time.Sleep(time.Second * 1)
		psk, _ := m.pass.GetText()
		wpa, _ := wpasupplicant.Unixgram("wlan0")

		if wpa != nil {
			wpa.RemoveAllNetworks()
			wpa.AddNetwork()
			wpa.SetNetwork(0, "ssid", m.SSID)
			wpa.SetNetwork(0, "psk", psk)

			go wpa.EnableNetwork(0)
			time.Sleep(time.Second * 1)
			go wpa.SaveConfig()
		}

		time.Sleep(time.Second * 1)
		m.UI.GoHistory()
	})

	ctx, _ := b.GetStyleContext()
	ctx.AddClass("color3")

	b.SetProperty("width-request", m.Scaled(150))

	return b
}

type keyButtonHander struct {
	char byte
	p    *connectionPanel
}

func (m *keyButtonHander) clicked() {
	m.p.pass.InsertText(string(m.char), m.p.cursorPosition)
	m.p.cursorPosition += 1
}
