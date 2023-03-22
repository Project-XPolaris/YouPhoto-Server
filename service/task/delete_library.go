package task

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
	option  RemoveLibraryTaskOption
	output  *RemoveLibraryTaskOutput
	Logger  *youlog.Scope
	Library *database.Library
}

func (t *RemoveLibraryTask) Stop() error {
	return nil
}

func (t *RemoveLibraryTask) Start() error {
	err := database.Instance.Model(&database.Image{}).Where("library_id = ?", t.Library.ID).Count(&t.output.Total).Error
	if err != nil {
		return t.AbortError(err)
	}
	// delete colors
	err = database.Instance.Unscoped().
		Model(&database.ImageColor{}).
		Where("image_colors.image_id in (?)", database.Instance.
			Table("images").
			Select("images.id as img_id").
			Where("library_id = ?", t.Library.ID),
		).
		Delete(&database.ImageColor{}).Error
	if err != nil {
		return t.AbortError(err)
	}
	// delete prediction
	err = database.Instance.Unscoped().
		Model(&database.Prediction{}).
		Where("predictions.image_id in (?)", database.Instance.
			Table("images").
			Select("images.id as img_id").
			Where("library_id = ?", t.Library.ID),
		).
		Delete(&database.Prediction{}).Error
	// delete deepdanbooru result
	err = database.Instance.Unscoped().
		Model(&database.DeepdanbooruResult{}).
		Where("deepdanbooru_results.image_id in (?)", database.Instance.
			Table("images").
			Select("images.id as img_id").
			Where("library_id = ?", t.Library.ID),
		).
		Delete(&database.Prediction{}).Error

	err = database.Instance.Unscoped().Model(&database.Image{}).Where("library_id = ?", t.Library.ID).Delete(database.Image{}).Error
	if err != nil {
		return t.AbortError(err)
	}
	for _, image := range t.Library.Images {
		t.output.Current++
		plugins.GetDefaultStorage().Delete(context.Background(), utils.DefaultBucket, utils.GetThumbnailsPath(image.Thumbnail))
	}
	err = database.Instance.Unscoped().Model(&database.Library{}).Where("id = ?", t.Library.ID).Delete(t.Library).Error
	if err != nil {
		return t.AbortError(err)
	}
	t.Done()
	return nil
}

func (t *RemoveLibraryTask) Output() (interface{}, error) {
	return t.output, nil
}

func (t *RemoveLibraryTask) Done() {
	if t.option.OnComplete != nil {
		t.option.OnComplete(t)
	}
	t.BaseTask.Done()
}
func (t *RemoveLibraryTask) AbortError(err error) error {
	if t.option.OnError != nil {
		t.option.OnError(t, err)
	}
	return t.BaseTask.AbortError(err)
}

func (t *RemoveLibraryTask) GetOutput() interface{} {
	return t.output
}

func CreateRemoveLibraryTask(option RemoveLibraryTaskOption) (*RemoveLibraryTask, error) {
	for _, existedTask := range module.Task.Pool.Tasks {
		if removeOutput, ok := existedTask.(*RemoveLibraryTask); ok && removeOutput.output.Id == option.LibraryId {
			if existedTask.GetStatus() == task.GetStatusText(nil, task.StatusDone) {
				return removeOutput, nil
			}
			module.Task.Pool.RemoveTaskById(existedTask.GetId())
			break
		}
	}
	newTask := RemoveLibraryTask{
		BaseTask: task.NewBaseTask(TypeRemove, "-1", task.GetStatusText(nil, task.StatusDone)),
		option:   option,
	}
	var library database.Library
	err := database.Instance.Preload("Images").Find(&library, option.LibraryId).Error
	if err != nil {
		return nil, err
	}
	newTask.Library = &library
	output := RemoveLibraryTaskOutput{
		Id:   library.ID,
		Path: library.Path,
		Name: library.Name,
	}

	newTask.output = &output

	newTask.Logger = plugins.DefaultYouLogPlugin.Logger.NewScope("Task").WithFields(youlog.Fields{
		"path":      library.Path,
		"libraryId": library.ID,
	})

	module.Task.Pool.AddTask(&newTask)
	return &newTask, nil
}
