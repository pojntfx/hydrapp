package renderers

import (
	"os"
	"path/filepath"
)

func WriteRenders(workdir string, renders []Renderer, overwrite, ejecting bool) error {
	for _, renderer := range renders {
		if path, content, err := renderer.Render(""); err != nil {
			return err
		} else {
			if err := os.MkdirAll(filepath.Dir(filepath.Join(workdir, path)), 0755); err != nil {
				return err
			}

			file := filepath.Join(workdir, path)
			exists := true
			if _, err := os.Stat(file); err != nil {
				exists = false
			}

			if !exists || overwrite {
				if err := os.WriteFile(file, []byte(content), 0664); err != nil {
					return err
				}
			}

			if exists && !ejecting {
				// Read existing file and allow for templating in them
				t, err := os.ReadFile(file)
				if err != nil {
					return err
				}

				if _, content, err = renderer.Render(string(t)); err != nil {
					return err
				}

				if err := os.WriteFile(file, []byte(content), 0664); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
