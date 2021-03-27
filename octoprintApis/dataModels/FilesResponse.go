package dataModels


// FilesResponse is the response to a FilesRequest.
type FilesResponse struct {
	// Files is the list of requested files.  Might be an empty list if no files are available
	Files    []*FileResponse

	//
	Children []*FileResponse

	// Free is the amount of disk space in bytes available in the local disk
	// space (refers to OctoPrint’s `uploads` folder).  Only returned if file
	// list was requested for origin `local` or all origins.
	Free uint64
}
