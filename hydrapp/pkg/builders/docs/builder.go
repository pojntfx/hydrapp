package docs

import (
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/builders"
	cconfig "github.com/pojntfx/hydrapp/hydrapp/pkg/config"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/executors"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers/docs"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/utils"
	uutils "github.com/pojntfx/hydrapp/hydrapp/pkg/utils"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-docs"
)

var (
	ErrInvalidRPMDistro = errors.New("invalid RPM distro")
)

func NewBuilder(
	ctx context.Context,
	cli *client.Client,

	image string, // OCI image to use
	pull bool, // Whether to pull the image or not
	src, // Input directory
	dst string, // Output directory
	onID func(id string), // Callback to handle container ID
	stdout io.Writer, // Writer to handle container output
	branchID, // Branch ID
	branchName, // Branch Name
	goMain string, // Directory with the main package to build
	cfg *cconfig.Root, // hydrapp config file
	overwrite bool, // Overwrite files even if they exist
) *Builder {
	return &Builder{
		ctx,
		cli,

		image,
		pull,
		src,
		dst,
		onID,
		stdout,
		branchID,
		branchName,
		goMain,
		cfg,
		overwrite,
	}
}

type Builder struct {
	ctx context.Context
	cli *client.Client

	image string
	pull  bool
	src,
	dst string
	onID   func(id string)
	stdout io.Writer
	branchID,
	branchName,
	goMain string
	cfg       *cconfig.Root
	overwrite bool
}

func (b *Builder) Render(workdir string, ejecting bool) error {
	appID := builders.GetAppIDForBranch(b.cfg.App.ID, b.branchID)
	appName := builders.GetAppNameForBranch(b.cfg.App.Name, b.branchName)
	macOSBinaryName := builders.GetAppIDForBranch(b.cfg.App.ID, b.branchID) + ".darwin.dmg"

	flatpaks := []docs.Artifact{}
	for _, f := range b.cfg.Flatpak {
		flatpaks = append(flatpaks, docs.Artifact{
			Architecture: f.Architecture,
			URL:          b.cfg.App.BaseURL + builders.GetPathForBranch(f.Path, b.branchID, "/") + "/hydrapp.flatpakrepo",
		})
	}

	msis := []docs.Artifact{}
	for _, m := range b.cfg.MSI {
		msis = append(msis, docs.Artifact{
			Architecture: m.Architecture,
			URL:          b.cfg.App.BaseURL + builders.GetPathForBranch(m.Path, b.branchID, "/") + "/" + builders.GetAppIDForBranch(b.cfg.App.ID, b.branchID) + ".windows-" + uutils.GetArchIdentifier(m.Architecture) + ".msi",
		})
	}

	rpms := []docs.DistroArtifact{}
	for _, r := range b.cfg.RPM {
		parts := strings.Split(r.Distro, "-")
		if len(parts) < 2 {
			return ErrInvalidRPMDistro
		}

		rpms = append(rpms, docs.DistroArtifact{
			Artifact: docs.Artifact{
				Architecture: r.Architecture,
				URL:          b.cfg.App.BaseURL + builders.GetPathForBranch(r.Path, b.branchID, "/") + "/repodata/hydrapp.repo",
			},
			DistroName:    parts[0],
			DistroVersion: parts[1],
		})
	}

	debs := []docs.DistroArtifact{}
	for _, d := range b.cfg.DEB {
		debs = append(debs, docs.DistroArtifact{
			Artifact: docs.Artifact{
				Architecture: d.Architecture,
				URL:          b.cfg.App.BaseURL + builders.GetPathForBranch(d.Path, b.branchID, "/"),
			},
			DistroName:    d.OS,
			DistroVersion: d.Distro,
		})
	}

	return utils.WriteRenders(
		filepath.Join(workdir, b.goMain),
		[]renderers.Renderer{
			docs.NewInstallationRenderer(
				appID,
				appName,
				b.cfg.App.BaseURL+builders.GetPathForBranch(b.cfg.APK.Path, b.branchID, "/")+"/repo", // F-Droid requires the path to end with `/repo`: `CRITICAL: repo_url needs to end with /repo`
				b.cfg.App.BaseURL+builders.GetPathForBranch(b.cfg.DMG.Path, b.branchID, "/")+"/"+macOSBinaryName,
				macOSBinaryName,
				b.cfg.App.BaseURL+builders.GetPathForBranch(b.cfg.Binaries.Path, b.branchID, "/"),
				flatpaks,
				msis,
				rpms,
				debs,
				strings.TrimSpace(b.cfg.DMG.Path) != "",
				strings.TrimSpace(b.cfg.APK.Path) != "",
				strings.TrimSpace(b.cfg.Binaries.Path) != "",
			),
		},
		b.overwrite,
		ejecting,
	)
}

func (b *Builder) Build() error {
	dst := builders.GetFilepathForBranch(b.dst, b.branchID)

	return executors.DockerRunImage(
		b.ctx,
		b.cli,
		b.image,
		b.pull,
		true,
		b.src,
		dst,
		b.onID,
		b.stdout,
		map[string]string{
			"GOMAIN": b.goMain,
		},
		b.Render,
		[]string{},
	)
}
