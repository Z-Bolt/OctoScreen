package ui

import (
	"fmt"
	"os"
	"sort"
	"strings"

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

func (this *filesPanel) getSortedFiles() []*octoprint.FileResponse {
	var files []*octoprint.FileResponse

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
		files = []*octoprint.FileResponse{}
	} else {
		files = filesResponse.Files
	}

	var filteredFiles []*octoprint.FileResponse
	for i := range files {
		if !strings.HasPrefix(files[i].Path, "trash") {
			filteredFiles = append(filteredFiles, files[i])
		}
	}

	sortedFiles := utils.FileResponsesSortedByDate(filteredFiles)
	// sortedFiles := utils.FileResponsesSortedByName(filteredFiles)
	sort.Sort(sortedFiles)

	return sortedFiles
}

func (this *filesPanel) addRootLocations() {
	this.addMessage("Select source location:")
	this.addRootLocation(octoprint.Local, 0)
	this.addRootLocation(octoprint.SDCard, 1)
}

func (this *filesPanel) addMessage(message string) {
	// nameLabel := utils.MustLabel(message)
	// nameLabel.SetMarkup(fmt.Sprintf("<big>%s</big>", utils.StrEllipsis(message)))
	// nameLabel.SetHExpand(true)
	// nameLabel.SetHAlign(gtk.ALIGN_START)
	nameLabel := this.createNameLabel(message)

	// labelsBox := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	// labelsBox.Add(nameLabel)
	// labelsBox.SetVExpand(false)
	// labelsBox.SetVAlign(gtk.ALIGN_CENTER)
	// labelsBox.SetHAlign(gtk.ALIGN_START)
	// labelsBoxStyleContext, _ := labelsBox.GetStyleContext()
	// labelsBoxStyleContext.AddClass("labels-box")
	labelsBox := this.createLabelsBox(nameLabel, nil)

	topBox := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	topBox.Add(labelsBox)

	listItemBox := this.createListItemBox()
	listItemBox.Add(topBox)

	listItemFrame, _ := gtk.FrameNew("")
	listItemFrame.Add(listItemBox)

	this.listBox.Add(listItemFrame)
}

func (this *filesPanel) addRootLocation(location octoprint.Location, index int) {
	var itemImage *gtk.Image
	if location == octoprint.Local {
		itemImage = utils.MustImageFromFileWithSize("octoprint-tentacle.svg", this.Scaled(35), this.Scaled(35))
	} else {
		itemImage = utils.MustImageFromFileWithSize("sd.svg", this.Scaled(35), this.Scaled(35))
	}

	topBox := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	topBox.Add(itemImage)

	name := ""
	if location == octoprint.Local {
		name = "  OctoPrint"
	} else {
		name = "  SD Card"
	}
	// nameLabel := utils.MustLabel(name)
	// nameLabel.SetMarkup(fmt.Sprintf("<big>%s</big>", utils.StrEllipsis(name)))
	// nameLabel.SetHExpand(true)
	// nameLabel.SetHAlign(gtk.ALIGN_START)
	nameLabel := this.createNameLabel(name)

	infoLabel := utils.MustLabel("")
	infoLabel.SetHAlign(gtk.ALIGN_START)
	infoLabel.SetMarkup("<small> </small>")

	// labelsBox := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	// labelsBox.Add(nameLabel)
	// labelsBox.Add(infoLabel)
	// labelsBox.SetVExpand(false)
	// labelsBox.SetVAlign(gtk.ALIGN_CENTER)
	// labelsBox.SetHAlign(gtk.ALIGN_START)
	// labelsBoxStyleContext, _ := labelsBox.GetStyleContext()
	// labelsBoxStyleContext.AddClass("labels-box")
	labelsBox := this.createLabelsBox(nameLabel, infoLabel)

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


	listItemBox := this.createListItemBox()
	listItemBox.Add(topBox)


	listItemFrame, _ := gtk.FrameNew("")
	listItemFrame.Add(listItemBox)

	this.listBox.Add(listItemFrame)
}

func (this *filesPanel) addSortedFiles(sortedFiles []*octoprint.FileResponse) {
	var index int = 0

	for _, fileResponse := range sortedFiles {
		if fileResponse.IsFolder() {
			this.addItem(fileResponse, index)
			index++
		}
	}

	for _, fileResponse := range sortedFiles {
		if !fileResponse.IsFolder() {
			this.addItem(fileResponse, index)
			index++
		}
	}
}


func (this *filesPanel) addItem(fileResponse *octoprint.FileResponse, index int) {
	isFolder := fileResponse.IsFolder()

	var itemImage *gtk.Image
	if isFolder {
		itemImage = utils.MustImageFromFileWithSize("folder.svg", this.Scaled(35), this.Scaled(35))
	} else {
		itemImage = utils.MustImageFromFileWithSize("file-gcode.svg", this.Scaled(35), this.Scaled(35))
	}

	topBox := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	topBox.Add(itemImage)


	name := fileResponse.Name
	// nameLabel := utils.MustLabel(name)
	// nameLabel.SetMarkup(fmt.Sprintf("<big>%s</big>", utils.TruncateString(name, 30)))
	// nameLabel.SetHExpand(true)
	// nameLabel.SetHAlign(gtk.ALIGN_START)
	nameLabel := this.createNameLabel(name)

	infoLabel := utils.MustLabel("")
	infoLabel.SetHAlign(gtk.ALIGN_START)
	if isFolder {
		infoLabel.SetMarkup(fmt.Sprintf("<small>Size: <b>%s</b></small>",
			humanize.Bytes(uint64(fileResponse.Size)),
		))
	} else {
		infoLabel.SetMarkup(fmt.Sprintf("<small>Uploaded: <b>%s</b> - Size: <b>%s</b></small>",
			humanize.Time(fileResponse.Date.Time), humanize.Bytes(uint64(fileResponse.Size)),
		))
	}

	// labelsBox := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	// labelsBox.Add(nameLabel)
	// labelsBox.Add(infoLabel)
	// labelsBox.SetVExpand(false)
	// labelsBox.SetVAlign(gtk.ALIGN_CENTER)
	// labelsBox.SetHAlign(gtk.ALIGN_START)
	// labelsBoxStyleContext, _ := labelsBox.GetStyleContext()
	// labelsBoxStyleContext.AddClass("labels-box")
	labelsBox := this.createLabelsBox(nameLabel, infoLabel)

	topBox.Add(labelsBox)




	var itemButton *gtk.Button
	if isFolder {
		itemButton = this.createOpenFolderButton(fileResponse)
	} else {
		itemButton = this.createLoadAndPrintButton("print.svg", fileResponse)
	}

	actionBox := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	actionBox.Add(itemButton)
	topBox.Add(actionBox)




	listItemBox := this.createListItemBox()
	listItemBox.Add(topBox)






	if !isFolder {
		this.addThumbnail(fileResponse, listItemBox)
	}




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



func (this *filesPanel) createLabelsBox(nameLabel, infoLabel *gtk.Label) *gtk.Box {
	labelsBox := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	if nameLabel != nil {
		labelsBox.Add(nameLabel)
	}
	if infoLabel != nil {
		labelsBox.Add(infoLabel)
	}
	labelsBox.SetVExpand(false)
	labelsBox.SetVAlign(gtk.ALIGN_CENTER)
	labelsBox.SetHAlign(gtk.ALIGN_START)
	labelsBoxStyleContext, _ := labelsBox.GetStyleContext()
	labelsBoxStyleContext.AddClass("labels-box")

	return labelsBox
}

func (this *filesPanel) createNameLabel(name string) *gtk.Label {
	nameLabel := utils.MustLabel(name)
	// nameLabel.SetMarkup(fmt.Sprintf("<big>%s</big>", utils.StrEllipsis(name)))
	nameLabel.SetMarkup(fmt.Sprintf("<big>%s</big>", utils.TruncateString(name, 30)))
	nameLabel.SetHExpand(true)
	nameLabel.SetHAlign(gtk.ALIGN_START)

	return nameLabel
}

func (this *filesPanel) createListItemBox() *gtk.Box {
	listItemBox := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	listItemBox.SetMarginTop(1)
	listItemBox.SetMarginBottom(1)
	listItemBox.SetMarginStart(15)
	listItemBox.SetMarginEnd(15)
	listItemBox.SetHExpand(true)

	return listItemBox
}







func (this *filesPanel) addThumbnail(fileResponse *octoprint.FileResponse, listItemBox *gtk.Box) {
	if fileResponse.Thumbnail != "" {
		utils.Logger.Infof("FilesPanel.addItem() - fileResponse.Thumbnail is %s", fileResponse.Thumbnail)

		thumbnailUrl := fmt.Sprintf("%s/%s", os.Getenv(utils.EnvBaseURL), fileResponse.Thumbnail)
		utils.Logger.Infof("FilesPanel.addItem() - thumbnailPath is: %q" , thumbnailUrl)

		previewImage, imageFromUrlErr := utils.ImageFromUrl(thumbnailUrl)
		if imageFromUrlErr == nil {
			utils.Logger.Infof("FilesPanel.addItem() - no error from ImageNewFromPixbuf, now trying to add it...")

			bottomBox := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 0)

			// Initially was setting the horizontal alignment with CSS, but different resolutions
			// (eg 800x480 vs 480x320) didn't align correctly, so I added a blank SVG to offset
			// the preview thumbnail image.
			spacerImage := utils.MustImageFromFileWithSize("blank.svg", this.Scaled(35), this.Scaled(35))
			bottomBox.Add(spacerImage)

			// Still need some CSS for the bottom margin.
			previewImageStyleContext, _ := previewImage.GetStyleContext()
			previewImageStyleContext.AddClass("preview-image-list-item")

			// OK, now add the preview image.
			bottomBox.Add(previewImage)

			// ...and finally add everything to the bottom box/container.
			listItemBox.Add(bottomBox)
		}
	}
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

func (this *filesPanel) createOpenFolderButton(fileResponse *octoprint.FileResponse) *gtk.Button {
	image := utils.MustImageFromFileWithSize("open.svg", this.Scaled(40), this.Scaled(40))
	button := utils.MustButton(image, func() {
		this.locationHistory.GoForward(fileResponse.Name)
		this.doLoadFiles()
	})

	ctx, _ := button.GetStyleContext()
	ctx.AddClass("color1")
	ctx.AddClass("file-list")

	return button
}

func (this *filesPanel) createLoadAndPrintButton(imageFileName string, fileResponse *octoprint.FileResponse) *gtk.Button {
	button := utils.MustButton(
		utils.MustImageFromFileWithSize(imageFileName, this.Scaled(40), this.Scaled(40)),
		utils.MustConfirmDialogBox(this.UI.window, "Do you wish to proceed?", func() {
			selectFileRequest := &octoprint.SelectFileRequest{}
			selectFileRequest.Location = octoprint.Local
			selectFileRequest.Path = fileResponse.Path
			selectFileRequest.Print = true

			utils.Logger.Infof("Loading file %q", fileResponse.Name)
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
