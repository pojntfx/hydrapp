package utils

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

var (
	ErrCouldNotStartFileWatcher     = errors.New("could not start file watcher")
	ErrCouldNotAddPathToFileWatcher = errors.New("could not add path to file watcher")
	ErrCouldNotWatchFile            = errors.New("could not watch file")
)

func SetupFileWatcher(path string, watchDir bool) (waitForRemoval func() error, close func() error, err error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return func() error { return nil }, func() error { return nil }, errors.Join(ErrCouldNotStartFileWatcher, err)
	}
	defer watcher.Close()

	// Wait until browser has exited
	watchPath := path
	if watchDir {
		watchPath = filepath.Dir(path)
	}

	if err = watcher.Add(watchPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if watchDir {
				if err := os.MkdirAll(watchPath, os.ModePerm); err != nil {
					return func() error { return nil }, watcher.Close, err
				}
			} else {
				return func() error { return nil }, watcher.Close, nil
			}
		}

		return func() error { return nil }, watcher.Close, errors.Join(ErrCouldNotAddPathToFileWatcher, err)
	}

	return func() error {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return nil
				}

				if (!watchDir || event.Name == path) && event.Op&fsnotify.Remove == fsnotify.Remove {
					return nil
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return nil
				}

				return errors.Join(ErrCouldNotWatchFile, err)
			}
		}
	}, watcher.Close, nil
}
