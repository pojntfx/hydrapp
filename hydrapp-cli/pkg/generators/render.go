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
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}

		return ioutil.WriteFile(path, []byte(tpl), os.ModePerm)
	}

	t, err := template.New(path).Parse(tpl)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(t.Name()), os.ModePerm); err != nil {
		return err
	}

	dst, err := os.OpenFile(t.Name(), os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer dst.Close()

	return t.Execute(dst, data)
}
