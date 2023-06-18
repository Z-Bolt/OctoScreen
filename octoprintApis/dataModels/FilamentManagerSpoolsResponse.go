package dataModels


type FilamentManagerSpoolsResponse struct {
	Spools []*FilamentManagerSpool `json: "spools"`
}
