package database

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Tag    string
	Images []*Image `gorm:"many2many:tag_images;"`
}

type TagImage struct {
	TagId   uint `gorm:"primaryKey"`
	ImageId uint `gorm:"primaryKey"`
	Source  string
	Rank    float64
}
