package builders

const UnstableIDSuffix = ".unstable"
const UnstableNameSuffix = " (Unstable)"
const UnstablePathSuffix = "unstable"

const StablePathSuffix = "stable"

type Builder interface {
	Build() error
	Render(workdir string, ejecting bool) error
}
