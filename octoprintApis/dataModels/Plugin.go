package dataModels


// Plugin -
type Plugin struct {
	Author                     string   `json:"author"`
	IsBlacklisted              bool     `json:"blacklisted"`
	Bundled                    bool     `json:"bundled"`
	Description                string   `json:"description"`
	// DisablingDiscouraged    bool     `json:"disabling_discouraged"`
	IsEnabled                  bool     `json:"enabled"`
	ForcedDisabled             bool     `json:"forced_disabled"`
	Incompatible               bool     `json:"incompatible"`
	Key                        string   `json:"key"`
	License                    string   `json:"license"`
	IsManagable                bool     `json:"managable"`
	Name                       string   `json:"name"`
	// notifications: []
	Origin                     string   `json:"origin"`
	PendingDisable             bool     `json:"pending_disable"`
	PendingEnable              bool     `json:"pending_enable"`
	PendingInstall             bool     `json:"pending_install"`
	PendingUninstall           bool     `json:"pending_uninstall"`
	Python                     string   `json:"python"`
	SafeModeVictim             bool     `json:"safe_mode_victim"`
	URL                        string   `json:"url"`
	Version                    string   `json:"version"`
}
