package database

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Image struct {
	gorm.Model
	LibraryId          uint
	Name               string
	Path               string
	Thumbnail          string
	LastModify         time.Time
	Size               uint
	Width              uint
	Height             uint
	Md5                string
	Library            *Library
	Domain             string
	BlurHash           string
	AvgHash            string
	DifHash            string
	PerHash            string
	Hentai             float64 `gorm:"default:0"`
	Drawings           float64 `gorm:"default:0"`
	Neutral            float64 `gorm:"default:0"`
	Sexy               float64 `gorm:"default:0"`
	Porn               float64 `gorm:"default:0"`
	ImageColor         []*ImageColor
	Prediction         []*Prediction
	DeepdanbooruResult []*DeepdanbooruResult
	Tags               []*Tag
}

func (i *Image) GetAvgHash() (uint64, error) {
	ui64, err := strconv.ParseUint(i.AvgHash, 10, 64)
	if err != nil {
		return 0, err
	}
	return ui64, nil
}

func (i *Image) GetDifHash() (uint64, error) {
	ui64, err := strconv.ParseUint(i.DifHash, 10, 64)
	if err != nil {
		return 0, err
	}
	return ui64, nil
}

func (i *Image) GetPerHash() (uint64, error) {
	ui64, err := strconv.ParseUint(i.PerHash, 10, 64)
	if err != nil {
		return 0, err
	}
	return ui64, nil
}
