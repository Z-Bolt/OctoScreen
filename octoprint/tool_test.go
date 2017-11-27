package octoprint

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToolResponse(t *testing.T) {
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

	r := &ToolResponse{}

	err := json.Unmarshal(js, r)
	assert.NoError(t, err)

	assert.Len(t, r.Current, 2)
	assert.Equal(t, r.Current["tool0"].Actual, 214.8821)
	assert.Equal(t, r.Current["tool1"].Actual, 25.3)

	assert.Len(t, r.History, 2)
	assert.Equal(t, r.History[0].Tools["tool0"].Actual, 214.8821)
	assert.Equal(t, r.History[1].Tools["tool0"].Actual, 212.32)
}
