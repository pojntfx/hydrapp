//go:build !selfupdate
// +build !selfupdate

package update

import "context"

func Update(
	ctx context.Context,

	apiURL,
	owner,
	repo,

	currentVersion,
	appID string,

	state *BrowserState,
	handlePanic func(msg string, err error),
) error {
	return nil
}
