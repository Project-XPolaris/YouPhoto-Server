package httpapi

import (
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/service"
)

type BaseImageTemplate struct {
	Id                  uint                   `json:"id"`
	Name                string                 `json:"name"`
	Thumbnail           string                 `json:"thumbnail"`
	Width               uint                   `json:"width"`
	Height              uint                   `json:"height"`
	Domain              string                 `json:"domain"`
	Created             string                 `json:"created"`
	Updated             string                 `json:"updated"`
	BlurHash            string                 `json:"blurHash"`
	ImageColors         []ImageColorTemplate   `json:"imageColors"`
	Classify            []PredictionTemplate   `json:"classify"`
	Hentai              float64                `json:"hentai"`
	Drawings            float64                `json:"drawings"`
	Neutral             float64                `json:"neutral"`
	Sexy                float64                `json:"sexy"`
	Porn                float64                `json:"porn"`
	DeepdanbooruResult  []DeepdanbooruTemplate `json:"deepdanbooruResult"`
	Tag                 []ImageTagTemplate     `json:"tag"`
	Lat                 float64                `json:"lat,omitempty"`
	Lng                 float64                `json:"lng,omitempty"`
	Fnumber             float64                `json:"fnumber,omitempty"`
	FocalLength         float64                `json:"focalLength,omitempty"`
	ISO                 float64                `json:"iso,omitempty"`
	Time                string                 `json:"time,omitempty"`
	LibraryId           uint                   `json:"libraryId,omitempty"`
	Country             string                 `json:"country,omitempty"`
	AdministrativeArea1 string                 `json:"administrativeArea1,omitempty"`
	AdministrativeArea2 string                 `json:"administrativeArea2,omitempty"`
	Locality            string                 `json:"locality,omitempty"`
	Route               string                 `json:"route,omitempty"`
	StreetNumber        string                 `json:"streetNumber,omitempty"`
	Premise             string                 `json:"premise,omitempty"`
	Address             string                 `json:"address,omitempty"`
}

func NewBaseImageTemplate(data *database.Image) BaseImageTemplate {
	return BaseImageTemplate{
		Id:                  data.Model.ID,
		Name:                data.Name,
		Thumbnail:           data.Thumbnail,
		Width:               data.Width,
		Height:              data.Height,
		Created:             data.Model.CreatedAt.Format(TimeFormat),
		Updated:             data.Model.UpdatedAt.Format(TimeFormat),
		Domain:              data.Domain,
		BlurHash:            data.BlurHash,
		ImageColors:         NewColorTemplateList(data.ImageColor),
		Classify:            NewPredictionTemplateList(data.Prediction),
		Hentai:              data.Hentai,
		Drawings:            data.Drawings,
		Neutral:             data.Neutral,
		Sexy:                data.Sexy,
		Porn:                data.Porn,
		DeepdanbooruResult:  NewDeepdanbooruTemplateList(data.DeepdanbooruResult),
		Tag:                 NewImageTagTemplateList(data.Tags),
		Lat:                 data.Lat,
		Lng:                 data.Lng,
		Fnumber:             data.Fnumber,
		FocalLength:         data.FocalLength,
		ISO:                 data.ISO,
		Time:                data.Time.Format(TimeFormat),
		LibraryId:           data.LibraryId,
		Country:             data.Country,
		AdministrativeArea1: data.AdministrativeArea1,
		AdministrativeArea2: data.AdministrativeArea2,
		Locality:            data.Locality,
		Route:               data.Route,
		StreetNumber:        data.StreetNumber,
		Premise:             data.Premise,
		Address:             data.Address,
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

type PredictionTemplate struct {
	Prob  float64 `json:"prob"`
	Label string  `json:"label"`
}

func NewPredictionTemplate(data *database.Prediction) PredictionTemplate {
	return PredictionTemplate{
		Prob:  data.Probability,
		Label: data.Label,
	}
}

func NewPredictionTemplateList(data []*database.Prediction) []PredictionTemplate {
	result := make([]PredictionTemplate, len(data))
	for i, v := range data {
		result[i] = NewPredictionTemplate(v)
	}
	return result
}

type DeepdanbooruTemplate struct {
	Tag  string  `json:"tag"`
	Prob float64 `json:"prob"`
}

func NewDeepdanbooruTemplate(data *database.DeepdanbooruResult) DeepdanbooruTemplate {
	return DeepdanbooruTemplate{
		Tag:  data.Tag,
		Prob: data.Prob,
	}
}

func NewDeepdanbooruTemplateList(data []*database.DeepdanbooruResult) []DeepdanbooruTemplate {
	result := make([]DeepdanbooruTemplate, len(data))
	for i, v := range data {
		result[i] = NewDeepdanbooruTemplate(v)
	}
	return result
}

type ImageTagTemplate struct {
	Tag     string  `json:"tag,omitempty"`
	Source  string  `json:"source,omitempty"`
	Rank    float64 `json:"rank,omitempty"`
	ImageId uint    `json:"imageId,omitempty"`
}

func NewImageTagTemplate(data *database.Tag) ImageTagTemplate {
	return ImageTagTemplate{
		Tag: data.Tag,
	}
}

func NewImageTagTemplateList(data []*database.Tag) []ImageTagTemplate {
	result := make([]ImageTagTemplate, len(data))
	for i, v := range data {
		result[i] = NewImageTagTemplate(v)
	}
	return result
}
