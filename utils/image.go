package utils

import (
	"context"
	"github.com/EdlinOrg/prominentcolor"
	"github.com/bbrks/go-blurhash"
	"github.com/projectxpolaris/youphoto/plugins"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
)

func loadImage(reader io.ReadCloser) (image.Image, error) {
	image, _, err := image.Decode(reader)
	return image, err
}
func GetImageDimension(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}
	return config.Width, config.Height, nil
}

func GetMostDomainColorFromImage(inputImage image.Image) ([]prominentcolor.ColorItem, error) {
	colorItems, err := prominentcolor.KmeansWithAll(9, inputImage, prominentcolor.ArgumentDefault, prominentcolor.DefaultSize, prominentcolor.GetDefaultMasks())
	if err != nil {
		return nil, err
	}
	return colorItems, nil
}

func GetBlurHash(thumbnailPath string) (string, error) {
	thumbnailStore := plugins.GetDefaultStorage()
	thumbnailSource, err := thumbnailStore.Get(context.Background(), DefaultBucket, GetThumbnailsPath(thumbnailPath))
	if err != nil {
		return "", err
	}
	thumbnailImage, err := loadImage(thumbnailSource)
	if err != nil {
		return "", err
	}
	str, err := blurhash.Encode(4, 3, thumbnailImage)
	if err != nil {
		return "", err
	}
	return str, nil
}
