package service

import (
	"github.com/projectxpolaris/youphoto/database"
	"gorm.io/gorm"
)

func CreateLibrary(name string, path string, userId uint) (*database.Library, error) {
	library := &database.Library{
		Name:   name,
		Path:   path,
		Public: userId == 0,
	}
	err := database.Instance.Create(&library).Error
	if err != nil {
		return nil, err
	}
	if userId != 0 {
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
		Joins("LEFT JOIN library_users lu on libraries.id = lu.library_id").
		Where("lu.user_id = ?", userId).
		Or("public = ?", true).
		First(library, id).Error
	return library, err
}
