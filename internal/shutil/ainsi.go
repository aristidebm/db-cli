package shutil

import (
	"fmt"
)

// Color output helpers
func ColorRed(text string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", text)
}

func ColorGreen(text string) string {
	return fmt.Sprintf("\033[1;32m%s\033[0m", text)
}

func ColorYellow(text string) string {
	return fmt.Sprintf("\033[1;33m%s\033[0m", text)
}
