package service

import (
	"github.com/corona10/goimagehash"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

type ImageHash struct {
	AvgHash uint64 `json:"avgHash"`
	DifHash uint64 `json:"difHash"`
	PerHash uint64 `json:"extHash"`
}

func getImageHash(filePath string) (*ImageHash, error) {
	image, err := getImageFromFilePath(filePath)
	if err != nil {
		return nil, err
	}
	result := &ImageHash{}
	avgHash, err := goimagehash.AverageHash(image)
	if err != nil {
		return nil, err
	}
	result.AvgHash = avgHash.GetHash()
	difHash, err := goimagehash.DifferenceHash(image)
	if err != nil {
		return nil, err
	}
	result.DifHash = difHash.GetHash()
	perHash, err := goimagehash.PerceptionHash(image)
	if err != nil {
		return nil, err
	}
	result.PerHash = perHash.GetHash()
	return result, nil
}
func getImageHashFromImage(inputImage image.Image) (*ImageHash, error) {
	result := &ImageHash{}
	avgHash, err := goimagehash.AverageHash(inputImage)
	if err != nil {
		return nil, err
	}
	result.AvgHash = avgHash.GetHash()
	difHash, err := goimagehash.DifferenceHash(inputImage)
	if err != nil {
		return nil, err
	}
	result.DifHash = difHash.GetHash()
	perHash, err := goimagehash.PerceptionHash(inputImage)
	if err != nil {
		return nil, err
	}
	result.PerHash = perHash.GetHash()
	return result, nil
}
