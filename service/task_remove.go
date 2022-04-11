package service

import (
	"github.com/project-xpolaris/youplustoolkit/youlog"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/utils"
	"os"
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
	OnError    func(task Task, err error)
	OnComplete func(task Task)
}

type RemoveLibraryTask struct {
	BaseTask
	option RemoveLibraryTaskOption
	output *RemoveLibraryTaskOutput
}

func (t *RemoveLibraryTask) Done() {
	t.UpdateDoneStatus()
	if t.option.OnComplete != nil {
		t.option.OnComplete(t)
	}
}
func (t *RemoveLibraryTask) AbortError(err error) {
	t.BaseTask.AbortError(err)
	if t.option.OnError != nil {
		t.option.OnError(t, err)
	}
}

func (t *RemoveLibraryTask) GetOutput() interface{} {
	return t.output
}

func CreateRemoveLibraryTask(option RemoveLibraryTaskOption) (Task, error) {
	for _, task := range DefaultTaskPool.Tasks {
		if removeOutput, ok := task.GetOutput().(*RemoveLibraryTaskOutput); ok && removeOutput.Id == option.LibraryId {
			if task.GetStatus() == TaskStatusRunning {
				return task, nil
			}
			DefaultTaskPool.RemoveTaskById(task.GetId())
			break
		}
	}
	task := RemoveLibraryTask{
		BaseTask: NewBaseTask(TaskTypeRemove),
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

	task.BaseTask.Logger = task.BaseTask.Logger.WithFields(youlog.Fields{
		"path":      library.Path,
		"libraryId": library.ID,
	})
	go func() {
		err = database.Instance.Model(&database.Image{}).Where("library_id = ?", library.ID).Count(&output.Total).Error
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
			os.Remove(utils.GetThumbnailsPath(image.Thumbnail))
		}
		err = database.Instance.Unscoped().Model(&database.Library{}).Where("id = ?", library.ID).Delete(library).Error
		if err != nil {
			task.AbortError(err)
			return
		}
		task.Status = TaskStatusDone
		task.Done()
	}()
	DefaultTaskPool.AddTask(&task)
	return &task, nil
}
