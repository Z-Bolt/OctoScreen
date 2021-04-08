package octoprintApis

import (
	"bytes"
	"encoding/json"
	// "fmt"
	// "io"
	// "strings"

	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


const PrinterPrintHeadApiUri     = "/api/printer/printhead"
const PrinterToolApiUri          = "/api/printer/tool"
const PrinterBedApiUri           = "/api/printer/bed"
const PrinterSdApiUri            = "/api/printer/sd"
const PrinterCommandApiUri       = "/api/printer/command"
const PrinterCommandCustomApiUri = "/api/printer/command/custom"


var (
	PrintErrors = StatusMapping {
		409: "Printer is not operational",
	}

	PrintHeadJobErrors = StatusMapping {
		400: "Invalid axis specified, invalid value for travel amount for a jog command or factor for feed rate or otherwise invalid request",
		409: "Printer is not operational or currently printing",
	}

	PrintToolErrors = StatusMapping {
		400: "Targets or offsets contains a property or tool contains a value not matching the format tool{n}, the target/offset temperature, extrusion amount or flow rate factor is not a valid number or outside of the supported range, or if the request is otherwise invalid",
		409: "Printer is not operational",
	}

	PrintBedErrors = StatusMapping {
		409: "Printer is not operational or the selected printer profile does not have a heated bed.",
	}

	PrintSdErrors = StatusMapping {
		404: "SD support has been disabled in OctoPrint’s settings.",
		409: "SD card has not been initialized.",
	}
)

// doCommandRequest can be used in any operation where the only required field is the `command` field.
func doCommandRequest(
	client *Client,
	uri string,
	command string,
	statusMapping StatusMapping,
) error {
	v := map[string]string{"command": command}

	buffer := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buffer).Encode(v); err != nil {
		return err
	}

	_, err := client.doJsonRequest("POST", uri, buffer, statusMapping, true)

	return err
}
