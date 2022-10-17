package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

func WriteRenders(workdir string, renders []*renderers.Renderer, overwrite bool) error {
	for _, renderer := range renders {
		if path, content, err := renderer.Render(); err != nil {
			return err
		} else {
			if err := os.MkdirAll(filepath.Dir(filepath.Join(workdir, path)), os.ModePerm); err != nil {
				return err
			}

			file := filepath.Join(workdir, path)
			exists := true
			if _, err := os.Stat(file); err != nil {
				exists = false
			}

			if !exists || overwrite {
				if err := ioutil.WriteFile(file, []byte(content), os.ModePerm); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
