package uiWidgets

import (
	"fmt"
	"sync"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type ToolHeatup struct {
	*gtk.Button
	sync.RWMutex

	isHeating		bool
	tool			string
	printer			*octoprint.Client
}

func CreteToolHeatupButton(index int, printer *octoprint.Client) *ToolHeatup {
	var (
		image string
		tool  string
	)

	if index < 0 {
		image = "bed.svg"
		tool = "bed"
	} else if index == 0 {
		image = "toolhead.svg"
		tool = "tool0"
	} else {
		image = fmt.Sprintf("hotend-with-color-%d.svg", index)
		tool = fmt.Sprintf("tool%d", index - 1)
	}

	instance := &ToolHeatup{
		Button:  utils.MustButtonImage("", image, nil),
		tool:    tool,
		printer: printer,
	}

	_, err := instance.Connect("clicked", instance.clicked)
	if err != nil {
		utils.LogError("idle_status.creteToolHeatupButton()", "t.Connect('clicked', t.clicked)", err)
	}

	return instance
}

func (this *ToolHeatup) UpdateStatus(heating bool) {
	ctx, _ := this.GetStyleContext()
	if heating {
		ctx.AddClass("active")
	} else {
		ctx.RemoveClass("active")
	}

	this.isHeating = heating
}

func (this *ToolHeatup) SetTemperatures(temperatureData octoprint.TemperatureData) {
	text := utils.GetTemperatureDataString(temperatureData)
	this.SetLabel(text)
	this.UpdateStatus(temperatureData.Target > 0)
}

func (this *ToolHeatup) GetProfileTemperature() float64 {
	temperature := 0.0

	settingsResponse, err := (&octoprint.SettingsRequest{}).Do(this.printer)
	if err != nil {
		utils.LogError("idle_status.getProfileTemperature()", "Do(SettingsRequest)", err)
		return 0
	}

	if len(settingsResponse.Temperature.TemperaturePresets) > 0 {
		if this.tool == "bed" {
			temperature = settingsResponse.Temperature.TemperaturePresets[0].Bed
		} else {
			temperature = settingsResponse.Temperature.TemperaturePresets[0].Extruder
		}
	} else {
		if this.tool == "bed" {
			temperature = 75
		} else {
			temperature = 220
		}
	}

	return temperature
}

func (this *ToolHeatup) clicked() {
	defer func() {
		this.UpdateStatus(!this.isHeating)
	}()

	var (
		target float64
		err    error
	)

	if this.isHeating {
		target = 0.0
	} else {
		target = this.GetProfileTemperature()
	}

	if this.tool == "bed" {
		cmd := &octoprint.BedTargetRequest{Target: target}
		err = cmd.Do(this.printer)
		if err != nil {
			utils.LogError("idle_status.clicked()", "Do(BedTargetRequest)", err)
		}
	} else {
		cmd := &octoprint.ToolTargetRequest{Targets: map[string]float64{this.tool: target}}
		err = cmd.Do(this.printer)
		if err != nil {
			utils.LogError("idle_status.clicked()", "Do(ToolTargetRequest)", err)
		}
	}
}
