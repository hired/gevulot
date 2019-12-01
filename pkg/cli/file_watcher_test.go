package cli

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFileWatcherWatch(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "config")

	if err != nil {
		t.Fatal(err)
	}

	// Cleanup
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte("foo")); err != nil {
		t.Fatal(err)
	}

	fileChanged := false

	watcher := newFileWatcher(tmpfile.Name())
	watcher.OnCreate = func() {
		panic("unexpected OnCreate callback fired from a file watcher")
	}
	watcher.OnWrite = func() {
		fileChanged = true
	}
	watcher.OnRemove = func() {
		panic("unexpected OnRemove callback fired from a file watcher")
	}
	watcher.Watch()

	// If we don't stop the watcher OnRemove will trigger a panic
	defer watcher.StopWatch()

	assert.False(t, fileChanged)

	if _, err := tmpfile.Write([]byte("bar")); err != nil {
		t.Fatal(err)
	}

	assert.Eventually(t, func() bool { return fileChanged }, time.Second, 10*time.Millisecond)
}
