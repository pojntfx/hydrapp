package generators

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

func RenderTemplate(path string, tpl string, data any) error {
	// Assume that templates without data are just files
	if data == nil {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		return ioutil.WriteFile(path, []byte(tpl), 0664)
	}

	t, err := template.New(path).Parse(tpl)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(t.Name()), 0755); err != nil {
		return err
	}

	dst, err := os.OpenFile(t.Name(), os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	defer dst.Close()

	return t.Execute(dst, data)
}
