package uiWidgets

import (
	// "fmt"

	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type OctoScreenPluginInfoBox struct {
	*SystemInfoBox
}

func CreateOctoScreenPluginInfoBox(
	client							*octoprintApis.Client,
	octoPrintPluginIsInstalled		bool,
) *OctoScreenPluginInfoBox {
	logoImage := utils.MustImageFromFile("logos/puzzle-piece.png")
	str1 := "OctoScreen plugin"

	str2 := ""
	if octoPrintPluginIsInstalled {
		getPluginManagerInfoResponse, err := (&octoprintApis.GetPluginManagerInfoRequest{}).Do(client)
		if err != nil {
			panic(err)
		}

		found := false
		for i := 0; i < len(getPluginManagerInfoResponse.Plugins) && !found; i++ {
			plugin := getPluginManagerInfoResponse.Plugins[i]
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

	return instance
}
