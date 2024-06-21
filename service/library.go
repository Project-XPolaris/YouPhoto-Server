package service

import (
	"errors"
	"github.com/projectxpolaris/youphoto/config"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/utils"
	"gorm.io/gorm"
	"os"
	path2 "path"
)

func CreateLibrary(name string, path string, userId uint, isPublic bool) (*database.Library, error) {
	libraryPath := path
	if len(libraryPath) == 0 {
		libraryPath = path2.Join(config.Instance.PrivateLibraryPath, name)
		isPublic = false
		if !utils.CheckFileExist(libraryPath) {
			err := os.Mkdir(libraryPath, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}
	}
	if !utils.CheckFileExist(libraryPath) {
		return nil, errors.New("path not exist")
	}
	// check library is exist
	var count int64
	err := database.Instance.Model(&database.Library{}).Where("path = ?", libraryPath).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count > 1 {
		return nil, errors.New("library is exist")
	}

	library := &database.Library{
		Name:   name,
		Path:   libraryPath,
		Public: isPublic,
	}
	err = database.Instance.Create(&library).Error
	if err != nil {
		return nil, err
	}
	if userId != 0 && !isPublic {
		err = database.Instance.Model(&library).Association("Users").Append(&database.User{
			Model: gorm.Model{ID: userId},
		})
		if err != nil {
			return nil, err
		}
	}
	return library, err
}

type LibraryQueryBuilder struct {
	Page     int
	PageSize int
	Order    string `hsource:"query" hname:"order"`
	Preview  string `hsource:"query" hname:"preview"`
	UserId   uint
}

func (q *LibraryQueryBuilder) Query() ([]*database.Library, int64, error) {
	var libraries []*database.Library
	var count int64
	query := database.Instance.Model(&database.Library{})
	if q.Page == 0 {
		q.Page = 1
	}
	if q.PageSize == 0 {
		q.PageSize = 10
	}
	if q.Order == "" {
		q.Order = "id desc"
	}
	err := query.
		Offset((q.Page-1)*q.PageSize).
		Limit(q.PageSize).
		Order(q.Order).
		Joins("LEFT JOIN library_users lu on libraries.id = lu.library_id").
		Where("lu.user_id = ?", q.UserId).
		Or("public = ?", true).
		Find(&libraries).Offset(-1).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return libraries, count, nil
}

func GetLibraryWithUser(id uint, userId uint) (*database.Library, error) {
	library := &database.Library{}
	err := database.Instance.
		Where("id = ?", id).
		Preload("Users", "id = ?", userId).
		First(library).Error
	if err != nil {
		return nil, err
	}
	for _, user := range library.Users {
		if user.ID == userId {
			return library, nil
		}
	}
	return library, err
}
