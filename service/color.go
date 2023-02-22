package service

import (
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/utils"
)

type MatchColorOption struct {
	Rank1Color  []string `json:"rank1Color"`
	Rank2Color  []string `json:"rank2Color"`
	Rank3Color  []string `json:"rank3Color"`
	MaxDistance float64  `json:"maxDistance"`
	Limit       int      `json:"limit"`
}

type MatchColorResult struct {
	Image *database.Image
	Color []*database.ImageColor
	Score float64
	Rank1 float64
	Rank2 float64
	Rank3 float64
}

func getMinDistanceColor(option []string, color *database.ImageColor) float64 {
	minDistance := 1000000.0
	for _, rankColor := range option {
		distance := utils.ColorDistance(color.Value, rankColor)
		if distance < minDistance {
			minDistance = distance
		}
	}
	return minDistance

}
func MatchColor(option MatchColorOption) ([]*MatchColorResult, error) {
	var images []*database.Image
	err := database.Instance.Preload("ImageColor").Find(&images).Error
	if err != nil {
		return nil, err
	}
	result := make([]*MatchColorResult, 0)
	for _, curImage := range images {
		rank1 := false
		rank2 := false
		rank3 := false
		if len(curImage.ImageColor) == 0 {
			continue
		}
		rank1Distance := 1000000.0
		rank2Distance := 1000000.0
		rank3Distance := 1000000.0
		for _, color := range curImage.ImageColor {
			switch color.Rank {
			case 0:
				if len(option.Rank1Color) > 0 {
					rank1Distance = getMinDistanceColor(option.Rank1Color, color)
					if rank1Distance < option.MaxDistance {
						rank1 = true
					}
				} else {
					rank1 = true
				}
			case 1:
				if len(option.Rank2Color) > 0 {
					rank2Distance = getMinDistanceColor(option.Rank2Color, color)
				} else {
					rank2 = true
				}
			case 2:
				if len(option.Rank3Color) > 0 {
					for _, rankColor := range option.Rank3Color {
						if utils.ColorDistance(color.Value, rankColor) < option.MaxDistance {
							rank3 = true
							break
						}
					}
				} else {
					rank3 = true
				}
			}
		}
		if rank1 && rank2 && rank3 {
			matchResult := &MatchColorResult{
				Image: curImage,
				Rank1: rank1Distance,
				Rank2: rank2Distance,
				Rank3: rank3Distance,
				Color: curImage.ImageColor,
			}

			if len(option.Rank1Color) > 0 {
				matchResult.Score += rank1Distance
			}
			if len(option.Rank2Color) > 0 {
				matchResult.Score += rank2Distance
			}
			if len(option.Rank3Color) > 0 {
				matchResult.Score += rank3Distance
			}

			result = append(result, matchResult)
		}
	}
	return result, nil
}
