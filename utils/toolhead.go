package utils

import (
	"strings"
	"strconv"
	"log"

	"github.com/mcuadros/go-octoprint"
)

var cachedToolheadCount = -1

func GetToolheadCount(client *octoprint.Client) int {
	if cachedToolheadCount != -1 {
		return cachedToolheadCount
	}

	c, err := (&octoprint.ConnectionRequest{}).Do(client)
	if err != nil {
		log.Println(err)
		return 0
	}

	profile, err := (&octoprint.PrinterProfilesRequest{Id: c.Current.PrinterProfile}).Do(client)
	if err != nil {
		log.Println(err)
		return 0
	}

	if profile.Extruder.SharedNozzle {
		cachedToolheadCount = 1
	}

	cachedToolheadCount = profile.Extruder.Count


	// TODO: uncomment to force all toolheads to display and use for testing
	// cachedToolheadCount = 2


	return cachedToolheadCount
}


func GetDisplayNameForTool(toolName string) string {
	// Since this is such a hack, lets add some bounds checking
	if toolName == "" {
		log.Println("ERROR: GetDisplayNameForTool() - toolName is empty")
		return ""
	}

	lowerCaseName := strings.ToLower(toolName)
	if strings.LastIndex(lowerCaseName, "tool") != 0 {
		log.Println("ERROR: GetDisplayNameForTool() - toolName is invalid, value passed in was: " + toolName)
		return ""
	}

	if len(toolName) != 5 {
		log.Println("ERROR: GetDisplayNameForTool() - toolName is invalid, value passed in was: " + toolName)
		return ""
	}

	toolIndexAsInt, _ := strconv.Atoi(string(toolName[4]))
	displayName := toolName[0:4]
	displayName = displayName + strconv.Itoa(toolIndexAsInt + 1)

	return displayName
}
