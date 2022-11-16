package httpapi

import (
	"errors"
	"github.com/allentom/haruka"
	"github.com/project-xpolaris/youplustoolkit/youlink"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/service"
	"net/http"
)

var generateAccessCodeWithYouAuthHandler haruka.RequestHandler = func(context *haruka.Context) {
	code := context.GetQueryString("code")
	accessToken, username, err := service.GenerateYouAuthToken(code)
	if err != nil {
		youlink.AbortErrorWithStatus(err, context, http.StatusInternalServerError)
		return
	}
	context.JSON(haruka.JSON{
		"success": true,
		"data": haruka.JSON{
			"accessToken": accessToken,
			"username":    username,
		},
	})
}

type LoginUserRequestBody struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	WithYouPlus bool   `json:"withYouPlus"`
}

var YouPlusLoginHandler haruka.RequestHandler = func(context *haruka.Context) {
	var err error
	requestBody := LoginUserRequestBody{}
	err = context.ParseJson(&requestBody)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	var user *database.User
	var sign string
	user, sign, err = service.YouPlusLogin(requestBody.Username, requestBody.Password)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	context.JSON(haruka.JSON{
		"success": true,
		"data": haruka.JSON{
			"accessToken": sign,
			"username":    user.Username,
		},
	})
}

type UserAuthRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var generateAccessCodeWithYouAuthPasswordHandler haruka.RequestHandler = func(context *haruka.Context) {
	var requestBody UserAuthRequestBody
	err := context.ParseJson(&requestBody)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	accessToken, username, err := service.GenerateYouAuthTokenByPassword(requestBody.Username, requestBody.Password)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	context.JSON(haruka.JSON{
		"success": true,
		"data": haruka.JSON{
			"accessToken": accessToken,
			"username":    username,
		},
	})
}
var youPlusTokenHandler haruka.RequestHandler = func(context *haruka.Context) {
	// check token is valid
	//token := context.GetQueryString("token")
	//_, err := plugins.DefaultYouAuthOauthPlugin.Client.GetCurrentUser(token)
	//if err != nil {
	//	AbortError(context, err, http.StatusBadRequest)
	//	return
	//}
	if claims, ok := context.Param["claim"]; ok {
		user := claims.(*database.User)
		context.JSON(haruka.JSON{
			"data": haruka.JSON{
				"username": user.Username,
				"id":       user.ID,
			},
			"success": true,
		})
	} else {
		AbortError(context, errors.New("need auth"), http.StatusBadRequest)
		return
	}

}
