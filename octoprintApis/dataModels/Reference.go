package dataModels


// Reference of a file.
type Reference struct {
	// Resource that represents the file or folder (e.g. for issuing commands
	// to or for deleting)
	Resource string `json:"resource"`

	// Download URL for the file. Never present for folders.
	Download string `json:"download"`

	// Model from which this file was generated (e.g. an STL, currently not
	// used). Never present for folders.
	Model string `json:"model"`
}
