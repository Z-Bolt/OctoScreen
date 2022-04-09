package uiWidgets

import (
	"fmt"

	"github.com/Z-Bolt/OctoScreen/logger"
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
	logger.TraceEnter("OctoPrintInfoBox.CreateOctoPrintInfoBox()")

	logoHeight := int(float64(logoWidth) * 1.25)
	logoImage := utils.MustImageFromFileWithSize("logos/logo-octoprint.png", logoWidth, logoHeight)

	server := "Unknown?"
	apiVersion := "Unknown?"

	connectionManager := utils.GetConnectionManagerInstance(client)
	if connectionManager.IsConnectedToOctoPrint == true {
		// Only call if we're connected to OctoPrint
		versionResponse, err := (&octoprintApis.VersionRequest{}).Do(client)
		if err != nil {
			logger.LogError("OctoPrintInfoBox.CreateOctoPrintInfoBox()", "VersionRequest.Do()", err)
		} else if versionResponse != nil {
			server = versionResponse.Server
			apiVersion = versionResponse.API
		}
	}

	base := CreateSystemInfoBox(
		client,
		logoImage,
		"OctoPrint",
		server,
		fmt.Sprintf("(API   %s)", apiVersion),   // Use 3 spaces here... 1 space doesn't have enough kerning.
	)

	instance := &OctoPrintInfoBox {
		SystemInfoBox:			base,
	}

	logger.TraceLeave("OctoPrintInfoBox.CreateOctoPrintInfoBox()")

	return instance
}
