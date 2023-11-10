package httpapi

import "github.com/projectxpolaris/youphoto/database"

type BaseAlbumTemplate struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

func NewBaseAlbumTemplate(data *database.Album) BaseAlbumTemplate {
	return BaseAlbumTemplate{
		Id:   data.Model.ID,
		Name: data.Name,
	}
}

func NewBaseAlbumTemplateList(data []*database.Album) []BaseAlbumTemplate {
	result := make([]BaseAlbumTemplate, len(data))

	for i, v := range data {
		result[i] = NewBaseAlbumTemplate(v)
	}
	return result
}
