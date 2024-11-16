package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Mkdir 在文件所在目录创建一个与文件同名的文件夹
func Mkdir(file *os.File) (path string, err error) {
	fileName := file.Name()
	filePath, err := filepath.Abs(fileName)
	if err != nil {
		return "", err
	}
	newDirPath := filepath.Join(filepath.Dir(filePath), strings.Split(filepath.Base(fileName), ".")[0])
	fmt.Println("OutputDir:", newDirPath)
	if err = os.Mkdir(newDirPath, os.ModePerm); err != nil {
		return "", err
	}
	return filepath.Abs(newDirPath)
}
