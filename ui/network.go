package ui

import (
	"fmt"
	"net"

	"github.com/gotk3/gotk3/gtk"
	"pifke.org/wpasupplicant"
)

var networkPanelInstance *networkPanel

type networkPanel struct {
	CommonPanel
	wpa wpasupplicant.Conn
}

func NetworkPanel(ui *UI, parent Panel) Panel {
	if networkPanelInstance == nil {
		m := &networkPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.initialize()
		networkPanelInstance = m
	}

	return networkPanelInstance
}

func (m *networkPanel) initialize() {
	m.wpa, _ = wpasupplicant.Unixgram("wlan0")

	m.Grid().Attach(m.createNetworkList(), 0, 0, 3, 1)
	m.Grid().Attach(m.createLeftBar(), 3, 0, 2, 1)
}

func (m *networkPanel) createNetworkList() gtk.IWidget {

	if m.wpa == nil {
		return MustLabel("Wifi management was disabled")
	}

	list := MustBox(gtk.ORIENTATION_VERTICAL, 0)
	list.SetVExpand(true)

	scroll, _ := gtk.ScrolledWindowNew(nil, nil)
	scroll.SetProperty("overlay-scrolling", false)
	scroll.Add(list)

	result, _ := m.wpa.ScanResults()

	fmt.Println("Getting list of networks")

	for _, bss := range result {
		m.addNetwork(list, bss.SSID())
	}

	return scroll
}

func (m *networkPanel) addNetwork(b *gtk.Box, n string) {
	frame, _ := gtk.FrameNew("")

	name := MustLabel(n)
	name.SetMarkup(fmt.Sprintf("<big>%s</big>", strEllipsisLen(n, 18)))
	name.SetHExpand(true)

	actions := MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	actions.Add(m.createConnectButton(n))

	network := MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	network.SetMarginTop(5)
	network.SetMarginEnd(15)
	network.SetMarginStart(15)
	network.SetMarginBottom(5)
	network.SetHExpand(true)

	network.Add(MustImageFromFileWithSize("network.svg", m.Scaled(35), m.Scaled(35)))

	network.Add(name)
	network.Add(actions)

	frame.Add(network)
	b.Add(frame)
}

func (m *networkPanel) createConnectButton(n string) gtk.IWidget {
	return MustButton(MustImageFromFileWithSize("open.svg", m.Scaled(40), m.Scaled(40)), func() {
		m.UI.Add(ConnectionPanel(m.UI, m, n))
	})
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
	addrs, _ := net.InterfaceAddrs()

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				info.Add(MustLabel("IP Address: <b>%s</b>", ipnet.IP.String()))
			}
		}
	}

	if m.wpa != nil {
		s, _ := m.wpa.Status()
		text := "<b>No Wifi Connection</b>"

		if s.WPAState() == "COMPLETED" {
			text = fmt.Sprintf("<b>Wifi Information</b>\nSSID: <b>%s</b>\nIP Address: <b>%s</b>",
				s.SSID(), s.IPAddr())
		}

		t2 := MustLabel(text)
		t2.SetHAlign(gtk.ALIGN_START)
		t2.SetMarginTop(25)
		info.Add(t2)
	}

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

	delButton := MustButton(MustImageFromFileWithSize("backspace.svg", m.Scaled(40), m.Scaled(40)), func() {
		if m.cursorPosition == 0 {
			return
		}

		m.pass.DeleteText(m.cursorPosition-1, m.cursorPosition)
		m.cursorPosition -= 1
	})

	top.Add(delButton)

	layout := MustBox(gtk.ORIENTATION_VERTICAL, 5)
	layout.SetHExpand(true)

	layout.Add(top)
	layout.Add(m.createKeyBoard())
	layout.Add(m.createActionBar())

	m.Grid().Add(layout)
}

func (m *connectionPanel) setSSID(SSID string) {
	m.SSID = SSID
	m.SSIDLabel.SetText(fmt.Sprintf("Enter password for \"%s\": ", strEllipsisLen(m.SSID, 18)))
}

func (m *connectionPanel) createActionBar() gtk.IWidget {
	back := MustButton(MustImageFromFileWithSize("back.svg", m.Scaled(40), m.Scaled(40)), func() {
		m.UI.GoHistory()
	})

	back.SetHAlign(gtk.ALIGN_END)

	connect := MustButtonText("Connect", func() {})
	connect.SetProperty("width-request", m.Scaled(150))

	layout := MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	layout.SetHAlign(gtk.ALIGN_END)
	layout.Add(connect)
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

		button.SetProperty("height-request", m.Scaled(40))
		keys.Attach(button, i%row, i/row, 1, 1)
	}

	return keyBoard
}

type keyButtonHander struct {
	char byte
	p    *connectionPanel
}

func (m *keyButtonHander) clicked() {
	m.p.pass.InsertText(string(m.char), m.p.cursorPosition)
	m.p.cursorPosition += 1
}
