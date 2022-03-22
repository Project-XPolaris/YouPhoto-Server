package httpapi

import "github.com/projectxpolaris/youphoto/database"

type BaseImageTemplate struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
}

func NewBaseImageTemplate(data *database.Image) BaseImageTemplate {
	return BaseImageTemplate{
		Id:        data.Model.ID,
		Name:      data.Name,
		Thumbnail: data.Thumbnail,
		Created:   data.Model.CreatedAt.Format(TimeFormat),
		Updated:   data.Model.UpdatedAt.Format(TimeFormat),
	}
}

func NewBaseImageTemplateList(data []*database.Image) []BaseImageTemplate {
	result := make([]BaseImageTemplate, len(data))
	for i, v := range data {
		result[i] = NewBaseImageTemplate(v)
	}
	return result
}
