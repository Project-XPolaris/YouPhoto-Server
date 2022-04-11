package service

import (
	"github.com/spf13/afero"
	"os"
	"strings"
)

var scanTargetExtensions = []string{
	".jpg", ".png", ".jpeg", ".bmp",
}

type ImageScanner struct {
	BasePath string
	OnHit    func(string)
}

func NewImageScanner(basePath string) *ImageScanner {
	return &ImageScanner{
		BasePath: basePath,
	}
}

func (s *ImageScanner) Scan() error {
	err := afero.Walk(AppFs, s.BasePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		match := strings.ToLower(info.Name())
		if strings.HasPrefix(match, ".") {
			return nil
		}
		for _, extension := range scanTargetExtensions {
			if strings.HasSuffix(match, extension) {
				s.OnHit(path)
			}
		}
		return nil
	})
	return err
}
