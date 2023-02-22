package database

import "gorm.io/gorm"

type ImageColor struct {
	gorm.Model
	Value   string
	Cnt     int
	ImageId uint
	Percent float64
	Rank    int
	Image   *Image
	R       int
	G       int
	B       int
}
