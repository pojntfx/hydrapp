package utils

import (
	"strings"
)

func Localhostize(url string) string {
	return strings.Replace(strings.Replace(strings.Replace(url, "127.0.0.1", "localhost", 1), "[::]", "localhost", 1), "0.0.0.0", "localhost", 1)
}
