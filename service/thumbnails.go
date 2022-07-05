package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/allentom/harukap/plugins/thumbnail"
	"github.com/nfnt/resize"
	"github.com/projectxpolaris/youphoto/config"
	"github.com/projectxpolaris/youphoto/plugins"
	"github.com/projectxpolaris/youphoto/utils"
	"github.com/rs/xid"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func GenerateThumbnail(source string) (string, error) {
	id := xid.New().String()
	thumbnailFileName := utils.ChangeFileNameWithoutExt(filepath.Base(source), id)
	output := utils.GetThumbnailsPath(thumbnailFileName)
	switch config.Instance.ThumbnailProvider {
	case "thumbnailservice":
		buf, err := plugins.DefaultThumbnailServicePlugin.Client.GenerateAsRaw(source, output, thumbnail.ThumbnailOption{
			MaxWidth:  320,
			MaxHeight: 320,
		})
		if err != nil {
			return "", err
		}
		storage := plugins.GetDefaultStorage()
		err = storage.Upload(context.Background(), buf, utils.DefaultBucket, output)
		if err != nil {
			return "", err
		}
		return thumbnailFileName, nil
	default:
		return thumbnailFileName, GenerateThumbnailWithResize(source, output)

	}
}
func GenerateThumbnailWithResize(source string, output string) error {
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()
	var img image.Image
	ext := strings.ToLower(filepath.Ext(source))
	switch ext {

	case ".jpg":
		img, err = jpeg.Decode(file)
		if err != nil {
			return err
		}
	case ".jpeg":
		img, err = jpeg.Decode(file)
		if err != nil {
			return err
		}
	case ".png":
		img, err = png.Decode(file)
		if err != nil {
			return err
		}
	}
	if img == nil {
		return errors.New("unexpect image format")
	}
	m := resize.Resize(320, 0, img, resize.NearestNeighbor)
	buf := new(bytes.Buffer)
	if err != nil {
		return err
	}
	switch ext {
	case ".jpg":
		err = jpeg.Encode(buf, m, nil)
		if err != nil {
			return err
		}
	case ".png":
		err = png.Encode(buf, m)
		if err != nil {
			return err
		}
	default:
		return errors.New("unknown output format")
	}
	storage := plugins.GetDefaultStorage()
	err = storage.Upload(context.Background(), buf, utils.DefaultBucket, output)
	if err != nil {
		return err
	}
	return nil
}
