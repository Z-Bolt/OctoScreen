package dataModels


// PluginManagerInfoResponse -
type PluginManagerInfoResponse struct {
	Octoprint			string `json:"octoprint"`
	IsOnline			bool `json:"online"`
	//orphan_data: { }
	OS					string `json:"os"`
	//pip: {}
	Plugins				[]Plugin `json:"plugins"`
}
