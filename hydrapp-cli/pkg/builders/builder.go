package builders

import (
	"path/filepath"
	"strings"
)

func GetAppIDForBranch(appID, branchID string) string {
	// Stable
	if strings.TrimSpace(branchID) == "" {
		return appID
	}

	return appID + "." + branchID
}

func GetAppNameForBranch(appName, branchName string) string {
	// Stable
	if strings.TrimSpace(branchName) == "" {
		return appName
	}

	return appName + " (" + branchName + ")"
}

func GetPathForBranch(path, branchID string) string {
	// Stable
	if strings.TrimSpace(branchID) == "" {
		return "/" + path + "/stable"
	}

	return "/" + path + "/" + branchID
}

func GetFilepathForBranch(path, branchID string) string {
	// Stable
	if strings.TrimSpace(branchID) == "" {
		return filepath.Join(path, "stable")
	}

	return filepath.Join(path, branchID)
}

type Builder interface {
	Build() error
	Render(workdir string, ejecting bool) error
}
