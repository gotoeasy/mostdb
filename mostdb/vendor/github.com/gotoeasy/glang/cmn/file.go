package cmn

import (
	"os"
	"path"
)

func PathSeparator() string {
	return string(os.PathSeparator)
}

// 取文件扩展名，如“.txt”
func FileExtName(name string) string {
	return path.Ext(name)
}

// 判断文件是否存在
func IsExistFile(file string) bool {
	s, err := os.Stat(file)
	if err == nil {
		return !s.IsDir()
	}
	if os.IsNotExist(err) {
		return false
	}
	return !s.IsDir()
}

// 判断文件夹是否存在
func IsExistDir(dir string) bool {
	s, err := os.Stat(dir)
	if err == nil {
		return s.IsDir()
	}
	if os.IsNotExist(err) {
		return false
	}
	return s.IsDir()
}
