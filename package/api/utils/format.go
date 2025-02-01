package utils

import "fmt"

func FormatTime(milliseconds int) string {
	// Convert milliseconds to minutes and seconds
	minutes := milliseconds / 60000
	seconds := (milliseconds % 60000) / 1000
	milliseconds = milliseconds % 1000

	// Return always two digits for the seconds and three for the milliseconds
	return fmt.Sprintf("%01d:%02d.%03d", minutes, seconds, milliseconds)
}
