package shutil

import (
	"os"
	"path/filepath"
)

func Getenv(name string, fallback ...string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return ""
}

func GetEditor() string {
	return Getenv("EDITOR", "vi")
}

func GetHomeDir() string {
	return Getenv("HOME")
}

func GetConfigDir() string {
	directory := Getenv("XDG_CONFIG", filepath.Join(GetHomeDir(), ".config"))
	return filepath.Join(directory, "db")
}

func GetScriptDirBySource(source string) string {
	return filepath.Join(GetConfigDir(), "scripts", source)
}
