package utils

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/ncruces/zenity"
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
		Capitalize(msg),
		Capitalize(err.Error()),
		string(debug.Stack()),
	)

	// Show error message visually using a dialog
	if err := zenity.Error(
		body,
		zenity.Title("Fatal error"),
		zenity.Width(320),
	); err != nil {
		log.Println("could not display fatal error dialog:", err)
	}

	// Log error message and exit with non-zero exit code
	log.Fatalln(body)
}
