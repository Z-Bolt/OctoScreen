package octoprint

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateResponse(t *testing.T) {
	js := []byte(`
		{
			"temperature":{
			   "tool0":{
				  "actual":214.8821,
				  "target":220.0,
				  "offset":0
			   },
			   "tool1":{
				  "actual":25.3,
				  "target":null,
				  "offset":0
			   },
			   "bed":{
				  "actual":50.221,
				  "target":70.0,
				  "offset":5
			   },
			   "history":[
				  {
					 "time":1395651928,
					 "tool0":{
						"actual":214.8821,
						"target":220.0
					 },
					 "tool1":{
						"actual":25.3,
						"target":null
					 },
					 "bed":{
						"actual":50.221,
						"target":70.0
					 }
				  },
				  {
					 "time":1395651926,
					 "tool0":{
						"actual":212.32,
						"target":220.0
					 },
					 "tool1":{
						"actual":25.1,
						"target":null
					 },
					 "bed":{
						"actual":49.1123,
						"target":70.0
					 }
				  }
			   ]
			},
			"sd":{
			   "ready":true
			},
			"state":{
			   "text":"Operational",
			   "flags":{
				  "operational":true,
				  "paused":false,
				  "printing":false,
				  "sdReady":true,
				  "error":false,
				  "ready":true,
				  "closedOrError":false
			   }
			}
		 }
	`)

	r := &StateResponse{}

	err := json.Unmarshal(js, r)
	assert.NoError(t, err)

	assert.Equal(t, r.State.Text, "Operational")
	assert.True(t, r.State.Flags.Ready)
	assert.True(t, r.SD.Ready)
	assert.Len(t, r.Temperature.Current, 3)
	assert.Len(t, r.Temperature.History, 2)
}
