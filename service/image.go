package service

import (
	"fmt"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/utils"
	"strings"
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
		Preload("ImageColor").
		Preload("Prediction").
		Preload("DeepdanbooruResult").
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
