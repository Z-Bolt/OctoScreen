package uiWidgets

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"

	"github.com/Z-Bolt/OctoScreen/logger"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)


func CreateFilesPreviewSubRow(
	fileResponse			*dataModels.FileResponse,
	fileSystemImageWidth	int,
	fileSystemImageHeight	int,
) *gtk.Box {
	previewSubRow := CreateVerticalLayoutBox()

	previewThumbnail := createPreviewThumbnail(
		fileResponse,
		fileSystemImageWidth,
		fileSystemImageHeight,
	)
	if previewThumbnail != nil {
		previewSubRow.Add(previewThumbnail)
	}

	return previewSubRow
}

func createPreviewThumbnail(
	fileResponse			*dataModels.FileResponse,
	fileSystemImageWidth	int,
	fileSystemImageHeight	int,
) *gtk.Box {
	if fileResponse.Thumbnail == "" {
		return nil;
	}

	logger.Debugf("FilesPreviewSubRow.createPreviewThumbnail() - fileResponse.Thumbnail is %s", fileResponse.Thumbnail)

	octoScreenConfig := utils.GetOctoScreenConfigInstance()
	octoPrintConfig := octoScreenConfig.OctoPrintConfig
	thumbnailUrl := fmt.Sprintf("%s/%s", octoPrintConfig.Server.Host, fileResponse.Thumbnail)
	logger.Debugf("FilesPreviewSubRow.createPreviewThumbnail() - thumbnailPath is: %q" , thumbnailUrl)

	previewImage, imageFromUrlErr := utils.ImageFromUrl(thumbnailUrl)
	if imageFromUrlErr != nil {
		return nil
	}

	logger.Debugf("FilesPreviewSubRow.createPreviewThumbnail() - no error from ImageNewFromPixbuf, now trying to add it...")

	bottomBox := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 0)

	// Initially was setting the horizontal alignment with CSS, but different resolutions
	// (eg 800x480 vs 480x320) didn't align correctly, so I added a blank SVG to offset
	// the preview thumbnail image.
	spacerImage := utils.MustImageFromFileWithSize("blank.svg", fileSystemImageWidth, fileSystemImageHeight)
	bottomBox.Add(spacerImage)

	// Still need some CSS for the bottom margin.
	previewImageStyleContext, _ := previewImage.GetStyleContext()
	previewImageStyleContext.AddClass("preview-image-list-item")

	// OK, now add the preview image.
	bottomBox.Add(previewImage)

	// ...and finally add everything to the bottom box/container.
	// listItemBox.Add(bottomBox)

	return bottomBox
}
