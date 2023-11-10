package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Uid      string `gorm:"unique"`
	Username string
	Token    string
	Library  []*Library `gorm:"many2many:library_users;"`
	Albums   []*Album   `gorm:"many2many:album_users;"`
	OwnAlbum []*Album   `gorm:"foreignKey:OwnerId"`
}
