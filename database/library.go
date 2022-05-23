package database

import "gorm.io/gorm"

type Library struct {
	gorm.Model
	Name   string
	Path   string
	Users  []*User `gorm:"many2many:library_users;"`
	Images []*Image
	Public bool
}
