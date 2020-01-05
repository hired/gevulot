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

	fileChanged := make(chan bool)

	watcher := newFileWatcher(tmpfile.Name())
	watcher.OnCreate = func() {
		panic("unexpected OnCreate callback fired from a file watcher")
	}
	watcher.OnWrite = func() {
		fileChanged <- true
	}
	watcher.OnRemove = func() {
		panic("unexpected OnRemove callback fired from a file watcher")
	}

	err = watcher.Watch()
	assert.NoError(t, err)

	// If we don't stop the watcher OnRemove will trigger a panic
	defer watcher.StopWatch()

	select {
	case <-fileChanged:
		assert.FailNow(t, "unexpected OnWrite callback")
	default:
	}

	if _, err := tmpfile.Write([]byte("bar")); err != nil {
		t.Fatal(err)
	}

	timer := time.NewTimer(5 * time.Second)

	select {
	case <-fileChanged:
		timer.Stop()

	case <-timer.C:
		assert.FailNow(t, "OnWrite hasn't been fired after file change")
	}
}
