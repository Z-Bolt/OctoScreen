package utils

import (
	"fmt"
	"strings"
	"strconv"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

var cachedToolheadCount = -1

func GetToolheadCount(client *octoprint.Client) int {
	if cachedToolheadCount != -1 {
		return cachedToolheadCount
	}

	c, err := (&octoprint.ConnectionRequest{}).Do(client)
	if err != nil {
		LogError("toolhaad.GetToolheadCount()", "Do(ConnectionRequest)", err)
		return 0
	}

	profile, err := (&octoprint.PrinterProfilesRequest{Id: c.Current.PrinterProfile}).Do(client)
	if err != nil {
		LogError("toolhaad.GetToolheadCount()", "Do(PrinterProfilesRequest)", err)
		return 0
	}

	if profile.Extruder.SharedNozzle {
		cachedToolheadCount = 1
	}

	cachedToolheadCount = profile.Extruder.Count


	// TESTING: uncomment to force all toolheads to display and use for testing
	// cachedToolheadCount = 2
	// cachedToolheadCount = 3
	// cachedToolheadCount = 4


	return cachedToolheadCount
}


func GetDisplayNameForTool(toolName string) string {
	// Since this is such a hack, lets add some bounds checking
	if toolName == "" {
		Logger.Error("toolhaad.GetDisplayNameForTool() - toolName is empty")
		return ""
	}

	lowerCaseName := strings.ToLower(toolName)
	if strings.LastIndex(lowerCaseName, "tool") != 0 {
		Logger.Errorf("toolhaad.GetDisplayNameForTool() - toolName is invalid, value passed in was: %q", toolName)
		return ""
	}

	if len(toolName) != 5 {
		Logger.Errorf("toolhaad.GetDisplayNameForTool() - toolName is invalid, value passed in was: %q", toolName)
		return ""
	}

	toolIndexAsInt, _ := strconv.Atoi(string(toolName[4]))
	displayName := toolName[0:4]
	displayName = displayName + strconv.Itoa(toolIndexAsInt + 1)

	return displayName
}


func GetToolTarget(client *octoprint.Client, tool string) (float64, error) {
	fullStateRequest, err := (&octoprint.FullStateRequest{
		Exclude: []string{"sd", "state"},
	}).Do(client)

	if err != nil {
		LogError("tools.GetToolTarget()", "Do(StateRequest)", err)
		return -1, err
	}

	currentTemperatureData, ok := fullStateRequest.Temperature.CurrentTemperatureData[tool]
	if !ok {
		return -1, fmt.Errorf("unable to find tool %q", tool)
	}

	return currentTemperatureData.Target, nil
}


func SetToolTarget(client *octoprint.Client, tool string, target float64) error {
	if tool == "bed" {
		cmd := &octoprint.BedTargetRequest{Target: target}
		return cmd.Do(client)
	}

	cmd := &octoprint.ToolTargetRequest{Targets: map[string]float64{tool: target}}

	return cmd.Do(client)
}


func GetCurrentTemperatureData(client *octoprint.Client) (map[string]octoprint.TemperatureData, error) {
	temperatureDataResponse, err := (&octoprint.TemperatureDataRequest{}).Do(client)
	if err != nil {
		LogError("tools.GetCurrentTemperatureData()", "Do(TemperatureDataRequest)", err)
		return nil, err
	}

	if temperatureDataResponse == nil {
		Logger.Error("tools.GetCurrentTemperatureData() - temperatureDataResponse is nil")
		return nil, err
	}

	// Can't test for temperatureDataResponse.TemperatureStateResponse == nil (type mismatch)

	if temperatureDataResponse.TemperatureStateResponse.CurrentTemperatureData == nil {
		Logger.Error("tools.GetCurrentTemperatureData() - temperatureDataResponse.TemperatureStateResponse.CurrentTemperatureData is nil")
		return nil, err
	}

	return temperatureDataResponse.TemperatureStateResponse.CurrentTemperatureData, nil
}


func CurrentHotendTemperatureIsTooLow(client *octoprint.Client, extruderId, action string, parentWindow *gtk.Window) bool {
	currentTemperatureData, err := GetCurrentTemperatureData(client)
	if err != nil {
		LogError("tools.CurrentHotendTemperatureIsTooLow()", "GetCurrentTemperatureData()", err)
		return true
	}

	temperatureData := currentTemperatureData[extruderId]

	if HotendTemperatureIsTooLow(temperatureData, action, parentWindow) {
		LogError("tools.CurrentHotendTemperatureIsTooLow()", "HotendTemperatureIsTooLow()", err)
		return true
	}

	return false
}


func GetToolheadFileName(hotendIndex, hotendCount int) string {
	strImageFileName := ""
	if hotendIndex == 1 && hotendCount == 1 {
		strImageFileName = "toolhead.svg"
	} else {
		strImageFileName = fmt.Sprintf("toolhead-%d.svg", hotendIndex)
	}

	return strImageFileName
}

func GetHotendFileName(hotendIndex, hotendCount int) string {
	strImageFileName := ""
	if hotendIndex == 1 && hotendCount == 1 {
		strImageFileName = "hotend.svg"
	} else {
		strImageFileName = fmt.Sprintf("hotend-%d.svg", hotendIndex)
	}

	return strImageFileName
}

func GetNozzleFileName(hotendIndex, hotendCount int) string {
	strImageFileName := ""
	if hotendIndex == 1 && hotendCount == 1 {
		strImageFileName = "nozzle.svg"
	} else {
		strImageFileName = fmt.Sprintf("nozzle-%d.svg", hotendIndex)
	}

	return strImageFileName
}

func GetTemperatureDataString(temperatureData octoprint.TemperatureData) string {
	return fmt.Sprintf("%.0f°C / %.0f°C", temperatureData.Actual, temperatureData.Target)
}