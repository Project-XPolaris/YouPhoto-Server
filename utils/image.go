package utils

import (
	"fmt"
	"github.com/EdlinOrg/prominentcolor"
	"github.com/bbrks/go-blurhash"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

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

func GetMostDomainColor(imagePath string) (string, error) {
	f, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	imageInput, _, err := image.Decode(f)
	if err != nil {
		return "", err
	}
	colorItems, err := prominentcolor.Kmeans(imageInput)
	if err != nil {
		return "", err
	}
	if len(colorItems) > 0 {
		return fmt.Sprintf("#%02x%02x%02x", colorItems[0].Color.R, colorItems[0].Color.G, colorItems[0].Color.B), nil
	}
	return "", nil
}

func GetBlurHash(imagePath string) (string, error) {
	f, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	imageInput, _, err := image.Decode(f)
	if err != nil {
		return "", err
	}
	str, _ := blurhash.Encode(4, 3, imageInput)
	if err != nil {
		return "", err
	}
	return str, nil
}
