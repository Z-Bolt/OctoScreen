package dataModels


// WebcamConfig settings to configure webcam support.
type WebcamConfig struct {
	// StreamUrl use this option to enable display of a webcam stream in the
	// UI, e.g. via MJPG-Streamer. Webcam support will be disabled if not
	// set. Maps to webcam.stream in config.yaml.
	StreamURL string `json:"streamUrl"`

	// SnapshotURL use this option to enable timelapse support via snapshot,
	// e.g. via MJPG-Streamer. Timelapse support will be disabled if not set.
	// Maps to webcam.snapshot in config.yaml.
	SnapshotURL string `json:"snapshotUrl"`

	// FFmpegPath path to ffmpeg binary to use for creating timelapse
	// recordings. Timelapse support will be disabled if not set. Maps to
	// webcam.ffmpeg in config.yaml.
	FFmpegPath string `json:"ffmpegPath"`

	// Bitrate to use for rendering the timelapse video. This gets directly
	// passed to ffmpeg.
	Bitrate int `json:"bitrate"`

	// FFmpegThreads number of how many threads to instruct ffmpeg to use for
	// encoding. Defaults to 1. Should be left at 1 for RPi1.
	FFmpegThreads int `json:"ffmpegThreads"`

	// Watermark whether to include a "created with OctoPrint" watermark in the
	// generated timelapse movies.
	Watermark string `json:"watermark"`

	// FlipH whether to flip the webcam horizontally.
	FlipH bool `json:"flipH"`

	// FlipV whether to flip the webcam vertically.
	FlipV bool `json:"flipV"`

	// Rotate90 whether to rotate the webcam 90° counter clockwise.
	Rotate90 bool `json:"rotate90"`
}