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

func GetSaveFileName(dirPath string, filename string) string {
	// Ensure the directory path ends with a separator
	if !filepath.IsAbs(dirPath) {
		dirPath = filepath.Join(".", dirPath)
	}
	if dirPath[len(dirPath)-1] != filepath.Separator {
		dirPath += string(filepath.Separator)
	}

	// Combine dirPath and filename to get the full path
	fullPath := filepath.Join(dirPath, filename)

	// Check if the file exists
	if !CheckFileExist(fullPath) {
		return fullPath
	}

	// Extract the extension and name without extension
	ext := filepath.Ext(filename)
	nameWithoutExt := filename[:len(filename)-len(ext)]

	// Try adding a prefix to find a unique filename
	for i := 1; ; i++ {
		newFilename := fmt.Sprintf("%s_%d%s", nameWithoutExt, i, ext)
		newFullPath := filepath.Join(dirPath, newFilename)
		if !CheckFileExist(newFullPath) {
			return newFullPath
		}
	}

}
