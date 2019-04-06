package octoprint

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHistoricTemperatureData(t *testing.T) {
	js := []byte(`{
		"time": 1395651928,
		"tool0": {
			"actual": 214.8821,
			"target": 220.0
		},
		"tool1": {
			"actual": 25.3,
			"target": null
		}
	}`)

	h := &HistoricTemperatureData{}

	err := json.Unmarshal(js, h)
	assert.NoError(t, err)

	assert.Len(t, h.Tools, 2)
	assert.False(t, h.Time.IsZero())
	assert.Equal(t, h.Tools["tool0"].Target, 220.)
	assert.Equal(t, h.Tools["tool1"].Actual, 25.3)
}

func TestTemperatureState(t *testing.T) {
	js := []byte(`{
		"tool0": {
			"actual": 214.8821,
			"target": 220.0,
			"offset": 0
		},
		"tool1": {
			"actual": 25.3,
			"target": null,
			"offset": 0
		},
		"history": [
			{
			"time": 1395651928,
			"tool0": {
				"actual": 214.8821,
				"target": 220.0
			},
			"tool1": {
				"actual": 25.3,
				"target": null
			}
			},
			{
			"time": 1395651926,
			"tool0": {
				"actual": 212.32,
				"target": 220.0
			},
			"tool1": {
				"actual": 25.1
			}
			}
		]
	}`)

	r := &TemperatureState{}

	err := json.Unmarshal(js, r)
	assert.NoError(t, err)

	assert.Len(t, r.Current, 2)
	assert.Equal(t, r.Current["tool0"].Actual, 214.8821)
	assert.Equal(t, r.Current["tool1"].Actual, 25.3)

	assert.Len(t, r.History, 2)
	assert.Equal(t, r.History[0].Tools["tool0"].Actual, 214.8821)
	assert.Equal(t, r.History[1].Tools["tool0"].Actual, 212.32)
}

func TestFullStateResponse(t *testing.T) {
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
				"history":[{
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
				},{
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
				}]
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

	r := &FullStateResponse{}

	err := json.Unmarshal(js, r)
	assert.NoError(t, err)

	assert.Equal(t, r.State.Text, "Operational")
	assert.True(t, r.State.Flags.Ready)
	assert.True(t, r.SD.Ready)
	assert.Len(t, r.Temperature.Current, 3)
	assert.Len(t, r.Temperature.History, 2)
}

func TestFileInformation_IsFolder(t *testing.T) {
	f := &FileInformation{TypePath: []string{"folder"}}
	assert.True(t, f.IsFolder())

	f = &FileInformation{}
	assert.False(t, f.IsFolder())
}

func TestJSONTime_UnmarshalJSONWithNull(t *testing.T) {
	time := &JSONTime{}
	err := time.UnmarshalJSON([]byte("null"))
	assert.NoError(t, err)
}
