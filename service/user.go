package service

import "github.com/projectxpolaris/youphoto/database"

var (
	PublicUid      = "-1"
	PublicUsername = "public"
)

func GetUserById(uid string) (*database.User, error) {
	var user database.User
	err := database.Instance.Where(map[string]string{"uid": uid}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
