package uiWidgets

import (
	"fmt"
	// "log"

	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)


type OctoPrintInfoBox struct {
	*SystemInfoBox
}

func CreateOctoPrintInfoBox(
	client				*octoprintApis.Client,
	logoWidth			int,
) *OctoPrintInfoBox {
	logoHeight := int(float64(logoWidth) * 1.25)
	logoImage := utils.MustImageFromFileWithSize("logos/logo-octoprint.png", logoWidth, logoHeight)

	versionResponse, err := (&octoprintApis.VersionRequest{}).Do(client)
	if err != nil {
		// TODO: should the error really trigger a panic?
		panic(err)
	}

	base := CreateSystemInfoBox(
		client,
		logoImage,
		"OctoPrint",
		versionResponse.Server,
		fmt.Sprintf("(API   %s)", versionResponse.API),   // Use 3 spaces... 1 space doesn't have enough kerning.
	)

	instance := &OctoPrintInfoBox {
		SystemInfoBox:			base,
	}

	return instance
}
