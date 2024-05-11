package xdg

import (
	"bytes"
	_ "embed"
	"os"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/utils"
)

type iconRenderer struct {
	inputFilePath  string
	outputFilePath string

	imageType utils.ImageType

	width  int
	height int
}

func NewIconRenderer(
	inputFilePath string,
	outputFilePath string,

	imageType utils.ImageType,

	width int,
	height int,
) renderers.Renderer {
	return &iconRenderer{
		inputFilePath:  inputFilePath,
		outputFilePath: outputFilePath,

		imageType: imageType,

		width:  width,
		height: height,
	}
}

func (r *iconRenderer) Render(templateOverride string) (filePath string, fileContent []byte, err error) {
	inputFile, err := os.Open(r.inputFilePath)
	if err != nil {
		return "", []byte{}, err
	}
	defer inputFile.Close()

	outputFile := &bytes.Buffer{}
	if err := utils.ConvertPNG(
		inputFile,
		outputFile,

		r.imageType,

		r.width,
		r.height,
	); err != nil {
		return "", []byte{}, err
	}

	return r.outputFilePath, outputFile.Bytes(), nil
}
