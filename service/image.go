package service

import (
	"fmt"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/utils"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"time"
)

func CreateImage(path string, libraryId uint, fullPath string) (*database.Image, error) {
	var image database.Image
	// check if it exists
	err := database.Instance.Where("library_id = ?", libraryId).Where("path = ?", path).First(&image).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
		image = database.Image{Path: path, LibraryId: libraryId, Name: filepath.Base(path)}
	}
	md5, err := utils.GetFileMD5(fullPath)
	if err != nil {
		return nil, err
	}
	isUpdate := md5 != image.Md5
	image.Md5 = md5

	// generate thumbnail
	if isUpdate && len(image.Thumbnail) > 0 {
		os.Remove(utils.GetThumbnailsPath(image.Thumbnail))
		image.Thumbnail = ""
	}
	if len(image.Thumbnail) == 0 {
		thumbnailTimestart := time.Now()
		image.Thumbnail, err = GenerateThumbnail(fullPath)
		if err != nil {
			return nil, err
		}
		thumbnailTime := time.Since(thumbnailTimestart)
		fmt.Printf("thumbnail time: %s\n", thumbnailTime)
	}
	// read image info
	imageInfoTimestart := time.Now()
	width, height, _ := utils.GetImageDimension(fullPath)
	imageInfoTime := time.Since(imageInfoTimestart)
	fmt.Printf("image info time: %s\n", imageInfoTime)
	image.Width = uint(width)
	image.Height = uint(height)
	// read lastModify
	fileStat, err := os.Stat(fullPath)
	if err == nil {
		image.LastModify = fileStat.ModTime()
		image.Size = uint(fileStat.Size())
	}
	err = database.Instance.Save(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, err
}

type ImagesQueryBuilder struct {
	Page      int
	PageSize  int
	LibraryId []string `hsource:"query" hname:"libraryId"`
	Orders    []string `hsource:"query" hname:"order"`
	Random    string   `hsource:"query" hname:"random"`
	UserId    uint
}

func (q *ImagesQueryBuilder) Query() ([]*database.Image, int64, error) {
	var images []*database.Image
	var count int64

	query := database.Instance.Model(&database.Image{})
	if q.Page == 0 {
		q.Page = 1
	}
	if q.PageSize == 0 {
		q.PageSize = 10
	}
	query = query.Joins("LEFT JOIN library_users lu on images.library_id = lu.library_id").
		Joins("LEFT JOIN libraries l on l.id = images.library_id")
	if len(q.LibraryId) > 0 {
		query = query.Where("images.library_id IN ? and (l.public = ? or lu.user_id = ?)", q.LibraryId, true, q.UserId)
	} else {
		query = query.Where("l.public = ? or lu.user_id = ?", true, q.UserId)
	}
	if len(q.Random) > 0 {
		if database.Instance.Dialector.Name() == "sqlite" {
			query = query.Order("random()")
		} else if database.Instance.Dialector.Name() == "mysql" {
			query = query.Order("RAND()")
		}
	} else {
		for _, order := range q.Orders {
			query = query.Order(fmt.Sprintf("%s", order))
		}
	}
	err := query.
		Offset((q.Page - 1) * q.PageSize).
		Limit(q.PageSize).
		Find(&images).
		Offset(-1).
		Count(&count).
		Error
	if err != nil {
		return nil, 0, err
	}
	return images, count, nil
}

func GetImageById(id uint, rels ...string) (*database.Image, error) {
	image := database.Image{}
	query := database.Instance
	for _, rel := range rels {
		query = query.Preload(rel)
	}
	err := query.Where("id = ?", id).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}
