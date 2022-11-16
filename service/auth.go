package service

import (
	"errors"
	"fmt"
	"github.com/allentom/harukap/plugins/youauth"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/plugins"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

const YouAuthProvider = "youauth"

func GenerateYouAuthToken(code string) (string, string, error) {
	tokens, err := plugins.DefaultYouAuthOauthPlugin.Client.GetAccessToken(code)
	if err != nil {
		return "", "", err
	}
	return LinkWithYouAuthToken(tokens)
}
func GenerateYouAuthTokenByPassword(username string, rawPassword string) (string, string, error) {
	authResult, err := plugins.DefaultYouAuthOauthPlugin.Client.GrantWithPassword(username, rawPassword)
	if err != nil {
		return "", "", err
	}
	return LinkWithYouAuthToken(authResult)
}
func LinkWithYouAuthToken(tokens *youauth.GenerateTokenResponse) (string, string, error) {
	currentUserResponse, err := plugins.DefaultYouAuthOauthPlugin.Client.GetCurrentUser(tokens.AccessToken)
	if err != nil {
		return "", "", err
	}
	// check if user exists
	uid := fmt.Sprintf("%d", currentUserResponse.Id)
	historyOauth := make([]database.Oauth, 0)
	err = database.Instance.Where("uid = ?", uid).
		Where("provider = ?", YouAuthProvider).
		Preload("User").
		Find(&historyOauth).Error
	if err != nil {
		return "", "", err
	}
	var user *database.User
	if len(historyOauth) == 0 {
		username := xid.New().String()
		// create new user
		user = &database.User{
			Uid:      xid.New().String(),
			Username: username,
		}
		err = database.Instance.Create(&user).Error
		if err != nil {
			return "", "", err
		}
	} else {
		user = historyOauth[0].User
	}

	oauthRecord := database.Oauth{
		Uid:          fmt.Sprintf("%d", currentUserResponse.Id),
		UserId:       user.ID,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		Provider:     YouAuthProvider,
	}
	err = database.Instance.Create(&oauthRecord).Error
	if err != nil {
		return "", "", err
	}
	return tokens.AccessToken, currentUserResponse.Username, nil
}
func refreshToken(accessToken string) (string, error) {
	tokenRecord := database.Oauth{}
	err := database.Instance.Where("access_token = ?", accessToken).First(&tokenRecord).Error
	if err != nil {
		return "", err
	}
	token, err := plugins.DefaultYouAuthOauthPlugin.Client.RefreshAccessToken(tokenRecord.RefreshToken)
	if err != nil {
		return "", err
	}
	err = database.Instance.Delete(&tokenRecord).Error
	if err != nil {
		return "", err
	}
	newOauthRecord := database.Oauth{
		UserId:       tokenRecord.UserId,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	err = database.Instance.Create(&newOauthRecord).Error
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

func YouPlusLogin(username string, rawPassword string) (*database.User, string, error) {
	authResult, err := plugins.DefaultYouPlusPlugin.Client.FetchUserAuth(username, rawPassword)
	if err != nil {
		return nil, "", err
	}
	if !authResult.Success {
		return nil, "", errors.New("user auth failed")
	}
	var oauthRecord database.Oauth
	err = database.Instance.Preload("User").Where("uid = ?", authResult.Uid).
		Where("provider = ?", "YouPlusServer").
		First(&oauthRecord).Error
	var user *database.User
	if oauthRecord.User != nil {
		user = oauthRecord.User
	}
	if err == gorm.ErrRecordNotFound {
		// create new user
		uid := xid.New().String()
		user = &database.User{
			Uid:      uid,
			Username: uid,
		}
		err = database.Instance.Create(&user).Error
		if err != nil {
			return nil, "", err
		}
	} else {
		if err != nil {
			return nil, "", err
		}
	}
	newOauth := database.Oauth{
		Uid:         authResult.Uid,
		Provider:    "YouPlusServer",
		AccessToken: authResult.Token,
		UserId:      user.ID,
	}
	err = database.Instance.Create(&newOauth).Error
	if err != nil {
		return nil, "", err
	}
	return user, authResult.Token, nil
}
