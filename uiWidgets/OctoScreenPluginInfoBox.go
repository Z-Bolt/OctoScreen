package uiWidgets

import (
	// "fmt"

	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type OctoScreenPluginInfoBox struct {
	*SystemInfoBox
}

func CreateOctoScreenPluginInfoBox(
	client							*octoprintApis.Client,
	uiState							string,
	octoPrintPluginIsInstalled		bool,
) *OctoScreenPluginInfoBox {
	logger.TraceEnter("OctoScreenPluginInfoBox.CreateOctoScreenPluginInfoBox()")

	logoImage := utils.MustImageFromFile("logos/puzzle-piece.png")
	str1 := "OctoScreen plugin"

	str2 := ""
	if octoPrintPluginIsInstalled {
		pluginManagerInfoResponse, err := (&octoprintApis.PluginManagerInfoRequest{}).Do(client, uiState)
		if err != nil {
			logger.LogError("CreateOctoScreenPluginInfoBox()", "PluginManagerInfoRequest.Do()", err)
			str2 = "Error"
		} else {
			found := false
			for i := 0; i < len(pluginManagerInfoResponse.Plugins) && !found; i++ {
				plugin := pluginManagerInfoResponse.Plugins[i]
				if plugin.Key == "zbolt_octoscreen" {
					found = true
					str2 = plugin.Version
				}
			}

			if !found {
				// OK, the plugin is there, we just can't get the info from a GET request.
				// Default to displaying, "Present"
				str2 = "Present"
			}
		}
	} else {
		str2 = "Not installed"
	}

	base := CreateSystemInfoBox(
		client,
		logoImage,
		str1,
		str2,
		"",
	)

	instance := &OctoScreenPluginInfoBox {
		SystemInfoBox:			base,
	}

	logger.TraceLeave("OctoScreenPluginInfoBox.CreateOctoScreenPluginInfoBox()")

	return instance
}
