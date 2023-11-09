package database

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Tag     string
	Source  string
	Rank    float64
	ImageId uint
	Image   *Image
}
