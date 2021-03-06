package plugins

import (
	"github.com/allentom/harukap/commons"
	"github.com/allentom/harukap/plugins/youplus"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/module"
)

var DefaultYouPlusPlugin *youplus.Plugin

func CreateDefaultYouPlusPlugin() {
	DefaultYouPlusPlugin = &youplus.Plugin{}
	DefaultYouPlusPlugin.AuthFromToken = func(token string) (commons.AuthUser, error) {
		return GetUserByYouPlusToken(token)
	}
	DefaultYouPlusPlugin.AuthUrl = "/oauth/youplus"
	module.Auth.Plugins = append(module.Auth.Plugins, DefaultYouPlusPlugin)
}
func GetUserByYouPlusToken(accessToken string) (*database.User, error) {
	var oauthRecord database.Oauth
	err := database.Instance.Model(&database.Oauth{}).Preload("User").Where("access_token = ?", accessToken).
		Where("provider = ?", "YouPlusServer").
		Find(&oauthRecord).Error
	if err != nil {
		return nil, err
	}
	_, err = DefaultYouPlusPlugin.Client.CheckAuth(accessToken)
	if err != nil {
		return nil, err
	}
	return oauthRecord.User, nil
}
