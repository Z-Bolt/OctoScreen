package dataModels


type PrinterProfileResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	Volume struct {
		FormFactor string  `json:"formFactor"`
		Origin     string  `json:"origin"`
		Width      float64 `json:"width"`
		Depth      float64 `json:"depth"`
		Height     float64 `json:"height"`
	} `json:"volume"`

	Extruder struct {
		Count           int  `json:"count"`
		HasSharedNozzle bool `json:"sharedNozzle"`
	} `json:"extruder"`
}
