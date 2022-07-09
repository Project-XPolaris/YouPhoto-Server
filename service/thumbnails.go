package service

import (
	"context"
	"github.com/allentom/harukap/plugins/thumbnail"
	"github.com/projectxpolaris/youphoto/plugins"
	"github.com/projectxpolaris/youphoto/utils"
	"github.com/rs/xid"
	"os"
	"path/filepath"
)

func GenerateThumbnail(source string) (string, error) {
	id := xid.New().String()
	thumbnailFileName := utils.ChangeFileNameWithoutExt(filepath.Base(source), id)
	output := utils.GetThumbnailsPath(thumbnailFileName)
	file, err := os.Open(source)
	if err != nil {
		return "", err
	}
	out, err := plugins.DefaultThumbnailServicePlugin.Resize(context.Background(), file, thumbnail.ThumbnailOption{
		MaxWidth:  320,
		MaxHeight: 320,
	})
	if err != nil {
		return "", err
	}
	storage := plugins.GetDefaultStorage()
	err = storage.Upload(context.Background(), out, utils.DefaultBucket, output)
	if err != nil {
		return "", err
	}
	return thumbnailFileName, nil
}
