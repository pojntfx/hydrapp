package utils

import "unicode"

func Capitalize(msg string) string {
	// Capitalize the first letter of the message if it is longer than two characters
	if len(msg) >= 2 {
		return string(unicode.ToUpper([]rune(msg)[0])) + msg[1:]
	}

	return msg
}
