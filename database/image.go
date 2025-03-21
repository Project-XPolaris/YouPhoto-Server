package database

import (
	"gorm.io/gorm"
	"path"
	"strconv"
	"time"
)

type Image struct {
	gorm.Model
	LibraryId           uint
	Name                string
	Path                string
	Thumbnail           string
	LastModify          time.Time
	Size                uint
	Width               uint
	Height              uint
	Md5                 string
	Library             *Library
	Domain              string
	BlurHash            string
	AvgHash             string
	DifHash             string
	PerHash             string
	Hentai              float64 `gorm:"default:0"`
	Drawings            float64 `gorm:"default:0"`
	Neutral             float64 `gorm:"default:0"`
	Sexy                float64 `gorm:"default:0"`
	Porn                float64 `gorm:"default:0"`
	ImageColor          []*ImageColor
	Prediction          []*Prediction
	DeepdanbooruResult  []*DeepdanbooruResult
	Tags                []*Tag   `gorm:"many2many:tag_images;"`
	Albums              []*Album `gorm:"many2many:album_image;"`
	Lat                 float64
	Lng                 float64
	Fnumber             float64
	FocalLength         float64
	ISO                 float64
	Time                time.Time
	Tagged              bool `gorm:"default:false"`
	Country             string
	AdministrativeArea1 string
	AdministrativeArea2 string
	Locality            string
	Route               string
	StreetNumber        string
	Premise             string
	Address             string
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
func (i *Image) GetRealPath() (string, error) {
	if i.Library == nil {
		return "", nil
	}
	return path.Join(i.Library.Path, i.Path), nil
}
