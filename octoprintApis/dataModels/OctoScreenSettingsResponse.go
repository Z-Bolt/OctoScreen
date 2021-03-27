package dataModels


type OctoScreenSettingsResponse struct {
	// Job contains information regarding the target of the current print job.
	FilamentInLength  float64    `json:"filament_in_length"`
	FilamentOutLength float64    `json:"filament_out_length"`
	ToolChanger       bool       `json:"toolchanger"`
	XAxisInverted     bool       `json:"x_axis_inverted"`
	YAxisInverted     bool       `json:"y_axis_inverted"`
	ZAxisInverted     bool       `json:"z_axis_inverted"`
	MenuStructure     []MenuItem `json:"menu_structure"`
	GCodes            struct {
		AutoBedLevel string `json:"auto_bed_level"`
	} `json:"gcodes"`
}
