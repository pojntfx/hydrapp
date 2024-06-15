package ui

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/ncruces/zenity"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/utils"
)

func HandlePanic(appName, msg string, err error) {
	// Create user-friendly error message
	body := fmt.Sprintf(`%v has encountered a fatal error and can't continue. The error message is:

%v

The following information might help you in fixing the problem:

%v

Strack trace:

%v`,
		appName,
		utils.Capitalize(msg),
		utils.Capitalize(err.Error()),
		string(debug.Stack()),
	)

	// Show error message visually using a dialog
	if err := zenity.Error(
		body,
		zenity.Title(fmt.Sprintf("Fatal error for %v", appName)),
		zenity.Width(320),
	); err != nil {
		log.Println("could not display fatal error dialog:", err)
	}

	// Log error message and exit with non-zero exit code
	log.Fatalln(body)
}
