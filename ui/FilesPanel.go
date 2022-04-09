package ui

import (
	"fmt"
	// "os"
	"sort"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gotk3/gotk3/gtk"

	// "github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	// "github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)


type filesPanel struct {
	CommonPanel

	listBox				*gtk.Box
	refreshButton		*gtk.Button
	backButton			*gtk.Button
	locationHistory		utils.LocationHistory
}

var filesPanelInstance *filesPanel

func GetFilesPanelInstance(
	ui					*UI,
) *filesPanel {
	if filesPanelInstance == nil {
		locationHistory := utils.LocationHistory {
			Locations: []dataModels.Location{},
		}

		instance := &filesPanel {
			CommonPanel: CreateCommonPanel("FilesPanel", ui),
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
	box.Add(this.createActionFooter())
	this.Grid().Add(box)

	this.doLoadFiles()
}

func (this *filesPanel) createActionFooter() *gtk.Box {
	actionBar := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	actionBar.SetHAlign(gtk.ALIGN_END)
	actionBar.SetHExpand(true)
	actionBar.SetMarginTop(5)
	actionBar.SetMarginBottom(5)
	actionBar.SetMarginEnd(5)

	this.refreshButton = this.createRefreshButton()
	actionBar.Add(this.refreshButton)

	this.backButton = this.createBackButton()
	actionBar.Add(this.backButton)

	return actionBar
}

func (this *filesPanel) createRefreshButton() *gtk.Button {
	image := utils.MustImageFromFileWithSize("refresh.svg", this.Scaled(40), this.Scaled(40))
	return utils.MustButton(image, this.doLoadFiles)
}

func (this *filesPanel) createBackButton() *gtk.Button {
	image := utils.MustImageFromFileWithSize("back.svg", this.Scaled(40), this.Scaled(40))
	return utils.MustButton(image, this.goBack)
}

func (this *filesPanel) doLoadFiles() {
	utils.EmptyTheContainer(&this.listBox.Container)
	atRootLevel := this.displayRootLocations()
	/*
	 * If we are at `root` (display the option for SD AND Local), but SD is not
	 * ready, push us up and into Local so the user doesn't have to work harder
	 * than they have to.
	 */
	if atRootLevel && !this.sdIsReady() {
		atRootLevel = false
		this.locationHistory = utils.LocationHistory {
			Locations: []dataModels.Location{dataModels.Local},
		}
	}
	if atRootLevel {
		this.addRootLocations()
	} else {
		sortedFiles := this.getSortedFiles()
		this.addSortedFiles(sortedFiles)
	}
	this.listBox.ShowAll()
}

func (this *filesPanel) sdIsReady() bool {
	err := (&octoprintApis.SdRefreshRequest {}).Do(this.UI.Client)
	if err == nil {
		sdState, err := (&octoprintApis.SdStateRequest {}).Do(this.UI.Client)
		if err == nil && sdState.IsReady == true {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (this *filesPanel) goBack() {
	if this.displayRootLocations() {
		this.UI.GoToPreviousPanel()
	} else if this.locationHistory.IsRoot() {
		this.locationHistory.GoBack()
		if this.sdIsReady() {
			this.doLoadFiles()
		} else {
			this.UI.GoToPreviousPanel()
		}
	} else {
		this.locationHistory.GoBack()
		this.doLoadFiles()
	}
}

func (this *filesPanel) displayRootLocations() bool {
	if this.locationHistory.Length() < 1 {
		return true
	} else {
		return false
	}
}

func (this *filesPanel) getSortedFiles() []*dataModels.FileResponse {
	var files []*dataModels.FileResponse

	if this.displayRootLocations() {
		return nil
	}

	current := this.locationHistory.CurrentLocation()
	logger.Infof("Loading list of files from: %s", string(current))

	if current == dataModels.SDCard {
		sdRefreshRequest := &octoprintApis.SdRefreshRequest {}
		err := sdRefreshRequest.Do(this.UI.Client)
		if err != nil {
			logger.LogError("getSortedFiles()", "sdRefreshRequest.Do()", err)
			return []*dataModels.FileResponse{}
		} else {
			// Pause here for a second, because the preceding call to filesRequest.Do()
			// doesn't work, and it returns a truncated list of files.  Pausing here
			// for a second seems to resolve the issue.
			time.Sleep(1 * time.Second)
		}
	}

	filesRequest := &octoprintApis.FilesRequest {
		Location: current,
		Recursive: false,
	}
	filesResponse, err := filesRequest.Do(this.UI.Client)
	if err != nil {
		logger.LogError("files.getSortedFiles()", "Do(FilesRequest)", err)
		files = []*dataModels.FileResponse{}
	} else {
		files = filesResponse.Files
	}

	var filteredFiles []*dataModels.FileResponse
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
	this.addRootLocation(dataModels.Local)
	this.addRootLocation(dataModels.SDCard)
}

func (this *filesPanel) addMessage(message string) {
	nameLabel := this.createNameLabel(message)
	labelsBox := this.createLabelsBox(nameLabel, nil)
	labelsBox.SetMarginStart(10)

	topBox := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	topBox.Add(labelsBox)

	listItemBox := this.createListItemBox()
	listItemBox.Add(topBox)

	listItemFrame, _ := gtk.FrameNew("")
	listItemFrame.Add(listItemBox)

	this.listBox.Add(listItemFrame)
}

func (this *filesPanel) addRootLocation(location dataModels.Location) {
	rootLocationButton := this.createRootLocationButton(location)

	listBoxRow, _ := gtk.ListBoxRowNew()
	listBoxRow.Add(rootLocationButton)

	this.listBox.Add(listBoxRow)
}

func (this *filesPanel) createRootLocationButton(location dataModels.Location) *gtk.Button {
	var itemImage *gtk.Image
	if location == dataModels.Local {
		itemImage = utils.MustImageFromFileWithSize("logos/octoprint-tentacle.svg", this.Scaled(35), this.Scaled(35))
	} else {
		itemImage = utils.MustImageFromFileWithSize("sd.svg", this.Scaled(35), this.Scaled(35))
	}

	topBox := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	topBox.Add(itemImage)

	name := ""
	if location == dataModels.Local {
		name = "  OctoPrint"
	} else {
		name = "  SD Card"
	}
	nameLabel := this.createNameLabel(name)

	infoLabel := utils.MustLabel("")
	infoLabel.SetHAlign(gtk.ALIGN_START)
	infoLabel.SetMarkup("<small> </small>")

	labelsBox := this.createLabelsBox(nameLabel, infoLabel)
	topBox.Add(labelsBox)


	var actionImage *gtk.Image
	if location == dataModels.Local {
		actionImage = this.createOpenLocationImage(0)
	} else {
		actionImage = this.createOpenLocationImage(1)
	}

	actionBox := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	actionBox.Add(actionImage)
	topBox.Add(actionBox)

	rootLocationButton, _ := gtk.ButtonNew()
	rootLocationButton.Connect("clicked", func() {
		this.locationHistory = utils.LocationHistory {
			Locations: []dataModels.Location{location},
		}

		this.doLoadFiles()
	})

	rootLocationButton.Add(topBox)

	return rootLocationButton
}

func (this *filesPanel) addSortedFiles(sortedFiles []*dataModels.FileResponse) {
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

func (this *filesPanel) addItem(
	fileResponse *dataModels.FileResponse,
	index int,
) {
	/*
		Object hierarchy:

		listBox
			listBoxRow
				listItemButton (to handle to click for the entire item amd all of the child controls)
					listItemBox (to layout the objects, in this case the two rows within the button)
						infoAndActionRow (a Box)
						previewRow (a Box)
	*/


	listItemBox := this.createListItemBox()

	isFolder := fileResponse.IsFolder()
	infoAndActionRow := this.createInfoAndActionRow(fileResponse, index, isFolder)
	listItemBox.Add(infoAndActionRow)
	if !isFolder {
		previewRow := this.createPreviewRow(fileResponse)
		listItemBox.Add(previewRow)
	}

	listItemButton := this.createListItemButton(fileResponse, index, isFolder)
	listItemButton.Add(listItemBox)

	listBoxRow, _ := gtk.ListBoxRowNew()
	listBoxRowStyleContext, _ := listBoxRow.GetStyleContext()
	listBoxRowStyleContext.AddClass("list-box-row")
	if index % 2 != 0 {
		listBoxRowStyleContext.AddClass("list-item-nth-child-even")
	}
	listBoxRow.Add(listItemButton)

	this.listBox.Add(listBoxRow)
}

func (this *filesPanel) createInfoAndActionRow(
	fileResponse *dataModels.FileResponse,
	index int,
	isFolder bool,
) *gtk.Box {
	infoAndActionRow := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)


	// Column 1: Folder or File icon
	var itemImage *gtk.Image
	if isFolder {
		itemImage = utils.MustImageFromFileWithSize("folder.svg", this.Scaled(35), this.Scaled(35))
	} else {
		itemImage = utils.MustImageFromFileWithSize("file-gcode.svg", this.Scaled(35), this.Scaled(35))
	}
	infoAndActionRow.Add(itemImage)


	// Column 2: File name and file info
	name := fileResponse.Name
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

	labelsBox := this.createLabelsBox(nameLabel, infoLabel)
	infoAndActionRow.Add(labelsBox)


	// Column 3: printer image
	var actionImage *gtk.Image
	if isFolder {
		actionImage = this.createOpenLocationImage(index)
	} else {
		actionImage = this.createPrintImage()
	}

	actionBox := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	actionBox.Add(actionImage)

	infoAndActionRow.Add(actionBox)

	return infoAndActionRow
}

func (this *filesPanel) createPreviewRow(fileResponse *dataModels.FileResponse) *gtk.Box {
	previewRow := this.createListItemBox()
	this.addThumbnail(fileResponse, previewRow)

	return previewRow
}

func (this *filesPanel) createListItemButton(
	fileResponse *dataModels.FileResponse,
	index int,
	isFolder bool,
) *gtk.Button {
	listItemButton, _ := gtk.ButtonNew()
	listItemButtonStyleContext, _ := listItemButton.GetStyleContext()
	listItemButtonStyleContext.AddClass("list-item-button")
	if index % 2 != 0 {
		listItemButtonStyleContext.AddClass("list-item-nth-child-even")
	}

	if isFolder {
		listItemButton.Connect("clicked", func() {
			this.locationHistory.GoForward(fileResponse.Name)
			this.doLoadFiles()
		})
	} else {
		message := ""
		strLen := len(fileResponse.Name)
		if strLen <= 20 {
			message = fmt.Sprintf("Do you wish to print %s?", fileResponse.Name)
		} else {
			truncatedFileName := utils.TruncateString(fileResponse.Name, 40)
			message = fmt.Sprintf("Do you wish to print\n%s?", truncatedFileName)
		}

		listItemButton.Connect("clicked", utils.MustConfirmDialogBox(this.UI.window, message, func() {
			selectFileRequest := &octoprintApis.SelectFileRequest{}

			// Set the location to "local" or "sdcard"
			selectFileRequest.Location = this.locationHistory.Locations[0]

			selectFileRequest.Path = fileResponse.Path
			selectFileRequest.Print = true

			logger.Infof("Loading file %q", fileResponse.Name)
			if err := selectFileRequest.Do(this.UI.Client); err != nil {
				logger.LogError("FilesPanel.createLoadAndPrintButton()", "Do(SelectFileRequest)", err)
				return
			}
		}))
	}

	return listItemButton
}

func (this *filesPanel) createLabelsBox(
	nameLabel *gtk.Label,
	infoLabel *gtk.Label,
) *gtk.Box {
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
	nameLabel.SetMarkup(fmt.Sprintf("<big>%s</big>", utils.TruncateString(name, 30)))
	nameLabel.SetHExpand(true)
	nameLabel.SetHAlign(gtk.ALIGN_START)

	return nameLabel
}

func (this *filesPanel) createListItemBox() *gtk.Box {
	listItemBox := utils.MustBox(gtk.ORIENTATION_VERTICAL, 0)
	listItemBox.SetMarginTop(0)
	listItemBox.SetMarginBottom(0)
	listItemBox.SetMarginStart(0)
	listItemBox.SetMarginEnd(0)

	return listItemBox
}

func (this *filesPanel) addThumbnail(
	fileResponse *dataModels.FileResponse,
	listItemBox *gtk.Box,
) {
	if fileResponse.Thumbnail != "" {
		logger.Debugf("FilesPanel.addItem() - fileResponse.Thumbnail is %s", fileResponse.Thumbnail)

		octoScreenConfig := utils.GetOctoScreenConfigInstance()
		octoPrintConfig := octoScreenConfig.OctoPrintConfig
		thumbnailUrl := fmt.Sprintf("%s/%s", octoPrintConfig.Server.Host, fileResponse.Thumbnail)
		logger.Debugf("FilesPanel.addItem() - thumbnailPath is: %q" , thumbnailUrl)

		previewImage, imageFromUrlErr := utils.ImageFromUrl(thumbnailUrl)
		if imageFromUrlErr == nil {
			logger.Debugf("FilesPanel.addItem() - no error from ImageNewFromPixbuf, now trying to add it...")

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

func (this *filesPanel) createOpenLocationImage(index int) *gtk.Image {
	colorClass := fmt.Sprintf("color%d", (index % 4) + 1)

	return this.createActionImage("open.svg", colorClass)
}

func (this *filesPanel) createPrintImage() *gtk.Image {
	return this.createActionImage("print.svg", "color-warning-sign-yellow")
}

func (this *filesPanel) createActionImage(
	imageFileName string,
	colorClass string,
) *gtk.Image {
	image := utils.MustImageFromFileWithSize(
		imageFileName,
		this.Scaled(40),
		this.Scaled(40),
	)

	imageStyleContext, _ := image.GetStyleContext()
	imageStyleContext.AddClass(colorClass)

	return image
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
