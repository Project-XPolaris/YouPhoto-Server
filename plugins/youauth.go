package plugins

import (
	"errors"
	"strconv"

	"github.com/allentom/harukap/commons"
	"github.com/allentom/harukap/plugins/youauth"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/module"
	"gorm.io/gorm"
)

var DefaultYouAuthOauthPlugin *youauth.OauthPlugin

func CreateYouAuthPlugin() {
	DefaultYouAuthOauthPlugin = &youauth.OauthPlugin{}
	DefaultYouAuthOauthPlugin.AuthFromToken = func(token string) (commons.AuthUser, error) {
		return GetUserByYouAuthToken(token)
	}
	DefaultYouAuthOauthPlugin.PasswordAuthUrl = "/oauth/youauth/password"
	module.Auth.Plugins = append(module.Auth.Plugins,
		DefaultYouAuthOauthPlugin.GetOauthPlugin(),
		DefaultYouAuthOauthPlugin.GetPasswordPlugin(),
	)
}
func SaveUserByYouAuthToken(accessToken string) (*database.User, error) {
	youAuthUser, err := DefaultYouAuthOauthPlugin.Client.GetCurrentUser(accessToken)
	if err != nil {
		return nil, err
	}
	// 创建用户和oauth记录
	uid := strconv.Itoa(youAuthUser.Id)
	var user *database.User
	err = database.Instance.Model(&database.User{}).Where("uid = ?", uid).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &database.User{
				Username: youAuthUser.Username,
				Uid:      uid,
			}
			err = database.Instance.Create(user).Error
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	oauthRecord := database.Oauth{
		Provider:    "youauth",
		AccessToken: accessToken,
		Uid:         strconv.Itoa(youAuthUser.Id),
		UserId:      user.ID,
	}
	err = database.Instance.Create(&oauthRecord).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByYouAuthToken(accessToken string) (*database.User, error) {
	var oauthRecord database.Oauth
	// 检查是否存在，不存在则创建用户和oauth记录
	err := database.Instance.Model(&database.Oauth{}).Preload("User").Where("access_token = ?", accessToken).
		Where("provider = ?", "youauth").
		Find(&oauthRecord).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 向youauth请求用户信息
			user, err := SaveUserByYouAuthToken(accessToken)
			if err != nil {
				return nil, err
			}
			return user, nil
		}
		return nil, err
	}
	if oauthRecord.User == nil {
		// 向youauth请求用户信息
		user, err := SaveUserByYouAuthToken(accessToken)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	_, err = DefaultYouAuthOauthPlugin.Client.GetCurrentUser(accessToken)
	if err != nil {
		return nil, err
	}
	return oauthRecord.User, nil

}
