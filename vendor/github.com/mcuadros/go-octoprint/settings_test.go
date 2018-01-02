package octoprint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSettingsRequest_Do(t *testing.T) {
	cli := NewClient("http://localhost:5000", "")

	r := &SettingsRequest{}
	settings, err := r.Do(cli)
	assert.NoError(t, err)

	assert.Equal(t, settings.API.Enabled, true)
}
