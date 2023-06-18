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
