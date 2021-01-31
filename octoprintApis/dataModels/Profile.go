package dataModels

// Profile describe a printer profile.
type Profile struct {
	// ID is the identifier of the profile.
	ID string `json:"id"`

	// Name is the display name of the profile.
	Name string `json:"name"`
}
