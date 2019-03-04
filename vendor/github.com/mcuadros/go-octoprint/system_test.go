package octoprint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemCommandsRequest_Do(t *testing.T) {
	cli := NewClient("http://localhost:5000", "")

	r := &SystemCommandsRequest{}
	state, err := r.Do(cli)
	assert.NoError(t, err)

	assert.Len(t, state.Core, 1)
	assert.Len(t, state.Custom, 0)
	assert.Equal(t, "shutdown", state.Core[0].Action)
}

func TestSystemExecuteCommandRequest_Do(t *testing.T) {
	cli := NewClient("http://localhost:5000", "")

	r := &SystemExecuteCommandRequest{}
	err := r.Do(cli)
	assert.Error(t, err)

	r = &SystemExecuteCommandRequest{Source: Core, Action: "shutdown"}
	err = r.Do(cli)
	assert.NoError(t, err)
}
