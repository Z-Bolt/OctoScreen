package octoprint

import (
	// "encoding/json"
	// "strconv"
	// "strings"
	// "time"
)

// Progress information
// https://docs.octoprint.org/en/master/api/datamodel.html#progress-information

// ProgressInformation contains information regarding the progress of the current print job.
type ProgressInformation struct {
	// Completion percentage of completion of the current print job.
	Completion float64 `json:"completion"`

	// FilePosition current position in the file being printed, in bytes from the beginning.
	FilePosition uint64 `json:"filepos"`

	// PrintTime is time already spent printing, in seconds.
	PrintTime float64 `json:"printTime"`

	// PrintTimeLeft is estimate of time left to print, in seconds.
	PrintTimeLeft float64 `json:"printTimeLeft"`

	// Origin of the current time left estimate.
	PrintTimeLeftOrigin string `json:"printTimeLeftOrigin"`
}
