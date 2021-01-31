package dataModels


// UploadFileResponse is the response to a UploadFileRequest.
type UploadFileResponse struct {
	// Abridged information regarding the file that was just uploaded. If only
	// uploaded to local this will only contain the local property. If uploaded
	// to SD card, this will contain both local and sdcard properties. Only
	// contained if a file was uploaded, not present if only a new folder was
	// created.
	File struct {
		// Local is the information regarding the file that was just uploaded
		// to the local storage.
		Local *FileResponse `json:"local"`

		// SDCard is the information regarding the file that was just uploaded
		// to the printer’s SD card.
		SDCard *FileResponse `json:"sdcard"`
	} `json:"files"`

	// Done whether any file processing after upload has already finished or
	// not, e.g. due to first needing to perform a slicing step. Clients may
	// use this information to direct progress displays related to the upload.
	IsDone bool `json:"done"`
}
