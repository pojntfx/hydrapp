package ui

// This needs to be `var`, not `const`, else Go can't overwrite them with
// `-ldflags="-X github.com/pojntfx/hydrapp/hydrapp/pkg/ui.SelfUpdaterBranchTimestampRFC3339=2024-10-09T23:14:53-07:00"` etc.
var (
	SelfUpdaterBranchTimestampRFC3339 = "" // Set by compiler
	SelfUpdaterBranchID               = "" // Set by compiler
	SelfUpdaterPackageType            = "" // Set by compiler
)
