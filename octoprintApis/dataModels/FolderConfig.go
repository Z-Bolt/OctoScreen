package dataModels


// FolderConfig settings to set custom paths for folders used by OctoPrint.
type FolderConfig struct {
	// Uploads absolute path where to store gcode uploads. Defaults to the
	// uploads folder in the OctoPrint settings folder.
	Uploads string `json:"uploads"`

	// Timelapse absolute path where to store finished timelapse recordings.
	// Defaults to the timelapse folder in the OctoPrint settings dir.
	Timelapse string `json:"timelapse"`

	// TimelapseTmp absolute path where to store temporary timelapse files.
	// Defaults to the timelapse/tmp folder in the OctoPrint settings dir Maps
	// to folder.timelapse_tmp in config.yaml.
	TimelapseTmp string `json:"timelapseTmp"`

	// Logs absolute path where to store log files. Defaults to the logs folder
	// in the OctoPrint settings dir
	Logs string `json:"logs"`

	// Watched absolute path to a folder being watched for new files which then
	// get automatically added to OctoPrint (and deleted from that folder).
	// Can e.g. be used to define a folder which can then be mounted from remote
	// machines and used as local folder for quickly adding downloaded and/or
	// sliced objects to print in the future.
	Watched string `json:"watched"`
}
