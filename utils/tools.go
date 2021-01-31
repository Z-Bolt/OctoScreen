package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)

var cachedToolheadCount = -1

func GetToolheadCount(client *octoprintApis.Client) int {
	if cachedToolheadCount != -1 {
		return cachedToolheadCount
	}

	connectionResponse, err := (&octoprintApis.ConnectionRequest{}).Do(client)
	if err != nil {
		LogError("Tools.GetToolheadCount()", "version.Get()", err)
		return 0
	}

	printerProfile, err := (&octoprintApis.PrinterProfilesRequest{Id: connectionResponse.Current.PrinterProfile}).Do(client)
	if err != nil {
		LogError("Tools.GetToolheadCount()", "Do(PrinterProfilesRequest)", err)
		return 0
	}

	cachedToolheadCount = printerProfile.Extruder.Count
	if printerProfile.Extruder.HasSharedNozzle {
		cachedToolheadCount = 1
	} else if cachedToolheadCount > 4 {
		cachedToolheadCount = 4
	}


	// TESTING: uncomment to force all toolheads to display and use for testing
	// cachedToolheadCount = 2
	// cachedToolheadCount = 3
	// cachedToolheadCount = 4


	return cachedToolheadCount
}


func GetDisplayNameForTool(toolName string) string {
	// Since this is such a hack, lets add some bounds checking
	if toolName == "" {
		Logger.Error("Tools..GetDisplayNameForTool() - toolName is empty")
		return ""
	}

	lowerCaseName := strings.ToLower(toolName)
	if strings.LastIndex(lowerCaseName, "tool") != 0 {
		Logger.Errorf("Tools.GetDisplayNameForTool() - toolName is invalid, value passed in was: %q", toolName)
		return ""
	}

	if len(toolName) != 5 {
		Logger.Errorf("Tools.GetDisplayNameForTool() - toolName is invalid, value passed in was: %q", toolName)
		return ""
	}

	toolIndexAsInt, _ := strconv.Atoi(string(toolName[4]))
	displayName := toolName[0:4]
	displayName = displayName + strconv.Itoa(toolIndexAsInt + 1)

	return displayName
}


func GetToolTarget(client *octoprintApis.Client, tool string) (float64, error) {
	Logger.Debug("entering Tools.GetToolTarget()")


	fullStateRespone, err := (&octoprintApis.FullStateRequest{
		Exclude: []string{"sd", "state"},
	}).Do(client)

	if err != nil {
		LogError("tools.GetToolTarget()", "Do(StateRequest)", err)

		Logger.Debug("leaving Tools.GetToolTarget()")
		return -1, err
	}

	currentTemperatureData, ok := fullStateRespone.Temperature.CurrentTemperatureData[tool]
	if !ok {
		Logger.Debug("leaving Tools.GetToolTarget()")
		return -1, fmt.Errorf("unable to find tool %q", tool)
	}

	Logger.Debug("leaving Tools.GetToolTarget()")
	return currentTemperatureData.Target, nil
}


func SetToolTarget(client *octoprintApis.Client, tool string, target float64) error {
	Logger.Debug("entering Tools.SetToolTarget()")

	if tool == "bed" {
		cmd := &octoprintApis.BedTargetRequest{Target: target}

		Logger.Debug("leaving Tools.SetToolTarget()")
		return cmd.Do(client)
	}

	cmd := &octoprintApis.ToolTargetRequest{Targets: map[string]float64{tool: target}}

	Logger.Debug("leaving Tools.SetToolTarget()")
	return cmd.Do(client)
}


func GetCurrentTemperatureData(client *octoprintApis.Client) (map[string]dataModels.TemperatureData, error) {
	Logger.Debug("entering Tools.GetCurrentTemperatureData()")

	temperatureDataResponse, err := (&octoprintApis.TemperatureDataRequest{}).Do(client)
	if err != nil {
		LogError("tools.GetCurrentTemperatureData()", "Do(TemperatureDataRequest)", err)

		Logger.Debug("leaving Tools.GetCurrentTemperatureData()")
		return nil, err
	}

	if temperatureDataResponse == nil {
		Logger.Error("tools.GetCurrentTemperatureData() - temperatureDataResponse is nil")

		Logger.Debug("leaving Tools.GetCurrentTemperatureData()")
		return nil, err
	}

	// Can't test for temperatureDataResponse.TemperatureStateResponse == nil (type mismatch)

	if temperatureDataResponse.TemperatureStateResponse.CurrentTemperatureData == nil {
		Logger.Error("tools.GetCurrentTemperatureData() - temperatureDataResponse.TemperatureStateResponse.CurrentTemperatureData is nil")

		Logger.Debug("leaving Tools.GetCurrentTemperatureData()")
		return nil, err
	}

	Logger.Debug("leaving Tools.GetCurrentTemperatureData()")
	return temperatureDataResponse.TemperatureStateResponse.CurrentTemperatureData, nil
}


func CheckIfHotendTemperatureIsTooLow(client *octoprintApis.Client, extruderId, action string, parentWindow *gtk.Window) bool {
	currentTemperatureData, err := GetCurrentTemperatureData(client)
	if err != nil {
		LogError("tools.CurrentHotendTemperatureIsTooLow()", "GetCurrentTemperatureData()", err)
		return true
	}

	temperatureData := currentTemperatureData[extruderId]

	// If the temperature of the hotend is too low, display an error.
	if HotendTemperatureIsTooLow(temperatureData, action, parentWindow) {
		errorMessage := fmt.Sprintf(
			"The temperature of the hotend is too low to %s.\n(the current temperature is only %.0f°C)\n\nPlease increase the temperature and try again.",
			action,
			temperatureData.Actual,
		)
		ErrorMessageDialogBox(parentWindow, errorMessage)

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

func GetTemperatureDataString(temperatureData dataModels.TemperatureData) string {
	return fmt.Sprintf("%.0f°C / %.0f°C", temperatureData.Actual, temperatureData.Target)
}


// TODO: maybe move HotendTemperatureIsTooLow into a hotend utils file?

const MIN_HOTEND_TEMPERATURE = 150.0

func HotendTemperatureIsTooLow(
	temperatureData			dataModels.TemperatureData,
	action					string,
	parentWindow			*gtk.Window,
) bool {
	targetTemperature := temperatureData.Target
	Logger.Infof("tools.HotendTemperatureIsTooLow() - targetTemperature is %.2f", targetTemperature)

	actualTemperature := temperatureData.Actual
	Logger.Infof("tools.HotendTemperatureIsTooLow() - actualTemperature is %.2f", actualTemperature)

	if targetTemperature <= MIN_HOTEND_TEMPERATURE || actualTemperature <= MIN_HOTEND_TEMPERATURE {
		return true
	}

	return false
}
