package ui

import (
	"fmt"
	"sort"
	//"time"

	"github.com/dustin/go-humanize"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

var filesPanelInstance *filesPanel

type filesPanel struct {
	CommonPanel

	list *gtk.Box
}

func FilesPanel(ui *UI, parent Panel) Panel {
	// if filesPanelInstance == nil {
		m := &filesPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.initialize()
		filesPanelInstance = m
	// }

	return filesPanelInstance
}

func (m *filesPanel) initialize() {
	m.list = MustBox(gtk.ORIENTATION_VERTICAL, 0)
	m.list.SetVExpand(true)

	scroll, _ := gtk.ScrolledWindowNew(nil, nil)
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
	//bar.Add(m.createInitReleaseSDButton())
	bar.Add(MustButton(MustImageFromFileWithSize("back.svg", 35, 35), m.UI.GoHistory))

	return bar
}

func (m *filesPanel) createRefreshButton() gtk.IWidget {
	return MustButton(MustImageFromFileWithSize("refresh.svg", 35, 35), m.doLoadFiles)
}

func (m *filesPanel) doLoadFiles() {
	Logger.Info("Refreshing list of files")
	//m.doRefreshSD()

	local := m.doLoadFilesFromLocation(octoprint.Local)
	//sdcard := m.doLoadFilesFromLocation(octoprint.SDCard)

	s := byDate(local)
	//s = append(s, sdcard...)
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

// func (m *filesPanel) doRefreshSD() {
// 	if err := (&octoprint.SDRefreshRequest{}).Do(m.UI.Printer); err != nil {
// 		Logger.Error(err)
// 	}
// }

func (m *filesPanel) doLoadFilesFromLocation(l octoprint.Location) []*octoprint.FileInformation {
	r := &octoprint.FilesRequest{Location: l, Recursive: true}
	files, err := r.Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return []*octoprint.FileInformation{}
	}

	return files.Files
}

func (m *filesPanel) addFile(b *gtk.Box, f *octoprint.FileInformation) {
	frame, _ := gtk.FrameNew("")

	name := MustLabel(f.Name)
	name.SetMarkup(fmt.Sprintf("<small>%s</small>", filenameEllipsis(f.Name)))
	name.SetHExpand(true)
	name.SetHAlign(gtk.ALIGN_START)
	name.SetMarginTop(5)

	info := MustLabel("")
	info.SetMarkup(fmt.Sprintf("<small>%s - %s</small>",
		humanize.Time(f.Date.Time), humanize.Bytes(uint64(f.Size)),
	))
	info.SetHAlign(gtk.ALIGN_START)

	labels := MustBox(gtk.ORIENTATION_VERTICAL, 5)
	labels.Add(name)
	labels.Add(info)

	actions := MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	actions.Add(m.createLoadAndPrintButton("status.svg", f))
	// actions.Add(m.createLoadAndPrintButton("load.svg", f, false))
	actions.Add(m.createDeleteButton("delete.svg", f))
	actions.SetHAlign(gtk.ALIGN_END)

	file := MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	file.SetMarginTop(5)
	file.SetMarginEnd(5)
	file.SetMarginStart(5)
	file.SetMarginBottom(5)
	file.SetHExpand(true)

	file.Add(labels)
	file.Add(actions)

	frame.Add(file)
	b.Add(frame)
}

func (m *filesPanel) createLoadAndPrintButton(img string, f *octoprint.FileInformation) gtk.IWidget {
	return MustButton(
		MustImageFromFileWithSize(img, 20, 20),
		PrintDialog(m.UI.w, "File loaded. Start printing?\n"+filenameEllipsis_short(f.Name), f.Path, m),
	)
}



func PrintDialog(parent *gtk.Window, msg string, pfad string, ma *filesPanel) func() {
	return func() {
		win := gtk.MessageDialogNewWithMarkup(
			parent,
			gtk.DIALOG_MODAL,
			gtk.MESSAGE_INFO,
			gtk.BUTTONS_YES_NO,
			"",
		)

		win.SetMarkup(CleanHTML(msg))
		defer win.Destroy()

		box, _ := win.GetContentArea()
		box.SetMarginStart(15)
		box.SetMarginEnd(15)
		box.SetMarginTop(15)
		box.SetMarginBottom(15)

		ctx, _ := win.GetStyleContext()
		ctx.AddClass("dialog")

		ergebnis := win.Run()

		if ergebnis == int(gtk.RESPONSE_NO) {
			ro := &octoprint.SelectFileRequest{}
			ro.Location = octoprint.Local
			ro.Path = pfad
			ro.Print = false

			Logger.Infof("Loading file %q", ro)
			if err := ro.Do(ma.UI.Printer); err != nil {
				Logger.Error(err)
				return
			}
		}else if ergebnis == int(gtk.RESPONSE_YES) {
			rt := &octoprint.SelectFileRequest{}
			rt.Location = octoprint.Local
			rt.Path = pfad
			rt.Print = true

			Logger.Infof("Printing file %q", rt)
			if err := rt.Do(ma.UI.Printer); err != nil {
				Logger.Error(err)
				return
			}
		}
	}
}

func (m *filesPanel) createDeleteButton(img string, de *octoprint.FileInformation) gtk.IWidget {
	return MustButton(
		MustImageFromFileWithSize(img, 20, 20),
		MustConfirmDialog(m.UI.w, "Delete file?\n"+filenameEllipsis_short(de.Name), func() {
			del := &octoprint.DeleteFileRequest{}
			del.Location = octoprint.Local
			del.Path = de.Path

			Logger.Infof("RM %q FROM %v", de.Path, octoprint.Local)
			if err := del.Do(m.UI.Printer); err != nil {
				Logger.Error(err)
				return
			}
			
			m.doLoadFiles()
		}),
	)
}

// func (m *filesPanel) createInitReleaseSDButton() gtk.IWidget {
// 	release := MustImageFromFileWithSize("sd_eject.svg", 40, 40)
// 	init := MustImageFromFileWithSize("sd.svg", 40, 40)
// 	b := MustButton(release, nil)

// 	state := func() {
// 		time.Sleep(50 * time.Millisecond)
// 		switch m.isReady() {
// 		case true:
// 			b.SetImage(release)
// 		case false:
// 			b.SetImage(init)
// 		}
// 	}

// 	b.Connect("clicked", func() {
// 		var err error
// 		if !m.isReady() {
// 			err = (&octoprint.SDInitRequest{}).Do(m.UI.Printer)
// 		} else {
// 			err = (&octoprint.SDReleaseRequest{}).Do(m.UI.Printer)
// 		}

// 		if err != nil {
// 			Logger.Error(err)
// 		}

// 		state()
// 	})

// 	return b
// }

// func (m *filesPanel) isReady() bool {
// 	state, err := (&octoprint.SDStateRequest{}).Do(m.UI.Printer)
// 	if err != nil {
// 		Logger.Error(err)
// 		return false
// 	}

// 	return state.Ready
// }

type byDate []*octoprint.FileInformation

func (s byDate) Len() int           { return len(s) }
func (s byDate) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byDate) Less(i, j int) bool { return s[j].Date.Time.Before(s[i].Date.Time) }
