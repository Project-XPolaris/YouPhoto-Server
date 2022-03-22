package service

import (
	"errors"
	"github.com/allentom/harukap/thumbnail"
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
		return thumbnailFileName, plugins.DefaultThumbnailServicePlugin.Client.Generate(source, output, thumbnail.ThumbnailOption{
			MaxWidth:  320,
			MaxHeight: 320,
		})
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
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()
	switch ext {
	case ".jpg":
		err = jpeg.Encode(out, m, nil)
		if err != nil {
			return err
		}
	case ".png":
		err = png.Encode(out, m)
		if err != nil {
			return err
		}
	default:
		return errors.New("unknown output format")
	}
	return nil
}
