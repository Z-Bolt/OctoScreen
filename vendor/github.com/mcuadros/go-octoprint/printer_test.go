package octoprint

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStateRequest_Do(t *testing.T) {
	cli := NewClient("http://localhost:5000", "")

	r := &StateRequest{}
	state, err := r.Do(cli)
	assert.NoError(t, err)

	assert.Equal(t, "Operational", state.State.Text)
	assert.Len(t, state.Temperature.Current, 2)
	assert.Len(t, state.Temperature.History, 0)
}

func TestStateRequest_DoWithHistory(t *testing.T) {
	cli := NewClient("http://localhost:5000", "")

	r := &StateRequest{
		IncludeHistory: true,
	}
	state, err := r.Do(cli)
	assert.NoError(t, err)

	assert.Equal(t, "Operational", state.State.Text)
	assert.Len(t, state.Temperature.Current, 2)
	assert.True(t, len(state.Temperature.History) > 0)
}

func TestStateRequest_DoWithExclude(t *testing.T) {
	cli := NewClient("http://localhost:5000", "")

	r := &StateRequest{Exclude: []string{"temperature"}}
	state, err := r.Do(cli)
	assert.NoError(t, err)

	assert.Equal(t, "Operational", state.State.Text)
	assert.Len(t, state.Temperature.Current, 0)
}

func TestSDInitRequest_Do(t *testing.T) {
	cli := NewClient("http://localhost:5000", "")

	r := &SDInitRequest{}
	err := r.Do(cli)
	assert.NoError(t, err)

	time.Sleep(50 * time.Millisecond)

	state, err := (&SDStateRequest{}).Do(cli)
	assert.NoError(t, err)
	assert.True(t, state.Ready)
}

func TestSDReleaseRequest_Do(t *testing.T) {
	cli := NewClient("http://localhost:5000", "")

	r := &SDReleaseRequest{}
	err := r.Do(cli)
	assert.NoError(t, err)

	state, err := (&SDStateRequest{}).Do(cli)
	assert.NoError(t, err)
	assert.False(t, state.Ready)
}

func TestSDRefreshRequest_Do(t *testing.T) {
	cli := NewClient("http://localhost:5000", "")

	r := &SDRefreshRequest{}
	err := r.Do(cli)
	assert.NoError(t, err)

	state, err := (&SDStateRequest{}).Do(cli)
	assert.NoError(t, err)
	assert.False(t, state.Ready)
}

func TestCustomCommandsRequest_Do(t *testing.T) {
	cli := NewClient("http://localhost:5000", "")

	r := &CustomCommandsRequest{}
	s, err := r.Do(cli)
	assert.NoError(t, err)

	assert.Len(t, s.Controls, 1)
	assert.Equal(t, s.Controls[0].Name, "Example for multiple commands")
	assert.Len(t, s.Controls[0].Children, 1)
	assert.Equal(t, s.Controls[0].Children[0].Name, "Move X (static)")
	assert.Len(t, s.Controls[0].Children[0].Commands, 3)
}
