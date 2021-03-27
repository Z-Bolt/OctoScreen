package dataModels



// VersionResponse is the response from a job command.
type VersionResponse struct {
	// API is the API version.
	API string `json:"api"`

	// Server is the server version.
	Server string `json:"server"`
}
