package httpapi

import (
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
var youPlusTokenHandler haruka.RequestHandler = func(context *haruka.Context) {
	// check token is valid
	//token := context.GetQueryString("token")
	//_, err := plugins.DefaultYouAuthOauthPlugin.Client.GetCurrentUser(token)
	//if err != nil {
	//	AbortError(context, err, http.StatusBadRequest)
	//	return
	//}
	context.JSON(haruka.JSON{
		"success": true,
	})
}
