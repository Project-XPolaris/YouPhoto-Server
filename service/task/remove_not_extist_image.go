package task

import (
	"github.com/allentom/harukap/module/task"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/service"
	"github.com/projectxpolaris/youphoto/utils"
	"gorm.io/gorm"
	"path/filepath"
)

type RemoveNotExistImageTaskOption struct {
	Uid          string
	libraryId    uint
	ParentTaskId string
}
type RemoveNotExistImageTaskOutput struct {
}
type RemoveNotExistImageTask struct {
	*task.BaseTask
	option     *RemoveNotExistImageTaskOption
	TaskOutput *RemoveNotExistImageTaskOutput
}

func (t *RemoveNotExistImageTask) Stop() error {
	return nil
}

func (t *RemoveNotExistImageTask) Start() error {
	var library *database.Library
	err := database.Instance.Model(database.Library{}).Where("id = ?", t.option.libraryId).First(&library).Error
	if err != nil {
		return t.AbortError(err)
	}
	var existImageCount int64
	err = database.Instance.Model(database.Image{}).
		Where("library_id = ?", t.option.libraryId).Count(&existImageCount).Error
	if err != nil {
		return t.AbortError(err)
	}
	for idx := 0; idx < int(existImageCount); idx += 20 {
		var images []database.Image
		err = database.Instance.Model(database.Image{}).
			Where("library_id = ?", t.option.libraryId).
			Offset(idx).
			Limit(20).
			Find(&images).Error
		if err != nil {
			return t.AbortError(err)
		}
		err = database.Instance.Transaction(func(tx *gorm.DB) error {
			for _, image := range images {
				if !utils.CheckFileExist(filepath.Join(library.Path, image.Path)) {
					err = service.DeleteImageById(image.ID)
					if err != nil {
						return err
					}
				}
			}
			return nil
		})
		if err != nil {
			return t.AbortError(err)
		}
	}
	t.Done()
	return nil
}

func (t *RemoveNotExistImageTask) Output() (interface{}, error) {
	return t.TaskOutput, nil
}

func NewRemoveNotExistImageTask(option *RemoveNotExistImageTaskOption) *RemoveNotExistImageTask {
	t := &RemoveNotExistImageTask{
		BaseTask:   task.NewBaseTask(TypeRemoveNotExistImage, option.Uid, task.GetStatusText(nil, task.StatusRunning)),
		TaskOutput: &RemoveNotExistImageTaskOutput{},
		option:     option,
	}
	t.ParentTaskId = option.ParentTaskId
	return t
}
