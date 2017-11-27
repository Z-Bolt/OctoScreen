package octoprint

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHistory(t *testing.T) {
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

	h := &History{}

	err := json.Unmarshal(js, h)
	assert.NoError(t, err)

	assert.Len(t, h.Tools, 2)
	assert.False(t, h.Time.IsZero())
	assert.Equal(t, h.Tools["tool0"].Target, 220.)
	assert.Equal(t, h.Tools["tool1"].Actual, 25.3)
}
