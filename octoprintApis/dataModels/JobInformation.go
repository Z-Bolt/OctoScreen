package dataModels

import (
	// "encoding/json"
	// "strconv"
	// "strings"
	// "time"
)


// Job information
// https://docs.octoprint.org/en/master/api/datamodel.html#job-information


// JobInformation contains information regarding the target of the current job.
type JobInformation struct {
	// File is the file that is the target of the current print job.
	File FileResponse `json:"file"`

	// EstimatedPrintTime is the estimated print time for the file, in seconds.
	EstimatedPrintTime float64 `json:"estimatedPrintTime"`

	// LastPrintTime is the print time of the last print of the file, in seconds.
	LastPrintTime float64 `json:"lastPrintTime"`

	// Filament contains Information regarding the estimated filament usage of the print job.
	Filament struct {
		// Length of filament used, in mm
		Length float64 `json:"length"`

		// Volume of filament used, in cm³
		Volume float64 `json:"volume"`
	} `json:"filament"`
}
