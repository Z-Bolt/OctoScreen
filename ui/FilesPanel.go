package ui

import (
	"fmt"
	"sort"

	"github.com/dustin/go-humanize"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/interfaces"
	// "github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)

var filesPanelInstance *filesPanel

type filesPanel struct {
	CommonPanel

	listBox				*gtk.Box
	locationHistory		locationHistory
}

func FilesPanel(
	ui					*UI,
	parentPanel			interfaces.IPanel,
) *filesPanel {
	if filesPanelInstance == nil {
		locationHistory := locationHistory {
			locations: []octoprint.Location{octoprint.Local},
		}

		instance := &filesPanel {
			CommonPanel: NewCommonPanel(ui, parentPanel),
			locationHistory: locationHistory,
		}
		instance.initialize()
		filesPanelInstance = instance
	}

	return filesPanelInstance
}

func (this *filesPanel) initialize() {
	this.listBox = utils.MustBox(gtk.ORIENTATION_VERTICAL, 0)
	this.listBox.SetVExpand(true)

	scroll, _ := gtk.ScrolledWindowNew(nil, nil)
	scroll.SetProperty("overlay-scrolling", false)
	scroll.Add(this.listBox)

	box := utils.MustBox(gtk.ORIENTATION_VERTICAL, 0)
	box.Add(scroll)
	box.Add(this.createActionBar())
	this.Grid().Add(box)

	this.doLoadFiles()
}

func (this *filesPanel) createActionBar() gtk.IWidget {
	actionBar := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	actionBar.SetHAlign(gtk.ALIGN_END)
	actionBar.SetHExpand(true)
	actionBar.SetMarginTop(5)
	actionBar.SetMarginBottom(5)
	actionBar.SetMarginEnd(5)

	actionBar.Add(this.createRefreshButton())
	actionBar.Add(this.createBackButton())

	return actionBar
}

func (this *filesPanel) createRefreshButton() gtk.IWidget {
	image := utils.MustImageFromFileWithSize("refresh.svg", this.Scaled(40), this.Scaled(40))
	return utils.MustButton(image, this.doLoadFiles)
}

func (this *filesPanel) createBackButton() gtk.IWidget {
	image := utils.MustImageFromFileWithSize("back.svg", this.Scaled(40), this.Scaled(40))
	return utils.MustButton(image, func() {
		if this.locationHistory.isRoot() {
			this.UI.GoHistory()
		} else {
			this.locationHistory.goBack()
			this.doLoadFiles()
		}
	})
}

func (this *filesPanel) doLoadFiles() {
	var files []*octoprint.FileInformation

	utils.Logger.Info("Loading list of files from: ", string(this.locationHistory.current()))

	filesRequest := &octoprint.FilesRequest{Location: this.locationHistory.current(), Recursive: false}
	filesResponse, err := filesRequest.Do(this.UI.Printer)
	if err != nil {
		utils.LogError("files.doLoadFiles()", "Do(FilesRequest)", err)
		files = []*octoprint.FileInformation{}
	} else {
		files = filesResponse.Files
	}

	s := byDate(files)
	sort.Sort(s)

	utils.EmptyTheContainer(&this.listBox.Container)

	for _, fileInformation := range s {
		if fileInformation.IsFolder() {
			this.addFolder(this.listBox, fileInformation)
		}
	}

	for _, fileInformation := range s {
		if !fileInformation.IsFolder() {
			this.addFile(this.listBox, fileInformation)
		}
	}

	this.listBox.ShowAll()
}

func (this *filesPanel) addFile(box *gtk.Box, fileInformation *octoprint.FileInformation) {
	frame, _ := gtk.FrameNew("")

	name := utils.MustLabel(fileInformation.Name)
	name.SetMarkup(fmt.Sprintf("<big>%s</big>", utils.StrEllipsis(fileInformation.Name)))
	name.SetHExpand(true)
	name.SetHAlign(gtk.ALIGN_START)

	info := utils.MustLabel("")
	info.SetHAlign(gtk.ALIGN_START)
	info.SetMarkup(fmt.Sprintf("<small>Uploaded: <b>%s</b> - Size: <b>%s</b></small>",
		humanize.Time(fileInformation.Date.Time), humanize.Bytes(uint64(fileInformation.Size)),
	))

	labels := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	labels.Add(name)
	labels.Add(info)
	labels.SetVExpand(true)
	labels.SetVAlign(gtk.ALIGN_CENTER)
	labels.SetHAlign(gtk.ALIGN_START)

	actions := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	actions.Add(this.createLoadAndPrintButton("print.svg", fileInformation, true))

	file := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	file.SetMarginTop(1)
	file.SetMarginEnd(15)
	file.SetMarginStart(15)
	file.SetMarginBottom(1)
	file.SetHExpand(true)

	image := utils.MustImageFromFileWithSize("file-stl.svg", this.Scaled(35), this.Scaled(35))
	file.Add(image)

	file.Add(labels)
	file.Add(actions)

	frame.Add(file)
	box.Add(frame)
}

func (this *filesPanel) addFolder(box *gtk.Box, fileInformation *octoprint.FileInformation) {
	frame, _ := gtk.FrameNew("")

	labels := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)

	nameLabel := utils.MustLabel(fileInformation.Name)
	nameLabel.SetMarkup(fmt.Sprintf("<big>%s</big>", utils.StrEllipsis(fileInformation.Name)))
	nameLabel.SetHExpand(true)
	nameLabel.SetHAlign(gtk.ALIGN_START)
	labels.Add(nameLabel)

	info := utils.MustLabel("")
	info.SetHAlign(gtk.ALIGN_START)
	info.SetMarkup(fmt.Sprintf("<small>Size: <b>%s</b></small>",
		humanize.Bytes(uint64(fileInformation.Size)),
	))
	labels.Add(info)

	labels.SetVExpand(true)
	labels.SetVAlign(gtk.ALIGN_CENTER)
	labels.SetHAlign(gtk.ALIGN_START)

	actionsBox := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	actionsBox.Add(this.createOpenFolderButton(fileInformation))

	fileBox := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	fileBox.SetMarginTop(1)
	fileBox.SetMarginEnd(15)
	fileBox.SetMarginStart(15)
	fileBox.SetMarginBottom(1)
	fileBox.SetHExpand(true)

	image := utils.MustImageFromFileWithSize("folder.svg", this.Scaled(35), this.Scaled(35))
	fileBox.Add(image)

	fileBox.Add(labels)
	fileBox.Add(actionsBox)

	frame.Add(fileBox)
	box.Add(frame)
}

func (this *filesPanel) createLoadAndPrintButton(imageFileName string, fileInformation *octoprint.FileInformation, shouldPrint bool) gtk.IWidget {
	button := utils.MustButton(
		utils.MustImageFromFileWithSize(imageFileName, this.Scaled(40), this.Scaled(40)),
		utils.MustConfirmDialogBox(this.UI.window, "Are you sure you want to proceed?", func() {
			selectFileRequest := &octoprint.SelectFileRequest{}
			selectFileRequest.Location = octoprint.Local
			selectFileRequest.Path = fileInformation.Path
			selectFileRequest.Print = shouldPrint

			utils.Logger.Infof("Loading file %q, printing: %v", fileInformation.Name, shouldPrint)
			if err := selectFileRequest.Do(this.UI.Printer); err != nil {
				utils.LogError("files.createLoadAndPrintButton()", "Do(SelectFileRequest)", err)
				return
			}
		}),
	)

	ctx, _ := button.GetStyleContext()
	ctx.AddClass("color-warning-sign-yellow")
	ctx.AddClass("file-list")

	return button
}

func (this *filesPanel) createOpenFolderButton(fileInformation *octoprint.FileInformation) gtk.IWidget {
	image := utils.MustImageFromFileWithSize("open.svg", this.Scaled(40), this.Scaled(40))
	button := utils.MustButton(image, func() {
		this.locationHistory.goForward(fileInformation.Path)
		this.doLoadFiles()
	})

	ctx, _ := button.GetStyleContext()
	ctx.AddClass("color1")
	ctx.AddClass("file-list")

	return button
}

/*
func (this *filesPanel) isReady() bool {
	state, err := (&octoprint.SDStateRequest{}).Do(this.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return false
	}

	return state.Ready
}
*/
