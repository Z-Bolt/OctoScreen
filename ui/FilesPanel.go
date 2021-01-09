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
	locationHistory		utils.LocationHistory
}

func FilesPanel(
	ui					*UI,
	parentPanel			interfaces.IPanel,
) *filesPanel {
	if filesPanelInstance == nil {
		locationHistory := utils.LocationHistory {
			Locations: []octoprint.Location{},
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
		if this.locationHistory.Length() < 1 {
			this.UI.GoToPreviousPanel()
		} else if this.locationHistory.IsRoot() {
			this.locationHistory.GoBack()
			this.doLoadFiles()
		} else {
			this.locationHistory.GoBack()
			this.doLoadFiles()
		}
	})
}

func (this *filesPanel) doLoadFiles() {
	utils.EmptyTheContainer(&this.listBox.Container)

	sortedFiles := this.getSortedFiles()
	if sortedFiles == nil {
		this.addRootLocations()
	} else {
		this.addSortedFiles(sortedFiles)
	}

	this.listBox.ShowAll()
}

func (this *filesPanel) getSortedFiles() []*octoprint.FileInformation {
	var files []*octoprint.FileInformation

	length := this.locationHistory.Length()
	if length < 1 {
		return nil
	}

	current := this.locationHistory.CurrentLocation()
	utils.Logger.Info("Loading list of files from: ", string(current))

	filesRequest := &octoprint.FilesRequest {
		Location: current,
		Recursive: false,
	}
	filesResponse, err := filesRequest.Do(this.UI.Client)
	if err != nil {
		utils.LogError("files.getSortedFiles()", "Do(FilesRequest)", err)
		files = []*octoprint.FileInformation{}
	} else {
		files = filesResponse.Files
	}

	sortedFiles := utils.FileInformationsByDate(files)
	sort.Sort(sortedFiles)

	return sortedFiles
}


func (this *filesPanel) addSortedFiles(sortedFiles []*octoprint.FileInformation) {
	var index int = 0

	for _, fileInformation := range sortedFiles {
		if fileInformation.IsFolder() {
			this.addItem(fileInformation, index)
			index++
		}
	}

	for _, fileInformation := range sortedFiles {
		if !fileInformation.IsFolder() {
			this.addItem(fileInformation, index)
			index++
		}
	}
}


func (this *filesPanel) addItem(fileInformation *octoprint.FileInformation, index int) {
	isFolder := fileInformation.IsFolder()

	topBox := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)

	var itemImage *gtk.Image
	if isFolder {
		itemImage = utils.MustImageFromFileWithSize("folder.svg", this.Scaled(35), this.Scaled(35))
	} else {
		itemImage = utils.MustImageFromFileWithSize("file-gcode.svg", this.Scaled(35), this.Scaled(35))
	}
	topBox.Add(itemImage)




	name := fileInformation.Name
	nameLabel := utils.MustLabel(name)
	nameLabel.SetMarkup(fmt.Sprintf("<big>%s</big>", utils.StrEllipsis(name)))
	nameLabel.SetHExpand(true)
	nameLabel.SetHAlign(gtk.ALIGN_START)

	infoLabel := utils.MustLabel("")
	infoLabel.SetHAlign(gtk.ALIGN_START)
	if isFolder {
		infoLabel.SetMarkup(fmt.Sprintf("<small>Size: <b>%s</b></small>",
			humanize.Bytes(uint64(fileInformation.Size)),
		))
	} else {
		infoLabel.SetMarkup(fmt.Sprintf("<small>Uploaded: <b>%s</b> - Size: <b>%s</b></small>",
			humanize.Time(fileInformation.Date.Time), humanize.Bytes(uint64(fileInformation.Size)),
		))
	}

	labelsBox := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	labelsBox.Add(nameLabel)
	labelsBox.Add(infoLabel)
	labelsBox.SetVExpand(false)
	labelsBox.SetVAlign(gtk.ALIGN_CENTER)
	labelsBox.SetHAlign(gtk.ALIGN_START)
	labelsBoxStyleContext, _ := labelsBox.GetStyleContext()
	labelsBoxStyleContext.AddClass("labels-box")

	topBox.Add(labelsBox)




	var itemButton *gtk.Button
	if isFolder {
		itemButton = this.createOpenFolderButton(fileInformation)
	} else {
		itemButton = this.createLoadAndPrintButton("print.svg", fileInformation)
	}

	actionBox := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	actionBox.Add(itemButton)
	topBox.Add(actionBox)




	listItemBox := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	listItemBox.SetMarginTop(1)
	listItemBox.SetMarginBottom(1)
	listItemBox.SetMarginStart(15)
	listItemBox.SetMarginEnd(15)
	listItemBox.SetHExpand(true)
	listItemBox.Add(topBox)


	listItemFrame, _ := gtk.FrameNew("")
	listItemFrame.Add(listItemBox)



	itemButtonStyleContext, _ := itemButton.GetStyleContext()
	listItemBoxStyleContext, _:= listItemBox.GetStyleContext()
	listItemFrameStyleContext, _ := listItemFrame.GetStyleContext()
	if index % 2 != 0 {
		itemButtonStyleContext.AddClass("list-item-nth-child-even")
		listItemBoxStyleContext.AddClass("list-item-nth-child-even")
		listItemFrameStyleContext.AddClass("list-item-nth-child-even")
	}

	this.listBox.Add(listItemFrame)
}


func (this *filesPanel) createOpenLocationButton(location octoprint.Location) *gtk.Button {
	image := utils.MustImageFromFileWithSize("open.svg", this.Scaled(40), this.Scaled(40))
	button := utils.MustButton(image, func() {
		this.locationHistory = utils.LocationHistory {
			Locations: []octoprint.Location{location},
		}

		this.doLoadFiles()
	})

	ctx, _ := button.GetStyleContext()
	ctx.AddClass("color1")
	ctx.AddClass("file-list")

	return button
}

func (this *filesPanel) createOpenFolderButton(fileInformation *octoprint.FileInformation) *gtk.Button {
	image := utils.MustImageFromFileWithSize("open.svg", this.Scaled(40), this.Scaled(40))
	button := utils.MustButton(image, func() {
		this.locationHistory.GoForward(fileInformation.Name)
		this.doLoadFiles()
	})

	ctx, _ := button.GetStyleContext()
	ctx.AddClass("color1")
	ctx.AddClass("file-list")

	return button
}

func (this *filesPanel) createLoadAndPrintButton(imageFileName string, fileInformation *octoprint.FileInformation) *gtk.Button {
	button := utils.MustButton(
		utils.MustImageFromFileWithSize(imageFileName, this.Scaled(40), this.Scaled(40)),
		utils.MustConfirmDialogBox(this.UI.window, "Do you wish to proceed?", func() {
			selectFileRequest := &octoprint.SelectFileRequest{}
			selectFileRequest.Location = octoprint.Local
			selectFileRequest.Path = fileInformation.Path
			selectFileRequest.Print = true

			utils.Logger.Infof("Loading file %q", fileInformation.Name)
			if err := selectFileRequest.Do(this.UI.Client); err != nil {
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

func (this *filesPanel) addRootLocations() {
	this.addMessage("Select source location:")
	this.addRootLocation(octoprint.Local, 0)
	this.addRootLocation(octoprint.SDCard, 1)
}

func (this *filesPanel) addMessage(message string) {
	topBox := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	nameLabel := utils.MustLabel(message)
	nameLabel.SetMarkup(fmt.Sprintf("<big>%s</big>", utils.StrEllipsis(message)))
	nameLabel.SetHExpand(true)
	nameLabel.SetHAlign(gtk.ALIGN_START)

	labelsBox := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	labelsBox.Add(nameLabel)
	labelsBox.SetVExpand(false)
	labelsBox.SetVAlign(gtk.ALIGN_CENTER)
	labelsBox.SetHAlign(gtk.ALIGN_START)
	labelsBoxStyleContext, _ := labelsBox.GetStyleContext()
	labelsBoxStyleContext.AddClass("labels-box")

	topBox.Add(labelsBox)



	listItemBox := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	listItemBox.SetMarginTop(1)
	listItemBox.SetMarginBottom(1)
	listItemBox.SetMarginStart(15)
	listItemBox.SetMarginEnd(15)
	listItemBox.SetHExpand(true)
	listItemBox.Add(topBox)

	listItemFrame, _ := gtk.FrameNew("")
	listItemFrame.Add(listItemBox)


	this.listBox.Add(listItemFrame)
}


func (this *filesPanel) addRootLocation(location octoprint.Location, index int) {
	topBox := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)

	var itemImage *gtk.Image
	if location == octoprint.Local {
		itemImage = utils.MustImageFromFileWithSize("octoprint-tentacle.svg", this.Scaled(35), this.Scaled(35))
	} else {
		itemImage = utils.MustImageFromFileWithSize("sd.svg", this.Scaled(35), this.Scaled(35))
	}
	topBox.Add(itemImage)


	name := ""
	if location == octoprint.Local {
		name = "  OctoPrint"
	} else {
		name = "  SD Card"
	}
	nameLabel := utils.MustLabel(name)
	nameLabel.SetMarkup(fmt.Sprintf("<big>%s</big>", utils.StrEllipsis(name)))
	nameLabel.SetHExpand(true)
	nameLabel.SetHAlign(gtk.ALIGN_START)

	infoLabel := utils.MustLabel("")
	infoLabel.SetHAlign(gtk.ALIGN_START)
	infoLabel.SetMarkup("<small> </small>")

	labelsBox := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	labelsBox.Add(nameLabel)
	labelsBox.Add(infoLabel)
	labelsBox.SetVExpand(false)
	labelsBox.SetVAlign(gtk.ALIGN_CENTER)
	labelsBox.SetHAlign(gtk.ALIGN_START)
	labelsBoxStyleContext, _ := labelsBox.GetStyleContext()
	labelsBoxStyleContext.AddClass("labels-box")

	topBox.Add(labelsBox)


	var itemButton *gtk.Button
	if location == octoprint.Local {
		itemButton = this.createOpenLocationButton(location)
	} else {
		itemButton = this.createOpenLocationButton(location)
	}

	actionBox := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	actionBox.Add(itemButton)
	topBox.Add(actionBox)


	listItemBox := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	listItemBox.SetMarginTop(1)
	listItemBox.SetMarginBottom(1)
	listItemBox.SetMarginStart(15)
	listItemBox.SetMarginEnd(15)
	listItemBox.SetHExpand(true)
	listItemBox.Add(topBox)


	listItemFrame, _ := gtk.FrameNew("")
	listItemFrame.Add(listItemBox)

	this.listBox.Add(listItemFrame)
}



/*
func (this *filesPanel) isReady() bool {
	state, err := (&octoprint.SDStateRequest{}).Do(this.UI.Client)
	if err != nil {
		Logger.Error(err)
		return false
	}

	return state.Ready
}
*/
