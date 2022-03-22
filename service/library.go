package service

import "github.com/projectxpolaris/youphoto/database"

func CreateLibrary(name string, path string) (*database.Library, error) {
	library := &database.Library{Name: name, Path: path}
	err := database.Instance.Create(&library).Error
	return library, err
}

type LibraryQueryBuilder struct {
	Page     int
	PageSize int
	Order    string `hsource:"query" hname:"order"`
	Preview  string `hsource:"query" hname:"preview"`
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
		Offset((q.Page - 1) * q.PageSize).
		Limit(q.PageSize).
		Order(q.Order).
		Find(&libraries).Offset(-1).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return libraries, count, nil
}
