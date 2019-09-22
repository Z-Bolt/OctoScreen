package ui

import (
	"fmt"
	"sort"

	"github.com/dustin/go-humanize"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

var filesPanelInstance *filesPanel

type filesPanel struct {
	CommonPanel

	list     *gtk.Box
	location locationHistory
}

func FilesPanel(ui *UI, parent Panel) Panel {
	if filesPanelInstance == nil {
		l := locationHistory{locations: []octoprint.Location{octoprint.Local}}
		m := &filesPanel{CommonPanel: NewCommonPanel(ui, parent), location: l}
		m.initialize()
		filesPanelInstance = m
	}

	return filesPanelInstance
}

func (m *filesPanel) initialize() {
	m.list = MustBox(gtk.ORIENTATION_VERTICAL, 0)
	m.list.SetVExpand(true)

	scroll, _ := gtk.ScrolledWindowNew(nil, nil)
	scroll.SetProperty("overlay-scrolling", false)
	scroll.Add(m.list)

	box := MustBox(gtk.ORIENTATION_VERTICAL, 0)
	box.Add(scroll)
	box.Add(m.createActionBar())
	m.Grid().Add(box)

	m.doLoadFiles()
}

func (m *filesPanel) createActionBar() gtk.IWidget {
	bar := MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	bar.SetHAlign(gtk.ALIGN_END)
	bar.SetHExpand(true)
	bar.SetMarginTop(5)
	bar.SetMarginBottom(5)
	bar.SetMarginEnd(5)

	bar.Add(m.createRefreshButton())
	bar.Add(m.createBackButton())

	return bar
}

func (m *filesPanel) createRefreshButton() gtk.IWidget {
	return MustButton(MustImageFromFileWithSize("refresh.svg", m.Scaled(40), m.Scaled(40)), m.doLoadFiles)
}

func (m *filesPanel) createBackButton() gtk.IWidget {
	return MustButton(MustImageFromFileWithSize("back.svg", m.Scaled(40), m.Scaled(40)), func() {
		if m.location.isRoot() {
			m.UI.GoHistory()
		} else {
			m.location.goBack()
			m.doLoadFiles()
		}
	})
}

func (m *filesPanel) doLoadFiles() {
	var files []*octoprint.FileInformation

	Logger.Info("Loading list of files from: ", string(m.location.current()))

	r := &octoprint.FilesRequest{Location: m.location.current(), Recursive: false}
	folder, err := r.Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		files = []*octoprint.FileInformation{}
	} else {
		files = folder.Files
	}

	s := byDate(files)
	sort.Sort(s)

	EmptyContainer(&m.list.Container)

	for _, f := range s {
		if f.IsFolder() {
			m.addFolder(m.list, f)
		}
	}

	for _, f := range s {
		if !f.IsFolder() {
			m.addFile(m.list, f)
		}
	}

	m.list.ShowAll()
}

func (m *filesPanel) addFile(b *gtk.Box, f *octoprint.FileInformation) {
	frame, _ := gtk.FrameNew("")

	name := MustLabel(f.Name)
	name.SetMarkup(fmt.Sprintf("<big>%s</big>", strEllipsis(f.Name)))
	name.SetHExpand(true)
	name.SetHAlign(gtk.ALIGN_START)

	info := MustLabel("")
	info.SetHAlign(gtk.ALIGN_START)
	info.SetMarkup(fmt.Sprintf("<small>Uploaded: <b>%s</b> - Size: <b>%s</b></small>",
		humanize.Time(f.Date.Time), humanize.Bytes(uint64(f.Size)),
	))

	labels := MustBox(gtk.ORIENTATION_VERTICAL, 5)
	labels.Add(name)
	labels.Add(info)
	labels.SetVExpand(true)
	labels.SetVAlign(gtk.ALIGN_CENTER)
	labels.SetHAlign(gtk.ALIGN_START)

	actions := MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	actions.Add(m.createLoadAndPrintButton("print.svg", f, true))

	file := MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	file.SetMarginTop(1)
	file.SetMarginEnd(15)
	file.SetMarginStart(15)
	file.SetMarginBottom(1)
	file.SetHExpand(true)

	file.Add(MustImageFromFileWithSize("file.svg", m.Scaled(35), m.Scaled(35)))

	file.Add(labels)
	file.Add(actions)

	frame.Add(file)
	b.Add(frame)
}

func (m *filesPanel) addFolder(b *gtk.Box, f *octoprint.FileInformation) {
	frame, _ := gtk.FrameNew("")

	name := MustLabel(f.Name)
	name.SetMarkup(fmt.Sprintf("<big>%s</big>", strEllipsis(f.Name)))
	name.SetHExpand(true)
	name.SetHAlign(gtk.ALIGN_START)

	info := MustLabel("")
	info.SetHAlign(gtk.ALIGN_START)
	info.SetMarkup(fmt.Sprintf("<small>Size: <b>%s</b></small>",
		humanize.Bytes(uint64(f.Size)),
	))

	labels := MustBox(gtk.ORIENTATION_VERTICAL, 5)
	labels.Add(name)
	labels.Add(info)
	labels.SetVExpand(true)
	labels.SetVAlign(gtk.ALIGN_CENTER)
	labels.SetHAlign(gtk.ALIGN_START)

	actions := MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	actions.Add(m.createOpenFolderButton(f))

	file := MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	file.SetMarginTop(1)
	file.SetMarginEnd(15)
	file.SetMarginStart(15)
	file.SetMarginBottom(1)
	file.SetHExpand(true)

	file.Add(MustImageFromFileWithSize("folder.svg", m.Scaled(35), m.Scaled(35)))

	file.Add(labels)
	file.Add(actions)

	frame.Add(file)
	b.Add(frame)
}

func (m *filesPanel) createLoadAndPrintButton(img string, f *octoprint.FileInformation, print bool) gtk.IWidget {
	b := MustButton(
		MustImageFromFileWithSize(img, m.Scaled(40), m.Scaled(40)),
		MustConfirmDialog(m.UI.w, "Are you sure you want to proceed?", func() {
			r := &octoprint.SelectFileRequest{}
			r.Location = octoprint.Local
			r.Path = f.Path
			r.Print = print

			Logger.Infof("Loading file %q, printing: %v", f.Name, print)
			if err := r.Do(m.UI.Printer); err != nil {
				Logger.Error(err)
				return
			}
		}),
	)
	ctx, _ := b.GetStyleContext()
	ctx.AddClass("color3")
	ctx.AddClass("file-list")
	return b
}

func (m *filesPanel) createOpenFolderButton(f *octoprint.FileInformation) gtk.IWidget {
	b := MustButton(MustImageFromFileWithSize("open.svg", m.Scaled(40), m.Scaled(40)), func() {
		m.location.goForward(f.Path)
		m.doLoadFiles()
	})

	ctx, _ := b.GetStyleContext()
	ctx.AddClass("color1")
	ctx.AddClass("file-list")

	return b
}

func (m *filesPanel) isReady() bool {
	state, err := (&octoprint.SDStateRequest{}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return false
	}

	return state.Ready
}

type byDate []*octoprint.FileInformation

func (s byDate) Len() int           { return len(s) }
func (s byDate) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byDate) Less(i, j int) bool { return s[j].Date.Time.Before(s[i].Date.Time) }

type locationHistory struct {
	locations []octoprint.Location
}

func (l *locationHistory) current() octoprint.Location {
	return l.locations[len(l.locations)-1]
}

func (l *locationHistory) goForward(folder string) {
	newLocation := string(l.current()) + "/" + folder
	l.locations = append(l.locations, octoprint.Location(newLocation))
}

func (l *locationHistory) goBack() {
	l.locations = l.locations[0 : len(l.locations)-1]
}

func (l *locationHistory) isRoot() bool {
	if len(l.locations) > 1 {
		return false
	} else {
		return true
	}
}
