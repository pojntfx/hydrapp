package ui

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"

	"github.com/ncruces/zenity"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/utils"
)

var (
	ErrCouldNotDisplayFatalErrorDialog = errors.New("could not display fatal error dial")
)

func HandlePanic(appName string, err error) {
	var (
		errHead = err
		errTail error
	)
	errs, ok := err.(interface{ Unwrap() []error })
	if ok {
		if e := errs.Unwrap(); len(e) > 1 {
			errHead = e[0]
			errTail = errors.Join(e[1:]...)
		}
	}

	body := ""
	if errTail == nil {
		body = fmt.Sprintf(`%v has encountered a fatal error and can't continue. The error message is:

%v

Strack trace:

%v`,
			appName,
			utils.Capitalize(errHead.Error()),
			string(debug.Stack()),
		)
	} else {
		body = fmt.Sprintf(`%v has encountered a fatal error and can't continue. The error message is:

%v

The following information might help you in fixing the problem:

%v

Strack trace:

%v`,
			appName,
			utils.Capitalize(errHead.Error()),
			utils.Capitalize(errTail.Error()),
			string(debug.Stack()),
		)
	}

	// Show error message visually using a dialog
	if err := zenity.Error(
		body,
		zenity.Title(fmt.Sprintf("Fatal error for %v", appName)),
		zenity.Width(320),
	); err != nil {
		log.Println(errors.Join(ErrCouldNotDisplayFatalErrorDialog, err))
	}

	// Log error message and exit with non-zero exit code
	log.Fatalln(body)
}
