package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)

var cachedExtruderCount = -1
var cachedHasSharedNozzle = false

func getCachedPrinterProfileData(client *octoprintApis.Client) {
	if cachedExtruderCount != -1 {
		return
	}

	connectionResponse, err := (&octoprintApis.ConnectionRequest{}).Do(client)
	if err != nil {
		logger.LogError("Tools.setCachedPrinterProfileData()", "version.Get()", err)
		return
	}

	printerProfile, err := (&octoprintApis.PrinterProfilesRequest{Id: connectionResponse.Current.PrinterProfile}).Do(client)
	if err != nil {
		logger.LogError("Tools.setCachedPrinterProfileData()", "Do(PrinterProfilesRequest)", err)
		return
	}

	cachedExtruderCount = printerProfile.Extruder.Count
	if cachedExtruderCount > 4 {
		cachedExtruderCount = 4
	}

	cachedHasSharedNozzle = printerProfile.Extruder.HasSharedNozzle
}


func GetExtruderCount(client *octoprintApis.Client) int {
	if cachedExtruderCount == -1 {
		getCachedPrinterProfileData(client)
	}

	return cachedExtruderCount
}

func GetHotendCount(client *octoprintApis.Client) int {
	if cachedExtruderCount == -1 {
		getCachedPrinterProfileData(client)
	}

	if cachedHasSharedNozzle {
		return 1
	} else if cachedExtruderCount > 4 {
		return 4
	}

	return cachedExtruderCount
}

func GetHasSharedNozzle(client *octoprintApis.Client) bool {
	if cachedExtruderCount == -1 {
		getCachedPrinterProfileData(client)
	}

	return cachedHasSharedNozzle
}








func GetDisplayNameForTool(toolName string) string {
	// Since this is such a hack, lets add some bounds checking
	if toolName == "" {
		logger.Error("Tools..GetDisplayNameForTool() - toolName is empty")
		return ""
	}

	lowerCaseName := strings.ToLower(toolName)
	if strings.LastIndex(lowerCaseName, "tool") != 0 {
		logger.Errorf("Tools.GetDisplayNameForTool() - toolName is invalid, value passed in was: %q", toolName)
		return ""
	}

	if len(toolName) != 5 {
		logger.Errorf("Tools.GetDisplayNameForTool() - toolName is invalid, value passed in was: %q", toolName)
		return ""
	}

	toolIndexAsInt, _ := strconv.Atoi(string(toolName[4]))
	displayName := toolName[0:4]
	displayName = displayName + strconv.Itoa(toolIndexAsInt + 1)

	return displayName
}


func GetToolTarget(client *octoprintApis.Client, tool string) (float64, error) {
	logger.TraceEnter("Tools.GetToolTarget()")

	fullStateRespone, err := (&octoprintApis.FullStateRequest{
		Exclude: []string{"sd", "state"},
	}).Do(client)

	if err != nil {
		logger.LogError("tools.GetToolTarget()", "Do(StateRequest)", err)
		logger.TraceLeave("Tools.GetToolTarget()")
		return -1, err
	}

	currentTemperatureData, ok := fullStateRespone.Temperature.CurrentTemperatureData[tool]
	if !ok {
		logger.TraceLeave("Tools.GetToolTarget()")
		return -1, fmt.Errorf("unable to find tool %q", tool)
	}

	logger.TraceLeave("Tools.GetToolTarget()")
	return currentTemperatureData.Target, nil
}


func SetToolTarget(client *octoprintApis.Client, tool string, target float64) error {
	logger.TraceEnter("Tools.SetToolTarget()")

	if tool == "bed" {
		cmd := &octoprintApis.BedTargetRequest{Target: target}
		logger.TraceLeave("Tools.SetToolTarget()")
		return cmd.Do(client)
	}

	cmd := &octoprintApis.ToolTargetRequest{Targets: map[string]float64{tool: target}}
	logger.TraceLeave("Tools.SetToolTarget()")
	return cmd.Do(client)
}


func GetCurrentTemperatureData(client *octoprintApis.Client) (map[string]dataModels.TemperatureData, error) {
	logger.TraceEnter("Tools.GetCurrentTemperatureData()")

	temperatureDataResponse, err := (&octoprintApis.TemperatureDataRequest{}).Do(client)
	if err != nil {
		logger.LogError("tools.GetCurrentTemperatureData()", "Do(TemperatureDataRequest)", err)
		logger.TraceLeave("Tools.GetCurrentTemperatureData()")
		return nil, err
	}

	if temperatureDataResponse == nil {
		logger.Error("tools.GetCurrentTemperatureData() - temperatureDataResponse is nil")
		logger.TraceLeave("Tools.GetCurrentTemperatureData()")
		return nil, err
	}

	// Can't test for temperatureDataResponse.TemperatureStateResponse == nil (type mismatch)

	if temperatureDataResponse.TemperatureStateResponse.CurrentTemperatureData == nil {
		logger.Error("tools.GetCurrentTemperatureData() - temperatureDataResponse.TemperatureStateResponse.CurrentTemperatureData is nil")
		logger.TraceLeave("Tools.GetCurrentTemperatureData()")
		return nil, err
	}

	logger.TraceLeave("Tools.GetCurrentTemperatureData()")
	return temperatureDataResponse.TemperatureStateResponse.CurrentTemperatureData, nil
}


func CheckIfHotendTemperatureIsTooLow(client *octoprintApis.Client, extruderId, action string, parentWindow *gtk.Window) bool {
	logger.TraceEnter("Tools.CheckIfHotendTemperatureIsTooLow()")

	currentTemperatureData, err := GetCurrentTemperatureData(client)
	if err != nil {
		logger.LogError("tools.CurrentHotendTemperatureIsTooLow()", "GetCurrentTemperatureData()", err)
		logger.TraceLeave("Tools.CheckIfHotendTemperatureIsTooLow()")
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

		logger.TraceLeave("Tools.CheckIfHotendTemperatureIsTooLow()")
		return true
	}

	logger.TraceLeave("Tools.CheckIfHotendTemperatureIsTooLow()")
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


func GetExtruderFileName(hotendIndex, hotendCount int) string {
	strImageFileName := ""
	if hotendIndex == 1 && hotendCount == 1 {
		strImageFileName = "extruder-typeB.svg"
	} else {
		strImageFileName = fmt.Sprintf("extruder-typeB-%d.svg", hotendIndex)
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
	logger.Infof("tools.HotendTemperatureIsTooLow() - targetTemperature is %.2f", targetTemperature)

	actualTemperature := temperatureData.Actual
	logger.Infof("tools.HotendTemperatureIsTooLow() - actualTemperature is %.2f", actualTemperature)

	if targetTemperature <= MIN_HOTEND_TEMPERATURE || actualTemperature <= MIN_HOTEND_TEMPERATURE {
		return true
	}

	return false
}
