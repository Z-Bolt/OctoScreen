package octoprintApis

import (
	// "bytes"
	// "encoding/json"
	// "fmt"
	// "io"
	// "mime/multipart"

	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


const FilesApiUri = "/api/files"


var (
	FilesLocationGETErrors = StatusMapping {
		404: "Location is neither local nor sdcard",
	}

	FilesLocationPOSTErrors = StatusMapping {
		400: "No file or foldername are included in the request, userdata was provided but could not be parsed as JSON or the request is otherwise invalid.",
		404: "Location is neither local nor sdcard or trying to upload to SD card and SD card support is disabled",
		409: "The upload of the file would override the file that is currently being printed or if an upload to SD card was requested and the printer is either not operational or currently busy with a print job.",
		415: "The file is neither a gcode nor an stl file (or it is an stl file but slicing support is disabled)",
		500: "The upload failed internally",
	}

	FilesLocationPathPOSTErrors = StatusMapping {
		400: "The command is unknown or the request is otherwise invalid",
		415: "A slice command was issued against something other than an STL file.",
		404: "Location is neither local nor sdcard or the requested file was not found",
		409: "Selected file is supposed to start printing directly but the printer is not operational or if a file to be sliced is supposed to be selected or start printing directly but the printer is not operational or already printing.",
	}

	FilesLocationDeleteErrors = StatusMapping {
		404: "Location is neither local nor sdcard",
		409: "The file to be deleted is currently being printed",
	}
)
