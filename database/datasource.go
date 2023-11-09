package database

import (
	"github.com/allentom/harukap/plugins/datasource"
	"gorm.io/gorm"
)

var DefaultPlugin = &datasource.Plugin{
	OnConnected: func(db *gorm.DB) {
		Instance = db
		Instance.AutoMigrate(&Library{}, &Image{}, &Oauth{}, &User{}, &ImageColor{}, &Prediction{}, &DeepdanbooruResult{}, SdwConfig{}, LoraConfig{}, &Tag{})
	},
}
