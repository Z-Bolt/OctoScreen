package ui

import (
	"fmt"
	"sort"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

type FilesPanel struct {
	CommonPanel

	list *gtk.Box
}

func NewFilesPanel(ui *UI) *FilesPanel {
	m := &FilesPanel{CommonPanel: NewCommonPanel(ui)}
	m.initialize()
	return m
}

func (m *FilesPanel) initialize() {
	m.list = MustBox(gtk.ORIENTATION_VERTICAL, 0)
	m.list.SetVExpand(true)

	scroll, _ := gtk.ScrolledWindowNew(nil, nil)
	scroll.Add(m.list)

	box := MustBox(gtk.ORIENTATION_VERTICAL, 0)
	box.Add(scroll)
	box.Add(m.createActionBar())
	m.grid.Add(box)

	m.doLoadFiles()
}

func (m *FilesPanel) createActionBar() gtk.IWidget {
	bar := MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	bar.SetHAlign(gtk.ALIGN_END)
	bar.SetHExpand(true)
	bar.SetMarginTop(5)
	bar.SetMarginBottom(5)
	bar.SetMarginEnd(5)

	bar.Add(m.createRefreshButton())
	bar.Add(m.createInitReleaseSDButton())
	bar.Add(MustButton(MustImageFromFileWithSize("back.svg", 40, 40), m.UI.ShowDefaultPanel))

	return bar
}

func (m *FilesPanel) createRefreshButton() gtk.IWidget {
	return MustButton(MustImageFromFileWithSize("refresh.svg", 40, 40), m.doLoadFiles)
}

func (m *FilesPanel) doLoadFiles() {
	Logger.Info("Refreshing list of files")
	m.doRefreshSD()

	local := m.doLoadFilesFromLocation(octoprint.Local)
	sdcard := m.doLoadFilesFromLocation(octoprint.SDCard)

	s := byDate(local)
	s = append(s, sdcard...)
	sort.Sort(s)

	EmptyContainer(&m.list.Container)
	for _, f := range s {
		if f.IsFolder() {
			continue
		}

		m.addFile(m.list, f)
	}

	m.list.ShowAll()
}

func (m *FilesPanel) doRefreshSD() {
	if err := (&octoprint.SDRefreshRequest{}).Do(m.UI.Printer); err != nil {
		Logger.Error(err)
	}
}

func (m *FilesPanel) doLoadFilesFromLocation(l octoprint.Location) []*octoprint.FileInformation {
	r := &octoprint.FilesRequest{Location: l, Recursive: true}
	files, err := r.Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return []*octoprint.FileInformation{}
	}

	return files.Files
}

func (m *FilesPanel) addFile(b *gtk.Box, f *octoprint.FileInformation) {
	frame, _ := gtk.FrameNew("")

	name := MustLabel(f.Name)
	name.SetMarkup(fmt.Sprintf("<big>%s</big>", filenameEllipsis(f.Name)))
	name.SetHExpand(true)

	info := MustLabel("")
	info.SetMarkup(fmt.Sprintf("<small>Uploaded: <b>%s</b> - Size: <b>%s</b></small>",
		humanize.Time(f.Date.Time), humanize.Bytes(uint64(f.Size)),
	))

	labels := MustBox(gtk.ORIENTATION_VERTICAL, 5)
	labels.Add(name)
	labels.Add(info)

	actions := MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	actions.Add(m.createLoadAndPrintButton("load.svg", f, false))
	actions.Add(m.createLoadAndPrintButton("status.svg", f, true))

	file := MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	file.SetMarginTop(15)
	file.SetMarginEnd(15)
	file.SetMarginStart(15)
	file.SetMarginBottom(15)
	file.SetHExpand(true)

	file.Add(MustImageFromFileWithSize("file.svg", 35, 35))

	file.Add(labels)
	file.Add(actions)

	frame.Add(file)
	b.Add(frame)
}

func (m *FilesPanel) createLoadAndPrintButton(img string, f *octoprint.FileInformation, print bool) gtk.IWidget {
	return MustButton(
		MustImageFromFileWithSize(img, 30, 30),
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
}

func (m *FilesPanel) createInitReleaseSDButton() gtk.IWidget {
	release := MustImageFromFileWithSize("sd_eject.svg", 40, 40)
	init := MustImageFromFileWithSize("sd.svg", 40, 40)
	b := MustButton(release, nil)

	state := func() {
		time.Sleep(50 * time.Millisecond)
		switch m.isReady() {
		case true:
			b.SetImage(release)
		case false:
			b.SetImage(init)
		}
	}

	b.Connect("clicked", func() {
		var err error
		if !m.isReady() {
			err = (&octoprint.SDInitRequest{}).Do(m.UI.Printer)
		} else {
			err = (&octoprint.SDReleaseRequest{}).Do(m.UI.Printer)
		}

		if err != nil {
			Logger.Error(err)
		}

		state()
	})

	return b
}

func (m *FilesPanel) isReady() bool {
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
