package octoprintApis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// TODO: OctoScreenSettingsRequest seems like it's practically the same as PluginManagerInfoRequest
// Need to clean up and consolidate, or add comments as to why the two different classes.

type OctoScreenSettingsRequest struct {
	Command string `json:"command"`
}

func (this *OctoScreenSettingsRequest) Do(client *Client, uiState string) (*dataModels.OctoScreenSettingsResponse, error) {
	// Pause for 2 seconds here.  This is in response to the "OctoScreen 2.7.0 doesn't recognize that OctoScreen plugin 2.6.0 is installed"
	// bug (see https://github.com/Z-Bolt/OctoScreen/issues/275).  I was able to repro the bug sometimes, yet other times this worked fine.
	// I examined the logs and they were exactly the same (up to the error, which was "HTTP/1.x transport connection broken: malformed MIME
	// header line: Not found").  This might not be an issue after the state machine rewrite in 2.8.
	time.Sleep(time.Second * 2)

	target := fmt.Sprintf("%s?command=get_settings", PluginZBoltOctoScreenApiUri)
	bytes, err := client.doJsonRequest("GET", target, nil, ConnectionErrors)
	if err != nil {
		return nil, err
	}

	response := &dataModels.OctoScreenSettingsResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, err
}
