package cmn

import (
	"os"
	"runtime"
	"strconv"
	"strings"
)

func IsWin() bool {
	return runtime.GOOS == "windows"
}

func IsMac() bool {
	return runtime.GOOS == "darwin"
}

func IsLinux() bool {
	return runtime.GOOS == "linux"
}

func GetEnvStr(name string, defaultValue string) string {
	s := os.Getenv(name)
	if s == "" {
		return defaultValue
	}
	return s
}

func GetEnvInt(name string, defaultValue int) int {
	s := os.Getenv(name)
	if s == "" {
		return defaultValue
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return v
}

func GetEnvBool(name string, defaultValue bool) bool {
	s := os.Getenv(name)
	if s == "" {
		return defaultValue
	}

	if strings.ToLower(s) == "true" {
		return true
	}
	if strings.ToLower(s) == "false" {
		return false
	}
	return defaultValue
}
