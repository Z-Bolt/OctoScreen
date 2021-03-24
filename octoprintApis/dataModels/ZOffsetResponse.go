package dataModels


type ZOffsetResponse struct {
	// Job contains information regarding the target of the current print job.
	Offset float64 `json:"offset"`
}
