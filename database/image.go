package database

import (
	"gorm.io/gorm"
	"time"
)

type Image struct {
	gorm.Model
	LibraryId  uint
	Name       string
	Path       string
	Thumbnail  string
	LastModify time.Time
	Size       uint
	Width      uint
	Height     uint
	Md5        string
	Library    *Library
}
