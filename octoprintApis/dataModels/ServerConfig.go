package dataModels


// ServerConfig settings to configure the server.
type ServerConfig struct {
	// Commands to restart/shutdown octoprint or the system it's running on.
	Commands struct {
		// ServerRestartCommand to restart OctoPrint, defaults to being unset
		ServerRestartCommand string `json:"serverRestartCommand"`

		//SystemRestartCommand  to restart the system OctoPrint is running on,
		// defaults to being unset
		SystemRestartCommand string `json:"systemRestartCommand"`

		// SystemShutdownCommand Command to shut down the system OctoPrint is
		// running on, defaults to being unset
		SystemShutdownCommand string `json:"systemShutdownCommand"`
	} `json:"commands"`

	// Diskspace settings of when to display what disk space warning
	Diskspace struct {
		// Warning threshold (bytes) after which to consider disk space becoming
		// sparse, defaults to 500MB.
		Warning uint64 `json:"warning"`

		// Critical threshold (bytes) after which to consider disk space becoming
		// critical, defaults to 200MB.
		Critical uint64 `json:"critical"`
	} `json:"diskspace"`

	// OnlineCheck configuration of the regular online connectivity check.
	OnlineCheck struct {
		// Enabled whether the online check is enabled, defaults to false due to
		// valid privacy concerns.
		IsEnabled bool `json:"enabled"`

		// Interval in which to check for online connectivity (in seconds)
		Interval int `json:"interval"`

		// Host DNS host against which to check (default: 8.8.8.8 aka Google's DNS)
		Host string `json:"host"`

		// DNS port against which to check (default: 53 - the default DNS port)
		Port int `json:"port"`
	} `json:"onlineCheck"`

	// PluginBlacklist configuration of the plugin blacklist
	PluginBlacklist struct {
		// Enabled whether use of the blacklist is enabled, defaults to false
		IsEnabled bool `json:"enabled"`

		// URL from which to fetch the blacklist
		URL string `json:"url"`

		// TTL is time to live of the cached blacklist, in secs (default: 15mins)
		TTL int `json:"ttl"`
	} `json:"pluginBlacklist"`
}
