package uiWidgets

import (
	// "fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type SystemInfoBox struct {
	*gtk.Box
}

func CreateSystemInfoBox(
	client				*octoprintApis.Client,
	image				*gtk.Image,
	str1				string,
	str2				string,
	str3				string,
) *SystemInfoBox {
	base := utils.MustBox(gtk.ORIENTATION_VERTICAL, 0)
	base.SetHExpand(true)
	base.SetHAlign(gtk.ALIGN_CENTER)
	base.SetVExpand(true)
	base.SetVAlign(gtk.ALIGN_CENTER)

	ctx, _ := image.GetStyleContext()
	ctx.AddClass("margin-top-5")
	base.Add(image)

	label1 := utils.MustLabel(str1)
	ctx, _ = label1.GetStyleContext()
	ctx.AddClass("margin-top-10")
	ctx.AddClass("font-size-18")
	base.Add(label1)

	label2 := utils.MustLabel(str2)
	ctx, _ = label2.GetStyleContext()
	ctx.AddClass("font-size-18")
	base.Add(label2)

	logLevel := logger.LogLevel()
	if logLevel == "debug" {
		label3 := utils.MustLabel(str3)
		ctx, _ = label3.GetStyleContext()
		ctx.AddClass("font-size-16")
		base.Add(label3)
	}

	instance := &SystemInfoBox {
		Box:			base,
	}

	return instance
}
