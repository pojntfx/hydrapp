//go:build !selfupdate
// +build !selfupdate

package update

func Update(repo string, version string, state *BrowserState) error {
	return nil
}
