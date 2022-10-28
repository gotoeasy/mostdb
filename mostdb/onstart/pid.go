package onstart

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type pidFile struct {
	Path  string // pid文件路径
	Pid   string // pid
	IsNew bool   // 是否新
	Err   error  // error
}

func checkPidFile(path string) *pidFile {
	pid := readPid(path)
	if pid == "" {
		return nil
	}
	if _, err := os.Stat(filepath.Join("/proc", pid)); err == nil {
		return &pidFile{
			Path:  path,
			Pid:   pid,
			IsNew: false,
		}
	}
	return nil
}

func readPid(path string) string {
	if pidByte, err := os.ReadFile(path); err == nil {
		return strings.TrimSpace(string(pidByte))
	}
	return ""
}

func savePid(path string, pid string) error {
	if err := os.MkdirAll(filepath.Dir(path), os.FileMode(0755)); err != nil {
		log.Println("create pid file failed", path)
		return err
	}

	if err := os.WriteFile(path, []byte(pid), 0644); err != nil {
		log.Println("save pid file failed", path)
		return err
	}

	return nil
}
