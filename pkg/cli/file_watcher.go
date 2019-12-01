package cli

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// fileWatcher watches given file path, delivering FS events via the callbacks.
// Example:
//
//   watcher, err := newFileWatcher("/path")
//   watcher.OnWrite = func() {
//     // file has changed
//   }
//   watcher.Watch()
//
type fileWatcher struct {
	Path     string // Watched path.
	OnCreate func() // Called when the file is created on the disk.
	OnWrite  func() // Called when the file is written to the disk.
	OnRemove func() // Called when the file is removed from the the disk.

	watchDone chan bool // Channel to notify watcher to stop watching.
}

// newFileWatcher creates new fileWatcher.
func newFileWatcher(path string) *fileWatcher {
	return &fileWatcher{Path: path}
}

// Watch starts watching for the file changes on the disk.
func (w *fileWatcher) Watch() error {
	// We have 2 nested goroutines in this function.
	// We use watcherError variable to pull any errors from them into the calling gorouting.
	var watcherError error = nil

	// Initialization sync group. We have to wait until the watcher starts the watching to
	// return any possible errors.
	initWG := sync.WaitGroup{}
	initWG.Add(1)

	// This goroutine starts and closes fsnotify watcher
	go func() {
		watcher, err := fsnotify.NewWatcher()

		if err != nil {
			watcherError = err
			initWG.Done()
			return
		}

		defer watcher.Close()

		// Event loop sync group
		eventsWG := sync.WaitGroup{}
		eventsWG.Add(1)

		// Channel to notify handler to stop
		w.watchDone = make(chan bool)

		// FS events handler
		go func() {
			// Continue execution of the parent goroutine
			defer eventsWG.Done()

			for {
				select {
				// fsnotify watcher event
				case event, ok := <-watcher.Events:
					// 'Events' channel is closed
					if !ok {
						return
					}

					// Run the callback
					switch {
					case event.Op&fsnotify.Create != 0:
						if w.OnCreate != nil {
							w.OnCreate()
						}

					case event.Op&fsnotify.Write != 0:
						if w.OnWrite != nil {
							w.OnWrite()
						}

					case event.Op&fsnotify.Remove != 0:
						if w.OnRemove != nil {
							w.OnRemove()
						}
					}

				// fsnotify watcher error
				case err, ok := <-watcher.Errors:
					// 'Errors' channel is open
					if ok {
						watcherError = err
					}

					return

				// StopWatch has been called
				case <-w.watchDone:
					return
				}
			}
		}()

		// Start the watcher
		watcher.Add(w.Path)

		// Exit from Watch
		initWG.Done()

		// Event loop
		eventsWG.Wait()

		// We are not in the context of calling goroutine anymore so we cannot just return error
		// and have to print the error by ourselves.
		if watcherError != nil {
			fmt.Printf("config watcher error: %v\n", err)
		}
	}()

	// Wait for watcher to initialize
	initWG.Wait()

	// Return initialization error if any
	return watcherError
}

// StopWatch stops the file watcher.
func (w *fileWatcher) StopWatch() {
	close(w.watchDone)
}
