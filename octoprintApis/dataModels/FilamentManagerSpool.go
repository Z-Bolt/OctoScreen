package dataModels

type FilamentManagerSpoolProfile struct {
	// Filament density
	Density float64 `json: "density"`

	// Filament diameter
	Diameter float64 `json: "diameter"`

	// Profile Id
	Id int `json: "id"`

	// Material name
	Material string `json: "material"`

	// Spool vendor
	Vendor string `json: "vendor"`
}

// ConnectionResponse is the response from a connection command.
type FilamentManagerSpool struct {
	// Cost of the spool
	Cost float64 `json:"cost"`

	// Spool ID
	Id int `json: id`

	// Name of the spool
	Name string `json: "name"`

	// Spool profile
	Profile FilamentManagerSpoolProfile `json: "profile"`

	// Temperature offset of spool
	TempOffset float64 `json: "temp_offset"`

	// Used filament
	Used float64 `json: "used"`

	// Starting weight of the spool
	Weight float64 `json: "Weight"`
}
