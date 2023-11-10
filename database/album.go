package database

import "gorm.io/gorm"

type Album struct {
	gorm.Model
	Name    string
	Images  []*Image `gorm:"many2many:album_image;"`
	Users   []*User  `gorm:"many2many:album_users;"`
	OwnerId uint
	Owner   *User
}
