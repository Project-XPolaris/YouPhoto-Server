package utils

import (
	"github.com/projectxpolaris/youphoto/config"
	"path/filepath"
)

func GetThumbnailsPath(thumbnailsName string) string {
	return filepath.Join(config.Instance.ThumbnailStorePath, thumbnailsName)
}
