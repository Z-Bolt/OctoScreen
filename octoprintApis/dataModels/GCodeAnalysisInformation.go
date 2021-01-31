package dataModels

// GCodeAnalysisInformation Information from the analysis of the GCODE file.
type GCodeAnalysisInformation struct {
	// EstimatedPrintTime is the estimated print time of the file, in seconds.
	EstimatedPrintTime float64 `json:"estimatedPrintTime"`

	// Filament estimated usage of filament
	Filament Filament `json:"filament"`
}
