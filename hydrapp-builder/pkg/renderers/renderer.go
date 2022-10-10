package renderers

type Renderer interface {
	Render() (filePath string, fileContent string, err error)
}
