package service

import (
	"github.com/corona10/goimagehash"
	"github.com/projectxpolaris/youphoto/database"
)

type DistanceImage struct {
	Image       *database.Image
	AvgDistance int
}

func GetNearImage(sourceId uint, maxDistance int) ([]*DistanceImage, error) {
	var sourceImage database.Image
	if err := database.Instance.First(&sourceImage, sourceId).Error; err != nil {
		return nil, err
	}
	if len(sourceImage.AvgHash) == 0 {
		return nil, nil
	}
	sourceHashVal, err := sourceImage.GetAvgHash()
	if err != nil {
		return nil, err
	}
	sourceHash := goimagehash.NewImageHash(sourceHashVal, goimagehash.AHash)
	result := make([]*DistanceImage, 0)
	var pickUp []*database.Image
	err = database.Instance.Find(&pickUp).Error
	if err != nil {
		return nil, err
	}
	for _, curImage := range pickUp {
		if curImage.ID == sourceImage.ID {
			continue
		}

		if len(curImage.AvgHash) == 0 {
			continue
		}
		compareHashVal, err := curImage.GetAvgHash()
		if err != nil {
			continue
		}
		compareHash := goimagehash.NewImageHash(compareHashVal, goimagehash.AHash)
		distance, err := sourceHash.Distance(compareHash)
		if err != nil {
			continue
		}

		if distance <= maxDistance {
			result = append(result, &DistanceImage{
				Image:       curImage,
				AvgDistance: distance,
			})
		}
	}
	return result, nil
}
