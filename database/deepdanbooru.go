package database

import "gorm.io/gorm"

type DeepdanbooruResult struct {
	gorm.Model
	Tag     string
	Prob    float64
	ImageId uint
	Image   *Image
}
