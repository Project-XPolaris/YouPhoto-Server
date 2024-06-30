package service

import (
	"errors"
	"fmt"
	"github.com/projectxpolaris/youphoto/database"
)

func CreateAlbum(name string, uid string) (*database.Album, error) {
	tx := database.Instance.Begin()
	user, err := GetUserById(uid)
	if err != nil {
		return nil, err
	}
	album := &database.Album{
		Name:    name,
		OwnerId: user.ID,
	}
	err = tx.FirstOrCreate(album, database.Album{
		Name:    name,
		OwnerId: user.ID,
	}).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Model(user).Association("Albums").Append(album)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return album, nil
}

func RemoveAlbum(albumId uint, uid string, deleteImage bool) error {
	tx := database.Instance.Begin()
	album := &database.Album{}
	err := tx.Preload("Owner").First(album, albumId).Error
	if err != nil {
		return err
	}
	if album.Owner.Uid != uid {
		return errors.New("permission denied")
	}
	images := make([]*database.Image, 0)
	err = tx.Model(album).Association("Images").Find(&images)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(album).Association("Images").Clear()
	if err != nil {
		tx.Rollback()
		return err
	}
	if deleteImage {
		for _, image := range images {
			DeleteImageById(image.ID, deleteImage)
		}
	}
	err = tx.Model(album).Association("Users").Clear()
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Unscoped().Delete(album).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func AddImageToAlbum(albumId uint, uid string, imageId ...uint) error {
	tx := database.Instance.Begin()
	album := &database.Album{}
	err := tx.Preload("Owner").First(album, albumId).Error
	if err != nil {
		return err
	}
	if album.Owner.Uid != uid {
		return errors.New("permission denied")
	}
	images := make([]*database.Image, len(imageId))
	for i, v := range imageId {
		image := &database.Image{}
		err = tx.First(image, v).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		images[i] = image
	}
	err = tx.Model(album).Association("Images").Append(images)
	if err != nil {
		tx.Rollback()
		return err
	}
	coverImage := images[0]
	err = tx.Model(album).Update("cover_id", coverImage.ID).Error
	tx.Commit()
	return nil
}

func RemoveImageFromAlbum(albumId uint, uid string, imageId ...uint) error {
	tx := database.Instance.Begin()
	album := &database.Album{}
	err := tx.Preload("Owner").First(album, albumId).Error
	if err != nil {
		return err
	}
	if album.Owner.Uid != uid {
		return errors.New("permission denied")
	}
	images := make([]*database.Image, len(imageId))
	for i, v := range imageId {
		image := &database.Image{}
		err = tx.First(image, v).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		images[i] = image
	}
	err = tx.Model(album).Association("Images").Delete(images)
	if err != nil {
		tx.Rollback()
		return err
	}
	// check if need to update cover
	if album.CoverId == images[0].ID {
		var newCover *database.Image
		err = tx.Model(album).Association("Images").Find(&newCover)
		if err != nil {
			tx.Rollback()
		}
		err = tx.Model(album).Update("cover_id", newCover.ID).Error
		if err != nil {
			tx.Rollback()
		}
	}
	tx.Commit()
	return nil
}

func UpdateAlbumName(albumId uint, uid string, name string) error {
	tx := database.Instance.Begin()
	album := &database.Album{}
	err := tx.Preload("Owner").First(album, albumId).Error
	if err != nil {
		return err
	}
	if album.Owner.Uid != uid {
		return errors.New("permission denied")
	}
	album.Name = name
	err = tx.Save(album).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

type AlbumQueryBuilder struct {
	Page       int
	PageSize   int
	NameSearch string `hsource:"query" hname:"nameSearch"`
	Uid        string `hsource:"query" hname:"uid"`
}

func (q *AlbumQueryBuilder) Query() ([]*database.Album, int64, error) {
	var albums []*database.Album
	var count int64
	query := database.Instance.Model(&database.Album{})
	if q.Page == 0 {
		q.Page = 1
	}
	if q.PageSize == 0 {
		q.PageSize = 10
	}
	if len(q.NameSearch) > 0 {
		query = query.Where("name like ?", fmt.Sprintf("%%%s%%", q.NameSearch))
	}
	if len(q.Uid) > 0 {
		query = query.Joins("left join album_users on album_users.album_id = albums.id").
			Joins("left join users on users.id = album_users.user_id").
			Where("users.uid = ?", q.Uid)
	}
	err := query.
		Offset((q.Page - 1) * q.PageSize).
		Limit(q.PageSize).
		Find(&albums).
		Offset(-1).
		Count(&count).
		Error
	if err != nil {
		return nil, 0, err
	}
	return albums, count, nil
}

func GetAlbumById(id uint, uid string) (*database.Album, error) {
	album := &database.Album{}
	err := database.Instance.Preload("Owner").First(album, id).Error
	if err != nil {
		return nil, err
	}
	if album.Owner.Uid != uid {
		return nil, errors.New("permission denied")
	}
	return album, nil
}
