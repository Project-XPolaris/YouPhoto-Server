package database

import "gorm.io/gorm"

type Prediction struct {
	gorm.Model
	Label       string
	Probability float64
	ImageId     uint
	Image       *Image
}
