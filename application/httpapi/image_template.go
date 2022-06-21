package httpapi

import "github.com/projectxpolaris/youphoto/database"

type BaseImageTemplate struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
	Width     uint   `json:"width"`
	Height    uint   `json:"height"`
	Domain    string `json:"domain"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
	BlurHash  string `json:"blurHash"`
}

func NewBaseImageTemplate(data *database.Image) BaseImageTemplate {
	return BaseImageTemplate{
		Id:        data.Model.ID,
		Name:      data.Name,
		Thumbnail: data.Thumbnail,
		Width:     data.Width,
		Height:    data.Height,
		Created:   data.Model.CreatedAt.Format(TimeFormat),
		Updated:   data.Model.UpdatedAt.Format(TimeFormat),
		Domain:    data.Domain,
		BlurHash:  data.BlurHash,
	}
}

func NewBaseImageTemplateList(data []*database.Image) []BaseImageTemplate {
	result := make([]BaseImageTemplate, len(data))
	for i, v := range data {
		result[i] = NewBaseImageTemplate(v)
	}
	return result
}
