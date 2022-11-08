package cmn

import (
	"os"
	"runtime"
	"strconv"
)

// 是否Windows系统
func IsWin() bool {
	return runtime.GOOS == "windows"
}

// 是否Mac系统
func IsMac() bool {
	return runtime.GOOS == "darwin"
}

// 是否Linux系统
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// 取环境变量
func GetEnvStr(name string, defaultValue string) string {
	s := os.Getenv(name)
	if s == "" {
		return defaultValue
	}
	return s
}

// 取环境变量
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

// 取环境变量
func GetEnvBool(name string, defaultValue bool) bool {
	s := os.Getenv(name)
	if s == "" {
		return defaultValue
	}

	if ToLower(s) == "true" {
		return true
	}
	if ToLower(s) == "false" {
		return false
	}
	return defaultValue
}
