package database

import "gorm.io/gorm"

type Album struct {
	gorm.Model
	Name    string
	Images  []*Image `gorm:"many2many:album_image;"`
	Users   []*User  `gorm:"many2many:album_users;"`
	Tags    []*Tag   `gorm:"many2many:album_tags;"`
	Cover   *Image
	CoverId *uint `gorm:"default:null"` // 设置为指针类型，允许为NULL
	OwnerId uint
	Owner   *User
}
