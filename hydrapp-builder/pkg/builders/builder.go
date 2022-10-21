package builders

type Builder interface {
	Build() error
	Render(workdir string) error
}
