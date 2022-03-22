package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func CheckFileExist(path string) bool {
	stat, _ := os.Stat(path)
	if stat != nil {
		return true
	}
	return false
}
func ChangeFileNameWithoutExt(filename string, newName string) string {
	baseName := filepath.Base(filename)
	ext := filepath.Ext(baseName)
	return fmt.Sprintf("%s%s", newName, ext)
}
