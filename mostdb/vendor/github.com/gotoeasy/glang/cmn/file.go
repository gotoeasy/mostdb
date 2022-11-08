package cmn

import (
	"io"
	"os"
	"path"
	"path/filepath"
)

// 路径分隔符
func PathSeparator() string {
	return string(os.PathSeparator)
}

// 取文件名，如“abc.txt”
func FileName(name string) string {
	return path.Base(ReplaceAll(name, "\\", "/"))
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

// 删除文件或目录(含全部子目录文件)
func RemoveAllFile(pathorfile string) error {
	return os.RemoveAll(pathorfile)
}

// 复制文件
func CopyFile(srcFilePath string, dstFilePath string) error {
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	distFile, err := os.Create(dstFilePath)
	if err != nil {
		return err
	}
	defer distFile.Close()

	var tmp = make([]byte, 1024*4)
	for {
		n, err := srcFile.Read(tmp)
		distFile.Write(tmp[:n])
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

// 写文件（指定目录不存在时先创建，不含目录时存当前目录）
func WriteFileString(filename string, content string) error {
	return WriteFileBytes(filename, StringToBytes(content))
}

// 写文件（指定目录不存在时先创建，不含目录时存当前目录）
func WriteFileBytes(filename string, data []byte) error {
	os.MkdirAll(filepath.Dir(filename), 0777)
	return os.WriteFile(filename, data, 0666)
}
