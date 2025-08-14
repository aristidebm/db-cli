package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// Driver and utility functions
func getSupportedDrivers() []string {
	return []string{
		"postgres://",
		"postgresql://",
		"mysql://",
		"sqlite3:///",
		"redis://",
		"influx://",
		"influxdb://",
	}
}

func getDriverFromURL(url string) string {
	if strings.Contains(url, "://") {
		return strings.Split(url, "://")[0]
	}
	if strings.Contains(url, ":") {
		return strings.Split(url, ":")[0]
	}
	return ""
}

func isDriverSupported(driver string) bool {
	supportedDrivers := map[string]bool{
		"postgres":   true,
		"postgresql": true,
		"mysql":      true,
		"sqlite3":    true,
		"redis":      true,
		"influx":     true,
		"influxdb":   true,
	}
	return supportedDrivers[driver]
}
