package httpapi

import "github.com/projectxpolaris/youphoto/database"

type BaseAlbumTemplate struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Cover *uint  `json:"cover,omitempty"` // 改为指针类型
}

func NewBaseAlbumTemplate(data *database.Album) BaseAlbumTemplate {
	tp := BaseAlbumTemplate{
		Id:    data.Model.ID,
		Name:  data.Name,
		Cover: data.CoverId, // 直接赋值指针
	}
	return tp
}

func NewBaseAlbumTemplateList(data []*database.Album) []BaseAlbumTemplate {
	result := make([]BaseAlbumTemplate, len(data))

	for i, v := range data {
		result[i] = NewBaseAlbumTemplate(v)
	}
	return result
}
