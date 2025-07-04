package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/allentom/harukap/plugins/upscaler"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/plugins"
	"github.com/projectxpolaris/youphoto/utils"
	"gorm.io/gorm"
)

type ImagesQueryBuilder struct {
	Page           int
	PageSize       int
	LibraryId      []string `hsource:"query" hname:"libraryId"`
	Orders         []string `hsource:"query" hname:"order"`
	Random         string   `hsource:"query" hname:"random"`
	MinWidth       int      `hsource:"query" hname:"minWidth"`
	MinHeight      int      `hsource:"query" hname:"minHeight"`
	MaxWidth       int      `hsource:"query" hname:"maxWidth"`
	MaxHeight      int      `hsource:"query" hname:"maxHeight"`
	UserId         uint
	ColorRank1     string   `hsource:"query" hname:"colorRank1"`
	ColorRank2     string   `hsource:"query" hname:"colorRank2"`
	ColorRank3     string   `hsource:"query" hname:"colorRank3"`
	NearAvgId      uint     `hsource:"query" hname:"nearAvgId"`
	MinAvgDistance int      `hsource:"query" hname:"minAvgDistance"`
	MaxDistance    int      `hsource:"query" hname:"maxDistance"`
	LabelSearch    string   `hsource:"query" hname:"labelSearch"`
	MaxProbability float64  `hsource:"query" hname:"maxProbability"`
	MinProbability float64  `hsource:"query" hname:"minProbability"`
	NSFW           bool     `hsource:"query" hname:"nsfw"`
	NSFWMax        float64  `hsource:"query" hname:"nsfwMax"`
	NSFWMin        float64  `hsource:"query" hname:"nsfwMin"`
	DbTag          []string `hsource:"query" hname:"dbTag"`
	DbTagNot       []string `hsource:"query" hname:"dbTagNot"`
	Md5            []string `hsource:"query" hname:"md5"`
	Orient         string   `hsource:"query" hname:"orient"`
	AlbumId        uint     `hsource:"query" hname:"albumId"`
	Tag            []string `hsource:"query" hname:"tag"`
	TagNot         []string `hsource:"query" hname:"tagNot"`
}

func (q *ImagesQueryBuilder) Query() ([]*database.Image, int64, error) {
	var images []*database.Image
	var count int64

	query := database.Instance.Model(&database.Image{})
	if q.Page == 0 {
		q.Page = 1
	}
	if q.PageSize == 0 {
		q.PageSize = 10
	}
	query = query.Joins("LEFT JOIN library_users lu on images.library_id = lu.library_id").
		Joins("LEFT JOIN libraries l on l.id = images.library_id")
	if len(q.LibraryId) > 0 {
		query = query.Where("images.library_id IN ? and (l.public = ? or lu.user_id = ?)", q.LibraryId, true, q.UserId)
	} else {
		query = query.Where("l.public = ? or lu.user_id = ?", true, q.UserId)
	}
	if q.MinWidth > 0 {
		query = query.Where("images.width >= ?", q.MinWidth)
	}
	if q.MinHeight > 0 {
		query = query.Where("images.height >= ?", q.MinHeight)
	}
	if q.MaxWidth > 0 {
		query = query.Where("images.width <= ?", q.MaxWidth)
	}
	if q.MaxHeight > 0 {
		query = query.Where("images.height <= ?", q.MaxHeight)
	}
	if len(q.Md5) > 0 {
		query = query.Where("images.md5 IN ?", q.Md5)
	}
	orientDelta := 200
	switch q.Orient {
	case "square":
		query = query.Where("COALESCE(images.width, 0) <= COALESCE(images.height, 0) + ? AND COALESCE(images.height, 0) <= COALESCE(images.width, 0) + ?", orientDelta, orientDelta)
	case "portrait":
		query = query.Where("COALESCE(images.width, 0) < COALESCE(images.height, 0) - 200")
	case "landscape":
		query = query.Where("COALESCE(images.width, 0) - 200 > COALESCE(images.height, 0)")
	default:
		break
	}
	if q.Tag != nil || q.TagNot != nil {
		tagFilterTable := database.Instance.
			Table("tags")
		if q.Tag != nil {
			orQuery := database.Instance
			for _, tag := range q.Tag {
				orQuery = orQuery.Or("tags.tag like ?", fmt.Sprintf("%%%s%%", tag))
			}
			tagFilterTable = tagFilterTable.Where(orQuery)
		}
		if q.TagNot != nil {
			notTagQuery := database.Instance
			for _, notTag := range q.TagNot {
				notTagQuery = notTagQuery.Where("tags.tag not like ?", fmt.Sprintf("%%%s%%", notTag))
			}
			tagFilterTable = tagFilterTable.Where(notTagQuery)
		}
		tagFilterTable = tagFilterTable.Group("tag_images.image_id").Having("count(distinct tags.tag) = ?", len(q.Tag))
		query = query.Where("images.id in (?)", tagFilterTable.Select("tag_images.image_id").Joins("INNER JOIN tag_images on tag_images.tag_id = tags.id"))
	}
	if q.AlbumId != 0 {
		query = query.Joins("INNER JOIN album_image on album_image.image_id = images.id").
			Where("album_image.album_id = ?", q.AlbumId)
	}
	colorTablesQueryStringParts := make([]string, 0)
	colorSubQueryTable := make([]interface{}, 0)
	colorQueryTableNames := make([]string, 0)
	if len(q.ColorRank1) > 0 {
		r, g, b := utils.HexToRGB(q.ColorRank1)
		colorTablesQueryStringParts = append(colorTablesQueryStringParts, "(?) as rank1")
		colorSubQueryTable = append(colorSubQueryTable, database.Instance.
			Table("image_colors").
			Select("sqrt(pow(image_colors.r - ?, 2) +pow(image_colors.g - ?, 2) +pow(image_colors.b - ?, 2)) as distance,image_id", r, g, b).
			Where("image_colors.rank = 0"))
		colorQueryTableNames = append(colorQueryTableNames, "rank1")
	}

	if len(q.ColorRank2) > 0 {
		colorTablesQueryStringParts = append(colorTablesQueryStringParts, "(?) as rank2")
		r, g, b := utils.HexToRGB(q.ColorRank2)
		colorSubQueryTable = append(colorSubQueryTable, database.Instance.
			Table("image_colors").
			Select("sqrt(pow(image_colors.r - ?, 2) +pow(image_colors.g - ?, 2) +pow(image_colors.b - ?, 2)) as distance,image_id", r, g, b).
			Where("image_colors.rank = 1"))
		colorQueryTableNames = append(colorQueryTableNames, "rank2")
	}
	if len(q.ColorRank3) > 0 {
		colorTablesQueryStringParts = append(colorTablesQueryStringParts, "(?) as rank3")
		r, g, b := utils.HexToRGB(q.ColorRank3)
		colorSubQueryTable = append(colorSubQueryTable, database.Instance.
			Table("image_colors").
			Select("sqrt(pow(image_colors.r - ?, 2) +pow(image_colors.g - ?, 2) +pow(image_colors.b - ?, 2)) as distance,image_id", r, g, b).
			Where("image_colors.rank = 2"))
	}
	if len(colorTablesQueryStringParts) > 0 {
		selectAdd := make([]string, 0)
		for i := 0; i < len(colorQueryTableNames); i++ {
			selectAdd = append(selectAdd, fmt.Sprintf("%s.distance", colorQueryTableNames[i]))
		}
		totalTable := database.Instance.Table(strings.Join(colorTablesQueryStringParts, ","), colorSubQueryTable...).
			Select(fmt.Sprintf("%s as total_distance, %s.image_id as id",
				strings.Join(selectAdd, "+"),
				colorQueryTableNames[0],
			))
		for i := 1; i < len(colorQueryTableNames); i++ {
			totalTable = totalTable.Where(fmt.Sprintf("%s.image_id = %s.image_id", colorQueryTableNames[i], colorQueryTableNames[i-1]))
		}
		//Where("rank1.image_id = rank2.image_id").
		//Where("rank2.image_id = rank3.image_id").
		//Where("rank1.image_id = rank3.image_id")
		query = query.Joins("INNER JOIN (?) as total_distance on total_distance.id = images.id", totalTable).
			Where("total_distance.total_distance < ?", q.MaxDistance)
		query = query.Order("total_distance.total_distance asc")
		query = query.Preload("ImageColor")
	}

	if q.NearAvgId != 0 {
		image := database.Image{}
		err := database.Instance.Where("id = ?", q.NearAvgId).First(&image).Error
		if err != nil {
			return nil, 0, err
		}
		avgHash := image.AvgHash
		if len(avgHash) == 0 {
			return nil, 0, nil
		}
		imageHashTable := database.Instance.
			Table("images").
			Select("BIT_COUNT(images.avg_hash ^ ?) as distance, images.id as id", avgHash)
		query = query.Joins("INNER JOIN (?) as image_hash_distance on image_hash_distance.id = images.id", imageHashTable).
			Where("image_hash_distance.distance < ?", q.MinAvgDistance).
			Order("image_hash_distance.distance asc")
	}
	if len(q.LabelSearch) > 0 {
		query = query.Joins("INNER JOIN predictions on predictions.image_id = images.id").
			Where("predictions.label like ?", fmt.Sprintf("%%%s%%", q.LabelSearch))
		if q.MaxProbability > 0 {
			query = query.Where("predictions.probability <= ?", q.MaxProbability)
		}
		if q.MinProbability > 0 {
			query = query.Where("predictions.probability >= ?", q.MinProbability)
		}
	}
	if q.NSFW {
		threshold := q.NSFWMax
		if threshold == 0 {
			threshold = 0.8
		}
		query = query.Where("hentai <= ?", threshold).
			Where("sexy <= ?", threshold).
			Where("porn <= ?", threshold)
	}
	if len(q.Random) > 0 {
		if database.Instance.Dialector.Name() == "sqlite" {
			query = query.Order("random()")
		} else if database.Instance.Dialector.Name() == "mysql" {
			query = query.Order("RAND()")
		}
	} else {
		for _, order := range q.Orders {
			query = query.Order(fmt.Sprintf("images.%s", order))
		}
	}

	if q.DbTag != nil || q.DbTagNot != nil {
		dprFilterTable := database.Instance.
			Table("deepdanbooru_results").
			Distinct("deepdanbooru_results.image_id")
		if q.DbTag != nil {
			orQuery := database.Instance
			for _, tag := range q.DbTag {
				orQuery = orQuery.Or("deepdanbooru_results.tag like ?", fmt.Sprintf("%%%s%%", tag))
			}
			dprFilterTable = dprFilterTable.Where(orQuery)
		}
		if q.DbTagNot != nil {
			notTagQuery := database.Instance
			for _, notTag := range q.DbTagNot {
				notTagQuery = notTagQuery.Where("deepdanbooru_results.tag not like ?", fmt.Sprintf("%%%s%%", notTag))
			}
			dprFilterTable = dprFilterTable.Where(notTagQuery)
		}

		query = query.Joins("INNER JOIN (?) as dbrf on dbrf.image_id = images.id", dprFilterTable)
	}
	err := query.
		Offset((q.Page - 1) * q.PageSize).
		Limit(q.PageSize).
		Find(&images).
		Offset(-1).
		Count(&count).
		Error
	if err != nil {
		return nil, 0, err
	}
	return images, count, nil
}

func GetImageById(id uint, rels ...string) (*database.Image, error) {
	image := database.Image{}
	query := database.Instance
	for _, rel := range rels {
		query = query.Preload(rel)
	}
	err := query.Where("id = ?", id).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func DeleteImageById(id uint, deleteFile bool) error {
	image := database.Image{}
	err := database.Instance.Where("id = ?", id).Preload("Library").First(&image).Error
	if err != nil {
		return err
	}
	// delete color pattern
	err = database.Instance.Unscoped().Where("image_id = ?", id).Delete(&database.ImageColor{}).Error
	if err != nil {
		return err
	}
	// delete prediction
	err = database.Instance.Unscoped().Where("image_id = ?", id).Delete(&database.Prediction{}).Error
	if err != nil {
		return err
	}
	// delete deepdanbooru result
	err = database.Instance.Unscoped().Where("image_id = ?", id).Delete(&database.DeepdanbooruResult{}).Error
	if err != nil {
		return err
	}
	// delete image
	plugins.GetDefaultStorage().Delete(context.Background(), utils.DefaultBucket, utils.GetThumbnailsPath(image.Thumbnail))
	// clear if it album cover
	err = database.Instance.Exec("update albums set cover_id = null where cover_id = ?", id).Error
	if err != nil {
		return err
	}
	// clear album
	err = database.Instance.Exec("delete from album_image where image_id = ?", id).Error
	if err != nil {
		return err
	}
	err = database.Instance.Unscoped().Delete(&image).Error
	if err != nil {
		return err
	}

	if deleteFile {
		// delete file
		err = os.Remove(path.Join(image.Library.Path, image.Path))
		if err != nil {
			return err
		}
	}
	return nil
}
func DeleteImageByIds(ids []uint, deleteFile bool) error {
	for _, id := range ids {
		err := DeleteImageById(id, deleteFile)
		if err != nil {
			return err
		}
	}
	return nil
}
func TagImageById(id uint, taggerModel string, threshold float64) ([]*database.Tag, error) {
	image := database.Image{}
	err := database.Instance.Where("id = ?", id).Preload("Library").First(&image).Error
	if err != nil {
		return nil, err
	}
	imagePath := path.Join(image.Library.Path, image.Path)
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	result, err := plugins.DefaultImageTaggerPlugin.TagImage(imageFile, taggerModel, threshold)
	if err != nil {
		return nil, err
	}
	tx := database.Instance.Begin()
	var tagToAddList = make([]*database.Tag, 0)
	for _, tag := range result {
		var tagToAdd *database.Tag
		err = database.Instance.FirstOrCreate(&tagToAdd, database.Tag{
			Tag: tag.Tag,
		}).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		tagToAddList = append(tagToAddList, tagToAdd)
	}
	err = tx.Model(&image).Association("Tags").Append(tagToAddList)
	tx.Commit()
	return tagToAddList, nil
}

type TagQueryBuilder struct {
	Page       int
	PageSize   int
	NameSearch string `hsource:"query" hname:"nameSearch"`
}

func (q *TagQueryBuilder) Query() ([]*database.Tag, int64, error) {
	var tags []*database.Tag
	var count int64
	query := database.Instance.Model(&database.Tag{})
	if q.Page == 0 {
		q.Page = 1
	}
	if q.PageSize == 0 {
		q.PageSize = 10
	}
	if len(q.NameSearch) > 0 {
		query = query.Where("tag like ?", fmt.Sprintf("%%%s%%", q.NameSearch))
	}
	err := query.
		Offset((q.Page - 1) * q.PageSize).
		Limit(q.PageSize).
		Find(&tags).
		Offset(-1).
		Count(&count).
		Error
	if err != nil {
		return nil, 0, err
	}
	return tags, count, nil
}

func GetTaggerList() ([]string, error) {
	if !plugins.DefaultImageTaggerPlugin.IsEnable() {
		return []string{}, nil
	}
	state, err := plugins.DefaultImageTaggerPlugin.GetTaggerState()
	if err != nil {
		return nil, err
	}
	return state.ModelList, nil
}
func SaveUploadFile(filename string, file io.Reader, libraryId uint) (*database.Image, error) {
	// save file
	fileByteArr := make([]byte, 0)
	for {
		buf := make([]byte, 1024)
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		fileByteArr = append(fileByteArr, buf[:n]...)
	}
	// check if exists
	library := database.Library{}
	err := database.Instance.Where("id = ?", libraryId).First(&library).Error
	if err != nil {
		return nil, err
	}
	// check it has same md5 file in same library
	md5, err := utils.GetMd5FromBytes(fileByteArr)
	if err != nil {
		return nil, err
	}
	var existImage database.Image
	err = database.Instance.Model(database.Image{}).
		Where("md5 = ?", md5).
		Where("library_id = ?", libraryId).
		Where("name = ?", filename).
		Find(&existImage).Error

	if err != nil {
		return nil, err
	}
	// if exists, return nil
	if existImage.ID != 0 {
		return &existImage, err
	}

	// save file
	fullPath := filepath.Join(library.Path, filename)
	saveFileName := filename
	if utils.CheckFileExist(fullPath) {
		saveFileName = fmt.Sprintf("%s_%d%s", strings.TrimSuffix(filename, filepath.Ext(filename)), time.Now().Unix(), filepath.Ext(filename))
	}
	// saveFile
	savePath := filepath.Join(library.Path, saveFileName)
	err = os.WriteFile(savePath, fileByteArr, 0644)
	if err != nil {
		// handle error
	}

	// save image
	image := database.Image{
		Name:      filename,
		Path:      saveFileName,
		LibraryId: libraryId,
		Md5:       md5,
	}

	width, height, _ := utils.GetImageDimension(savePath)
	image.Width = uint(width)
	image.Height = uint(height)
	// read lastModify
	fileStat, err := os.Stat(savePath)
	if err == nil {
		image.LastModify = fileStat.ModTime()
		image.Size = uint(fileStat.Size())
	}

	//sourceImage, err := GetImageFromFilePath(savePath)
	//if err != nil {
	//	return err
	//}
	// generate thumbnail
	thumbnail, err := GenerateThumbnail(savePath)
	if err != nil {
		return nil, err
	}
	image.Thumbnail = thumbnail

	err = database.Instance.Create(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

type UpscaleImageOption struct {
	OutScale    float64 `hsource:"query" hname:"out_scale"`
	ModelName   string  `hsource:"query" hname:"model_name"`
	FaceEnhance bool    `hsource:"query" hname:"face_enhance"`
}

func UpscaleImage(imageId uint, option *UpscaleImageOption) ([]byte, string, error) {
	if !plugins.DefaultImageUpscalerPlugin.IsEnable() {
		return nil, "", fmt.Errorf("no image upscaler plugin")
	}
	var image database.Image
	err := database.Instance.Where("id = ?", imageId).Preload("Library").First(&image).Error
	if err != nil {
		return nil, "", err
	}
	realPath, err := image.GetRealPath()
	if err != nil {
		return nil, "", err
	}
	file, err := os.Open(realPath)
	result, err := plugins.DefaultImageUpscalerPlugin.Client.Upscale(file, &upscaler.UpscaleOptions{
		Model:       option.ModelName,
		FaceEnhance: option.FaceEnhance,
		OutScale:    option.OutScale,
	})
	if err != nil {
		return nil, "", err
	}
	return result, image.Name, nil
}

func GetImagesWithoutTags() ([]database.Image, error) {
	var imagesWithoutTags []database.Image
	// Assuming `database.Instance` is your GORM database connection
	// Adjust table and column names as per your schema
	err := database.Instance.
		Joins("LEFT JOIN tag_images ON tag_images.image_id = images.id").
		Where("tag_images.tag_id IS NULL").
		Where("images.tagged = ?", false).
		Group("images.id").
		Find(&imagesWithoutTags).Limit(1).Error

	if err != nil {
		return nil, err
	}
	return imagesWithoutTags, nil
}

func GetTopPhotosByMostOccurringField(db *gorm.DB, field string) ([]database.Image, error) {
	var results []struct {
		FieldValue string
		Count      int
	}
	var topPhotos []database.Image

	// Step 1 & 2: Group by the field and count the occurrences
	err := db.Model(&database.Image{}).
		Select(field + " AS field_value, COUNT(*) AS count").
		Group(field).
		Order("count DESC").
		Limit(5).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Step 3: For each top group, retrieve the photos
	for _, result := range results {
		var photos []database.Image
		err := db.Where(field+" = ?", result.FieldValue).
			Find(&photos).Error
		if err != nil {
			return nil, err
		}
		topPhotos = append(topPhotos, photos...)
	}

	return topPhotos, nil
}
