package dataModels


type PrinterProfilesResponse struct {
	Profiles []*PrinterProfileResponse `json:"profiles"`
}
