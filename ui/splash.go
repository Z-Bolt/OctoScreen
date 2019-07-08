package ui

import "github.com/gotk3/gotk3/gtk"

type SplashPanel struct {
	CommonPanel
	Label *gtk.Label
}

func NewSplashPanel(ui *UI) *SplashPanel {
	m := &SplashPanel{CommonPanel: NewCommonPanel(ui, nil)}
	m.initialize()
	return m
}

func (m *SplashPanel) initialize() {
	logo := MustImageFromFile("logo.png")
	m.Label = MustLabel("Initializing printer...")
	m.Label.SetLineWrap(true)

	box := MustBox(gtk.ORIENTATION_VERTICAL, 15)
	box.SetVAlign(gtk.ALIGN_CENTER)
	box.SetVExpand(true)
	box.SetHExpand(true)

	box.Add(logo)
	box.Add(m.Label)

	m.Grid().Attach(box, 1, 0, 3, 2)
	m.Grid().Attach(MustButtonImageStyle("System", "info.svg", "color4", m.showSystem), 4, 0, 1, 1)
	m.Grid().Attach(MustButtonImageStyle("Network", "network.svg", "color2", m.showNetwork), 4, 1, 1, 1)
}

func (m *SplashPanel) showNetwork() {
	m.UI.Add(NetworkPanel(m.UI, m))
}

func (m *SplashPanel) showSystem() {
	m.UI.Add(SystemPanel(m.UI, m))
}
