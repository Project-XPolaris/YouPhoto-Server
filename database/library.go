package database

import "gorm.io/gorm"

type Library struct {
	gorm.Model
	Name   string
	Path   string
	Images []*Image
}
