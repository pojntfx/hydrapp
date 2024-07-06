package utils

import (
	"errors"
	"os"

	"github.com/fsnotify/fsnotify"
)

var (
	ErrCouldNotStartFileWatcher     = errors.New("could not start file watcher")
	ErrCouldNotAddPathToFileWatcher = errors.New("could not add path to file watcher")
	ErrCouldNotWatchFile            = errors.New("could not watch file")
)

func WaitForFileRemoval(path string) error {
	if _, err := os.Stat(path); err == nil {
		// Wait until browser has exited
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return errors.Join(ErrCouldNotStartFileWatcher, err)
		}
		defer watcher.Close()

		if err = watcher.Add(path); err != nil {
			return errors.Join(ErrCouldNotAddPathToFileWatcher, err)
		}

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return nil
				}

				if event.Op&fsnotify.Remove == fsnotify.Remove {
					return nil
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return nil
				}

				return errors.Join(ErrCouldNotWatchFile, err)
			}
		}
	}

	return nil
}
