package dataModels


// APIConfig REST API settings.
type APIConfig struct {
	// Enabled whether to enable the API.
	IsEnabled bool `json:"enabled"`

	// Key current API key needed for accessing the API
	Key string `json:"key"`
}
