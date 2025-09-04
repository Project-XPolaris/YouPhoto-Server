package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/projectxpolaris/youphoto/database"
	"gorm.io/gorm"
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
		CoverId: nil, // 设置为nil，表示没有封面
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
	// 如果相册还没有封面，将第一个图片设为封面
	if album.CoverId == nil {
		coverImage := images[0]
		err = tx.Model(album).Update("cover_id", coverImage.ID).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
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
	if album.CoverId != nil && *album.CoverId == images[0].ID {
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
	Tags       string `hsource:"query" hname:"tags"`
	Order      string `hsource:"query" hname:"order"`
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
	// 添加标签过滤
	if len(q.Tags) > 0 {
		// 解析标签字符串，支持逗号分隔的多个标签
		tagNames := []string{}
		for _, tag := range strings.Split(q.Tags, ",") {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tagNames = append(tagNames, tag)
			}
		}
		if len(tagNames) > 0 {
			// 使用EXISTS子查询来查找包含指定标签的相册
			subQuery := database.Instance.Table("album_tags").
				Select("1").
				Joins("left join tags on tags.id = album_tags.tag_id").
				Where("album_tags.album_id = albums.id").
				Where("tags.tag IN ?", tagNames)
			query = query.Where("EXISTS (?)", subQuery)
		}
	}
	// 添加排序
	orderBy := "id desc" // 默认排序
	if len(q.Order) > 0 {
		switch q.Order {
		case "id asc":
			orderBy = "id asc"
		case "id desc":
			orderBy = "id desc"
		case "name asc":
			orderBy = "name asc"
		case "name desc":
			orderBy = "name desc"
		default:
			orderBy = "id desc"
		}
	}
	query = query.Order(orderBy)

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

// CreateOrGetTag 创建或获取标签
func CreateOrGetTag(tagName string) (*database.Tag, error) {
	var tag database.Tag
	err := database.Instance.Where("tag = ?", tagName).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新标签
			tag = database.Tag{Tag: tagName}
			err = database.Instance.Create(&tag).Error
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &tag, nil
}

// AddTagsToAlbum 为相册添加标签
func AddTagsToAlbum(albumId uint, uid string, tagNames []string) error {
	tx := database.Instance.Begin()
	album := &database.Album{}
	err := tx.Preload("Owner").First(album, albumId).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if album.Owner.Uid != uid {
		tx.Rollback()
		return errors.New("permission denied")
	}

	// 创建或获取标签
	var tags []*database.Tag
	for _, tagName := range tagNames {
		if tagName == "" {
			continue
		}
		tag, err := CreateOrGetTag(tagName)
		if err != nil {
			tx.Rollback()
			return err
		}
		tags = append(tags, tag)
	}

	// 添加标签到相册
	if len(tags) > 0 {
		err = tx.Model(album).Association("Tags").Append(tags)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}
