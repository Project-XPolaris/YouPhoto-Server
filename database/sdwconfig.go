package database

import "gorm.io/gorm"

type SdwConfig struct {
	gorm.Model
	Name   string
	Config string
	UserId uint
	User   *User
}
