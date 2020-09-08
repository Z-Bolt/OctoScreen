package uiWidgets

import (
	"fmt"

	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/utils"
)


type OctoPrintInfoBox struct {
	*SystemInfoBox
}

func CreateOctoPrintInfoBox(
	client				*octoprint.Client,
	logoWidth			int,
) *OctoPrintInfoBox {
	logoHeight := int(float64(logoWidth) * 1.25)
	logoImage := utils.MustImageFromFileWithSize("logos/logo-octoprint.png", logoWidth, logoHeight)

	versionResponse, err := (&octoprint.VersionRequest{}).Do(client)
	if err != nil {
		panic(err)
	}

	base := CreateSystemInfoBox(
		client,
		logoImage,
		"OctoPrint",
		versionResponse.Server,
		fmt.Sprintf("(API   %s)", versionResponse.API),
	)

	instance := &OctoPrintInfoBox {
		SystemInfoBox:			base,
	}

	return instance
}
