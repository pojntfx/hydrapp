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
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.Join(ErrCouldNotStartFileWatcher, err)
	}
	defer watcher.Close()

	// Wait until browser has exited
	if err = watcher.Add(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

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
