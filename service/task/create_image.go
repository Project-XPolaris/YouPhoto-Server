package task

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/allentom/harukap/module/task"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/plugins"
	"github.com/projectxpolaris/youphoto/service"
	"github.com/projectxpolaris/youphoto/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	image2 "image"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type ProcessImageOption struct {
	EnableDomainColor         bool   `json:"enableDomainColor"`
	ForceRefreshDomainColor   bool   `json:"forceRefreshDomainColor"`
	EnableImageClassification bool   `json:"enableImageClassification"`
	ForceImageClassification  bool   `json:"forceImageClassification"`
	EnableNsfwCheck           bool   `json:"enableNsfwCheck"`
	ForceNsfwCheck            bool   `json:"forceNsfwCheck"`
	EnableDeepdanbooruCheck   bool   `json:"enableDeepdanbooruCheck"`
	ForceDeepdanbooruCheck    bool   `json:"forceDeepdanbooruCheck"`
	EnableTagger              bool   `json:"enableTagger"`
	ForceTagger               bool   `json:"forceTagger"`
	TaggerModel               string `json:"taggerModel"`
}
type CreateImageTaskOption struct {
	Uid          string
	LibraryId    uint
	Path         string
	FullPath     string
	ParentTaskId string
	CreateOption *ProcessImageOption
}
type CreateImageTaskOutput struct {
	Filename string `json:"filename"`
	FilePath string `json:"filePath"`
}
type CreateImageTask struct {
	*task.BaseTask
	option     *CreateImageTaskOption
	TaskOutput *CreateImageTaskOutput
	Image      *database.Image
}

func (t *CreateImageTask) Stop() error {
	return nil
}

func (t *CreateImageTask) Start() error {
	option := t.option.CreateOption
	libraryId := t.option.LibraryId
	path := t.option.Path
	fullPath := t.option.FullPath
	if option == nil {
		option = &ProcessImageOption{
			EnableImageClassification: true,
			EnableDomainColor:         true,
			EnableNsfwCheck:           true,
			EnableDeepdanbooruCheck:   true,
			EnableTagger:              true,
		}
	}
	var image database.Image
	// check if it exists
	err := database.Instance.Where("library_id = ?", libraryId).Where("Path = ?", path).First(&image).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return t.AbortError(err)
		}
		image = database.Image{
			Path: path, LibraryId: libraryId, Name: filepath.Base(path), LastModify: time.Now(),
		}
		err = database.Instance.Create(&image).Error
		if err != nil {
			return t.AbortError(err)
		}
	}
	md5, err := utils.GetFileMD5(fullPath)
	if err != nil {
		return t.AbortError(err)
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
		generateThumbnailTask := NewGenerateThumbnailTask(&GenerateThumbnailTaskOption{
			Uid:          t.option.Uid,
			ParentTaskId: t.GetId(),
			FullPath:     fullPath,
		})
		err = task.RunTask(generateThumbnailTask)
		if err != nil {
			return t.AbortError(err)
		}
		image.Thumbnail = generateThumbnailTask.thumbnail
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
	needGenerateAvgHash := isUpdate || len(image.AvgHash) == 0
	needReadDomainColor := (isUpdate || len(image.Domain) == 0 || option.ForceRefreshDomainColor) && option.EnableDomainColor

	// for reuse read image
	if needGenerateAvgHash || needReadDomainColor {
		readImageTask := NewReadImageTask(&ReadImageTaskOption{
			Uid:          t.option.Uid,
			ParentTaskId: t.GetId(),
			Path:         fullPath,
		})
		err = task.RunTask(readImageTask)
		if err != nil {
			return t.AbortError(err)
		}
		source = readImageTask.Image

	}

	// read image hash
	if needGenerateAvgHash {
		imageHash, _ := service.GetImageHashFromImage(source)
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
		if needReadDomainColor {
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
	if ((isUpdate || option.ForceImageClassification) && plugins.DefaultImageClassifyPlugin.Enable) && option.EnableImageClassification {
		imageClassifyTask := NewImageClassifyTask(&ImageClassifyTaskOption{
			Uid:          t.option.Uid,
			ParentTaskId: t.GetId(),
			Path:         fullPath,
			ImageId:      image.ID,
		})
		err = task.RunTask(imageClassifyTask)
		if err != nil {
			log.Error(err)
		}
	}
	go func() {
		if ((isUpdate || option.ForceNsfwCheck) && plugins.DefaultNSFWCheckPlugin.Enable) && option.EnableNsfwCheck {
			// read image classification
			nsfwCheckTask := NewNSFWCheckTask(&NSFWCheckTaskOption{
				Uid:          t.option.Uid,
				ParentTaskId: t.GetId(),
				Path:         fullPath,
				Image:        image,
			})
			err = task.RunTask(nsfwCheckTask)
			if err != nil {
				log.Error(err)
			}
		}
	}()
	go func() {
		if ((isUpdate || option.ForceDeepdanbooruCheck) && plugins.DefaultDeepDanbooruPlugin.Enable) && option.EnableDeepdanbooruCheck {
			deepDanbooruCheckTask := NewDeepdanbooruTask(
				&DeepdanbooruTaskOption{
					Uid:          t.option.Uid,
					ParentTaskId: t.GetId(),
					FullPath:     fullPath,
					ImageId:      image.ID,
				})
			err = task.RunTask(deepDanbooruCheckTask)
			if err != nil {
				log.Error(err)
			}
		}
	}()

	if ((isUpdate || option.ForceTagger) && plugins.DefaultImageTaggerPlugin.Enable) && option.EnableTagger {
		taggerTask := NewTaggerTask(
			&TaggerTaskOption{
				Uid:          t.option.Uid,
				ParentTaskId: t.GetId(),
				FullPath:     fullPath,
				ImageId:      image.ID,
				TaggerModel:  t.option.CreateOption.TaggerModel,
			})
		err = task.RunTask(taggerTask)
		if err != nil {
			log.Error(err)
		}
	}

	err = database.Instance.Save(&image).Error
	if err != nil {
		return t.AbortError(err)
	}
	t.Image = &image
	t.Done()
	return nil
}

func (t *CreateImageTask) Output() (interface{}, error) {
	return t.TaskOutput, nil
}

func NewCreateImageTask(option *CreateImageTaskOption) *CreateImageTask {

	t := &CreateImageTask{
		BaseTask: task.NewBaseTask(TypeCreateImage, option.Uid, task.GetStatusText(nil, task.StatusRunning)),
		TaskOutput: &CreateImageTaskOutput{
			Filename: filepath.Base(option.FullPath),
			FilePath: option.FullPath,
		},
		option: option,
	}
	t.ParentTaskId = option.ParentTaskId

	return t
}

func SaveImageByBase64(rawImage string, filename string, libraryId uint) (*database.Image, error) {
	// decode base64
	dec, err := base64.StdEncoding.DecodeString(rawImage)
	if err != nil {
		return nil, err
	}
	var library database.Library
	err = database.Instance.Where("id = ?", libraryId).First(&library).Error
	if err != nil {
		return nil, err
	}
	savePath := filepath.Join(library.Path, filename)

	// write image to file
	err = os.WriteFile(savePath, dec, 0644)
	if err != nil {
		return nil, err
	}
	// save image to database
	createImageTask := NewCreateImageTask(&CreateImageTaskOption{
		Uid:       "-1",
		LibraryId: libraryId,
		FullPath:  savePath,
		Path:      filename,
	})

	err = createImageTask.Start()
	if err != nil {
		return nil, err
	}
	return createImageTask.Image, nil
}
