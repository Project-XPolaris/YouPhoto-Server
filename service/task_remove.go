package service

import (
	"context"
	"github.com/allentom/harukap/module/task"
	"github.com/project-xpolaris/youplustoolkit/youlog"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/module"
	"github.com/projectxpolaris/youphoto/plugins"
	"github.com/projectxpolaris/youphoto/utils"
)

type RemoveLibraryTaskOutput struct {
	Id      uint   `json:"id"`
	Path    string `json:"path"`
	Name    string `json:"name"`
	Total   int64  `json:"total"`
	Current int64  `json:"current"`
}

type RemoveLibraryTaskOption struct {
	LibraryId  uint
	OnError    func(task *RemoveLibraryTask, err error)
	OnComplete func(task *RemoveLibraryTask)
}

type RemoveLibraryTask struct {
	*task.BaseTask
	option RemoveLibraryTaskOption
	output *RemoveLibraryTaskOutput
	Logger *youlog.Scope
}

func (t *RemoveLibraryTask) Stop() error {
	return nil
}

func (t *RemoveLibraryTask) Start() error {
	return nil
}

func (t *RemoveLibraryTask) Output() (interface{}, error) {
	return t.output, nil
}

func (t *RemoveLibraryTask) Done() {
	t.Status = TaskStatusDone
	if t.option.OnComplete != nil {
		t.option.OnComplete(t)
	}
}
func (t *RemoveLibraryTask) AbortError(err error) {
	t.Status = TaskStatusError
	t.Err = err
	if t.option.OnError != nil {
		t.option.OnError(t, err)
	}
}

func (t *RemoveLibraryTask) GetOutput() interface{} {
	return t.output
}

func CreateRemoveLibraryTask(option RemoveLibraryTaskOption) (*RemoveLibraryTask, error) {
	for _, task := range module.Task.Pool.Tasks {
		if removeOutput, ok := task.(*RemoveLibraryTask); ok && removeOutput.output.Id == option.LibraryId {
			if task.GetStatus() == TaskStatusRunning {
				return removeOutput, nil
			}
			module.Task.Pool.RemoveTaskById(task.GetId())
			break
		}
	}
	task := RemoveLibraryTask{
		BaseTask: task.NewBaseTask(TaskTypeRemove, "-1", TaskStatusRunning),
		option:   option,
	}
	var library database.Library
	err := database.Instance.Preload("Images").Find(&library, option.LibraryId).Error
	if err != nil {
		return nil, err
	}
	output := RemoveLibraryTaskOutput{
		Id:   library.ID,
		Path: library.Path,
		Name: library.Name,
	}
	task.output = &output

	task.Logger = plugins.DefaultYouLogPlugin.Logger.NewScope("Task").WithFields(youlog.Fields{
		"path":      library.Path,
		"libraryId": library.ID,
	})
	go func() {
		err = database.Instance.Model(&database.Image{}).Where("library_id = ?", library.ID).Count(&output.Total).Error
		if err != nil {
			task.AbortError(err)
			return
		}
		// delete colors
		err = database.Instance.Unscoped().
			Model(&database.ImageColor{}).
			Where("image_colors.image_id in (?)", database.Instance.
				Table("images").
				Select("images.id as img_id").
				Where("library_id = ?", library.ID),
			).
			Delete(&database.ImageColor{}).Error
		if err != nil {
			task.AbortError(err)
			return
		}
		err = database.Instance.Unscoped().Model(&database.Image{}).Where("library_id = ?", library.ID).Delete(database.Image{}).Error
		if err != nil {
			task.AbortError(err)
			return
		}
		for _, image := range library.Images {
			output.Current++
			plugins.GetDefaultStorage().Delete(context.Background(), utils.DefaultBucket, utils.GetThumbnailsPath(image.Thumbnail))
		}
		err = database.Instance.Unscoped().Model(&database.Library{}).Where("id = ?", library.ID).Delete(library).Error
		if err != nil {
			task.AbortError(err)
			return
		}
		task.Status = TaskStatusDone
		task.Done()
	}()
	module.Task.Pool.AddTask(&task)
	return &task, nil
}
