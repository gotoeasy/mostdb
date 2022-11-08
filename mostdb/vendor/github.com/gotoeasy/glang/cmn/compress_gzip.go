package cmn

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

// 指定文件压缩为gzip文件（文件名不支持中文）
func Gzip(srcFile string, gzipFile string) error {
	inFile, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	outFile, err := os.Create(gzipFile)
	if err != nil {
		return err
	}

	gz := gzip.NewWriter(outFile)
	defer gz.Close()

	gz.Name = FileName(srcFile)

	_, err = io.Copy(gz, inFile)
	return err
}

// 解压gzip文件到指定目录
func UnGzip(gzipPathFile string, destPath string) error {
	err := os.MkdirAll(destPath, 0666)
	if err != nil {
		return err
	}
	gzFile, err := os.Open(gzipPathFile)
	if err != nil {
		return err
	}
	defer gzFile.Close()

	gz, err := gzip.NewReader(gzFile)
	if err != nil {
		return err
	}
	defer gz.Close()

	outFile, err := os.Create(filepath.Join(destPath, gz.Header.Name))
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, gz)
	return err
}
