package service

import (
	"context"
	"fmt"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/plugins"
	"github.com/projectxpolaris/youphoto/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	image2 "image"
	"os"
	"sort"
	"strings"
	"time"

	"path/filepath"
)

type ProcessImageOption struct {
	ForceRefreshDomainColor  bool `json:"forceRefreshDomainColor"`
	ForceImageClassification bool `json:"forceImageClassification"`
}

func CreateImage(path string, libraryId uint, fullPath string, option *ProcessImageOption) (*database.Image, error) {
	if option == nil {
		option = &ProcessImageOption{}
	}
	var image database.Image
	// check if it exists
	err := database.Instance.Where("library_id = ?", libraryId).Where("path = ?", path).First(&image).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
		image = database.Image{Path: path, LibraryId: libraryId, Name: filepath.Base(path)}
	}
	md5, err := utils.GetFileMD5(fullPath)
	if err != nil {
		return nil, err
	}
	isUpdate := md5 != image.Md5
	image.Md5 = md5
	fmt.Println("md5: ", md5)
	// generate thumbnail
	if isUpdate && len(image.Thumbnail) > 0 {
		plugins.GetDefaultStorage().Delete(context.Background(), utils.DefaultBucket, utils.GetThumbnailsPath(image.Thumbnail))
		image.Thumbnail = ""
	}
	if len(image.Thumbnail) == 0 {
		thumbnailTimestart := time.Now()
		image.Thumbnail, err = GenerateThumbnail(fullPath)
		if err != nil {
			return nil, err
		}
		thumbnailTime := time.Since(thumbnailTimestart)
		fmt.Printf("thumbnail time: %s\n", thumbnailTime)
	}
	// read image info
	imageInfoTimestart := time.Now()
	width, height, _ := utils.GetImageDimension(fullPath)
	imageInfoTime := time.Since(imageInfoTimestart)
	fmt.Printf("image info time: %s\n", imageInfoTime)
	image.Width = uint(width)
	image.Height = uint(height)
	// read lastModify
	fileStat, err := os.Stat(fullPath)
	if err == nil {
		image.LastModify = fileStat.ModTime()
		image.Size = uint(fileStat.Size())
	}
	var source image2.Image
	// read image hash
	if isUpdate || len(image.AvgHash) == 0 {
		if source == nil {
			source, err = getImageFromFilePath(fullPath)
			if err != nil {
				return nil, err
			}
		}
		imageHash, _ := getImageHashFromImage(source)
		if imageHash != nil {
			image.AvgHash = fmt.Sprintf("%d", imageHash.AvgHash)
			image.DifHash = fmt.Sprintf("%d", imageHash.DifHash)
			image.PerHash = fmt.Sprintf("%d", imageHash.PerHash)
		}
	}
	go func() {
		// read blur hash
		if isUpdate || len(image.BlurHash) == 0 {
			blurHash, _ := utils.GetBlurHash(image.Thumbnail)
			if len(blurHash) > 0 {
				image.BlurHash = blurHash
			}
		}
		err = database.Instance.Save(&image).Error
		if err != nil {
			log.Error(err)
		}

		// read dominant color
		if isUpdate || len(image.Domain) == 0 || option.ForceRefreshDomainColor {
			if source == nil {
				source, err = getImageFromFilePath(fullPath)
				if err != nil {
					log.Error(err)
				}
			}
			domainColors, _ := utils.GetMostDomainColorFromImage(source)
			if len(domainColors) > 0 {
				sort.Slice(domainColors, func(i, j int) bool {
					return domainColors[i].Cnt > domainColors[j].Cnt
				})
				image.Domain = fmt.Sprintf("#%02x%02x%02x", domainColors[0].Color.R, domainColors[0].Color.G, domainColors[0].Color.B)
			}
			if domainColors != nil && len(domainColors) > 0 {
				colorToInsert := make([]database.ImageColor, 0)
				totalCnt := 0
				for _, color := range domainColors {
					totalCnt += color.Cnt
				}
				for idx, color := range domainColors {
					colorToInsert = append(colorToInsert, database.ImageColor{
						ImageId: image.ID,
						Value:   fmt.Sprintf("#%02x%02x%02x", color.Color.R, color.Color.G, color.Color.B),
						Cnt:     color.Cnt,
						Rank:    idx,
						Percent: float64(color.Cnt) / float64(totalCnt),
						R:       int(color.Color.R),
						G:       int(color.Color.G),
						B:       int(color.Color.B),
					})
				}
				database.Instance.Unscoped().Where("image_id = ?", image.ID).Delete(&database.ImageColor{ImageId: image.ID})
				database.Instance.Save(&colorToInsert)
			}
		}
		err = database.Instance.Save(&image).Error
	}()
	if isUpdate || option.ForceImageClassification {
		// read image classification
		rawFile, err := os.Open(fullPath)
		if err != nil {
			log.Error(err)
		}
		if rawFile != nil {
			predictions, _ := plugins.DefaultImageClassifyPlugin.Client.Predict(rawFile)
			savePredictionList := make([]*database.Prediction, 0)
			for _, prediction := range predictions {
				savePredictionList = append(savePredictionList, &database.Prediction{
					ImageId:     image.ID,
					Label:       prediction.Label,
					Probability: prediction.Prob,
				})
			}
			err = database.Instance.Where("image_id = ?", image.ID).Delete(&database.Prediction{}).Error
			if err != nil {
				log.Error(err)
			}
			err = database.Instance.Create(&savePredictionList).Error
			if err != nil {
				log.Error(err)
			}

		}
	}
	err = database.Instance.Save(&image).Error
	return &image, err
}

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
	ColorRank1     string  `hsource:"query" hname:"colorRank1"`
	ColorRank2     string  `hsource:"query" hname:"colorRank2"`
	ColorRank3     string  `hsource:"query" hname:"colorRank3"`
	NearAvgId      uint    `hsource:"query" hname:"nearAvgId"`
	MinAvgDistance int     `hsource:"query" hname:"minAvgDistance"`
	MaxDistance    int     `hsource:"query" hname:"maxDistance"`
	LabelSearch    string  `hsource:"query" hname:"labelSearch"`
	MaxProbability float64 `hsource:"query" hname:"maxProbability"`
	MinProbability float64 `hsource:"query" hname:"minProbability"`
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
	err := query.
		Preload("ImageColor").
		Preload("Prediction").
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
