package utils

import (
	"os"

	"github.com/fsnotify/fsnotify"
)

func WaitForFileRemoval(path string, handlePanic func(msg string, err error)) {
	if _, err := os.Stat(path); err == nil {
		// Wait until browser has exited
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			handlePanic("could not start lockfile watcher", err)
		}
		defer watcher.Close()

		done := make(chan struct{})
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}

					// Stop the app
					if event.Op&fsnotify.Remove == fsnotify.Remove {
						done <- struct{}{}

						return
					}

				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}

					handlePanic("could not continue watching lockfile", err)
				}
			}
		}()

		err = watcher.Add(path)
		if err != nil {
			handlePanic("could not watch lockfile", err)
		}

		<-done
	}
}
