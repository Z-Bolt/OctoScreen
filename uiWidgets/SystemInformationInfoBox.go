package uiWidgets

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type SystemInformationInfoBox struct {
	*gtk.Box
}

func CreateSystemInformationInfoBox(
	client			*octoprint.Client,
) *SystemInformationInfoBox {
	base := utils.MustBox(gtk.ORIENTATION_VERTICAL, 0)
	base.SetVExpand(true)
	base.SetVAlign(gtk.ALIGN_CENTER)

	title := utils.MustLabel("<b>System Information</b>")
	title.SetMarginBottom(5)
	title.SetMarginTop(15)
	base.Add(title)

	virtualMemoryStat, _ := mem.VirtualMemory()
	memoryString := fmt.Sprintf(
		"Memory Total / Free: <b>%s / %s</b>",
		humanize.Bytes(virtualMemoryStat.Total),
		humanize.Bytes(virtualMemoryStat.Free),
	)
	base.Add(utils.MustLabel(memoryString))

	avgStat, _ := load.Avg()
	loadAverageString := fmt.Sprintf(
		"Load Average: <b>%.2f, %.2f, %.2f</b>",
		avgStat.Load1,
		avgStat.Load5,
		avgStat.Load15,
	)
	base.Add(utils.MustLabel(loadAverageString))

	instance := &SystemInformationInfoBox {
		Box:			base,
	}

	return instance
}
