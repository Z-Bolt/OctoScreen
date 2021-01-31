package dataModels


// PrintStats information from the print stats of a file.
type PrintStats struct {
	// Failure number of failed prints.
	Failure int `json:"failure"`

	// Success number of success prints.
	SuccessfullPrintCount int `json:"success"`

	// Last print information.
	Last struct {
		// Date of the last print.
		Date JsonTime `json:"date"`

		// Success or not.
		IsSuccess bool `json:"success"`
	} `json:"last"`
}
