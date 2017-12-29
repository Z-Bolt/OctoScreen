package octoprint

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadFileRequest_Do(t *testing.T) {
	cli := NewClient("http://localhost:5000", "")

	r := &UploadFileRequest{Location: Local}
	err := r.AddFile("foo.gcode", bytes.NewBufferString("foo"))
	assert.NoError(t, err)

	state, err := r.Do(cli)
	assert.NoError(t, err)
	assert.Equal(t, "foo.gcode", state.File.Local.Name)

	err = (&DeleteFileRequest{Location: Local, Path: "foo.gcode"}).Do(cli)
	assert.NoError(t, err)
}

func TestUploadFileRequest_DoWithFolder(t *testing.T) {
	cli := NewClient("http://localhost:5000", "")

	r := &UploadFileRequest{Location: Local}
	err := r.AddFolder("qux")
	assert.NoError(t, err)

	state, err := r.Do(cli)
	assert.NoError(t, err)
	assert.Equal(t, true, state.Done)
}

func TestFilesRequest_Do(t *testing.T) {
	cli := NewClient("http://localhost:5000", "")

	ur := &UploadFileRequest{Location: Local}
	err := ur.AddFile("foo.gcode", bytes.NewBufferString("foo"))
	assert.NoError(t, err)

	_, err = ur.Do(cli)
	assert.NoError(t, err)

	files, err := (&FilesRequest{}).Do(cli)
	assert.NoError(t, err)

	assert.True(t, len(files.Files) >= 1)
	err = (&DeleteFileRequest{Location: Local, Path: "foo.gcode"}).Do(cli)
	assert.NoError(t, err)
	return

	r := &FileRequest{Location: Local, Filename: "foo.gcode"}
	file, err := r.Do(cli)
	assert.NoError(t, err)

	assert.Equal(t, "foo.gcode", file.Name)

	err = (&DeleteFileRequest{Location: Local, Path: "foo.gcode"}).Do(cli)
	assert.NoError(t, err)
}

func TestSelectFileRequest_Do(t *testing.T) {
	cli := NewClient("http://localhost:5000", "")

	ur := &UploadFileRequest{Location: Local}
	err := ur.AddFile("foo2.gcode", bytes.NewBufferString("foo"))
	assert.NoError(t, err)
	_, err = ur.Do(cli)
	assert.NoError(t, err)

	r := &SelectFileRequest{Location: Local, Path: "foo2.gcode"}
	err = r.Do(cli)
	assert.NoError(t, err)
}

func xxxTestFilesRequest_DoWithLocation(t *testing.T) {
	cli := NewClient("http://localhost:5000", "")

	files, err := (&FilesRequest{Location: SDCard}).Do(cli)
	assert.NoError(t, err)

	assert.Len(t, files.Files, 0)
}
