package httpapi

import (
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/service"
)

type BaseImageTemplate struct {
	Id          uint                 `json:"id"`
	Name        string               `json:"name"`
	Thumbnail   string               `json:"thumbnail"`
	Width       uint                 `json:"width"`
	Height      uint                 `json:"height"`
	Domain      string               `json:"domain"`
	Created     string               `json:"created"`
	Updated     string               `json:"updated"`
	BlurHash    string               `json:"blurHash"`
	ImageColors []ImageColorTemplate `json:"imageColors"`
}

func NewBaseImageTemplate(data *database.Image) BaseImageTemplate {
	return BaseImageTemplate{
		Id:          data.Model.ID,
		Name:        data.Name,
		Thumbnail:   data.Thumbnail,
		Width:       data.Width,
		Height:      data.Height,
		Created:     data.Model.CreatedAt.Format(TimeFormat),
		Updated:     data.Model.UpdatedAt.Format(TimeFormat),
		Domain:      data.Domain,
		BlurHash:    data.BlurHash,
		ImageColors: NewColorTemplateList(data.ImageColor),
	}
}

func NewBaseImageTemplateList(data []*database.Image) []BaseImageTemplate {
	result := make([]BaseImageTemplate, len(data))

	for i, v := range data {
		result[i] = NewBaseImageTemplate(v)
	}
	return result
}

type NearImageTemplate struct {
	Image       BaseImageTemplate `json:"image"`
	AvgDistance int               `json:"avgDistance"`
}

func NewNearImageTemplate(data *service.DistanceImage) NearImageTemplate {
	return NearImageTemplate{
		Image:       NewBaseImageTemplate(data.Image),
		AvgDistance: data.AvgDistance,
	}
}

func NewNearImageTemplateList(data []*service.DistanceImage) []NearImageTemplate {
	result := make([]NearImageTemplate, len(data))
	for i, v := range data {
		result[i] = NewNearImageTemplate(v)
	}
	return result
}

type ImageColorTemplate struct {
	Value   string  `json:"value"`
	ImageId uint    `json:"imageId"`
	Percent float64 `json:"percent"`
	Rank    int     `json:"rank"`
	Cnt     int     `json:"cnt"`
}

func NewColorTemplate(data *database.ImageColor) ImageColorTemplate {
	return ImageColorTemplate{
		Value:   data.Value,
		ImageId: data.ImageId,
		Percent: data.Percent,
		Rank:    data.Rank,
		Cnt:     data.Cnt,
	}
}
func NewColorTemplateList(data []*database.ImageColor) []ImageColorTemplate {
	result := make([]ImageColorTemplate, len(data))
	for i, v := range data {
		result[i] = NewColorTemplate(v)
	}
	return result
}

type ColorMatchTemplate struct {
	Image  BaseImageTemplate    `json:"image"`
	Colors []ImageColorTemplate `json:"colors"`
	Score  float64              `json:"score"`
	Rank1  float64              `json:"rank1"`
	Rank2  float64              `json:"rank2"`
	Rank3  float64              `json:"rank3"`
}

func NewColorMatchTemplate(data *service.MatchColorResult) ColorMatchTemplate {
	return ColorMatchTemplate{
		Image:  NewBaseImageTemplate(data.Image),
		Colors: NewColorTemplateList(data.Color),
		Score:  data.Score,
		Rank1:  data.Rank1,
		Rank2:  data.Rank2,
		Rank3:  data.Rank3,
	}
}

func NewColorMatchTemplateList(data []*service.MatchColorResult) []ColorMatchTemplate {
	result := make([]ColorMatchTemplate, len(data))
	for i, v := range data {
		result[i] = NewColorMatchTemplate(v)
	}
	return result
}
