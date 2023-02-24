package utils

import (
	"math/rand"
	"time"
)

const (
	ColorRed           = "\u001b[31m"
	ColorGreen         = "\u001b[32m"
	ColorYellow        = "\u001b[33m"
	ColorBlue          = "\u001b[34m"
	ColorMagenta       = "\u001b[35m"
	ColorCyan          = "\u001b[36m"
	ColorWhite         = "\u001b[37m"
	ColorBrightRed     = "\u001b[31;1m"
	ColorBrightGreen   = "\u001b[32;1m"
	ColorBrightYellow  = "\u001b[33;1m"
	ColorBrightBlue    = "\u001b[34;1m"
	ColorBrightMagenta = "\u001b[35;1m"
	ColorBrightCyan    = "\u001b[36;1m"
	ColorBrightWhite   = "\u001b[37;1m"

	ColorBackgroundBlack = "\u001b[40m"
	ColorReset           = "\u001b[0m"
)

var (
	colors = []string{
		ColorRed,
		ColorGreen,
		ColorYellow,
		ColorBlue,
		ColorMagenta,
		ColorCyan,
		ColorWhite,
		ColorBrightRed,
		ColorBrightGreen,
		ColorBrightYellow,
		ColorBrightBlue,
		ColorBrightMagenta,
		ColorBrightCyan,
		ColorBrightWhite,
	}

	s = rand.NewSource(time.Now().UnixNano())
)

func GetRandomANSIColor() string {
	return colors[rand.New(s).Intn(len(colors))]
}
