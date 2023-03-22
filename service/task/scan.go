package task

import (
	"errors"
	"github.com/allentom/harukap/module/task"
	"github.com/project-xpolaris/youplustoolkit/youlog"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/module"
	"github.com/projectxpolaris/youphoto/plugins"
	"github.com/projectxpolaris/youphoto/service"
	"path/filepath"
)

type SyncLibraryTask struct {
	*task.BaseTask
	option   CreateScanTaskOption
	output   *ScanTaskOutput
	library  *database.Library
	Logger   *youlog.Scope
	stopFlag bool
}

func (t *SyncLibraryTask) Stop() error {
	return nil
}

func (t *SyncLibraryTask) Start() error {
	t.Logger.Info("start sync remove missing images")
	removeNotExistsImageTask := NewRemoveNotExistImageTask(&RemoveNotExistImageTaskOption{
		libraryId:    t.library.ID,
		ParentTaskId: t.Id,
		Uid:          t.Owner,
	})
	t.SubTaskList = append(t.SubTaskList, removeNotExistsImageTask)
	err := task.RunTask(removeNotExistsImageTask)
	if err != nil {
		return t.AbortError(err)
	}

	//count total
	if t.stopFlag {
		t.Status = task.GetStatusText(nil, task.StatusDone)
		return nil
	}
	scanFileTask := NewScanImageFileTask(&ScanImageFileTaskOption{
		ParentTaskId: t.Id,
		Uid:          t.Owner,
		Path:         t.library.Path,
	})
	t.SubTaskList = append(t.SubTaskList, scanFileTask)
	t.Logger.Info("start scan library")
	err = task.RunTask(scanFileTask)
	if err != nil {
		return t.AbortError(err)
	}
	t.output.Total = int64(len(scanFileTask.pathList))
	for idx, path := range scanFileTask.pathList {
		if t.stopFlag {
			t.Done()
			return service.StopError
		}
		t.output.Current = int64(idx + 1)
		t.output.CurrentPath = path
		t.output.CurrentName = filepath.Base(path)
		imagePath, err := filepath.Rel(t.library.Path, path)
		if err != nil {
			t.AbortFileError(path, err)
			return nil
		}
		createImageTask := NewCreateImageTask(&CreateImageTaskOption{
			Uid:          t.Owner,
			path:         imagePath,
			fullPath:     path,
			CreateOption: t.option.ProcessOption,
			libraryId:    t.library.ID,
		})
		t.SubTaskList = append(t.SubTaskList, createImageTask)
		err = task.RunTask(createImageTask)
		if err != nil {

			t.AbortFileError(path, err)
			return nil

		} else {
			t.OnFileComplete()
		}
	}
	t.Done()
	return nil
}

func (t *SyncLibraryTask) Output() (interface{}, error) {
	return t.output, nil
}

type ScanTaskOutput struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Current     int64  `json:"current"`
	CurrentPath string `json:"currentPath"`
	CurrentName string `json:"currentName"`
	Total       int64  `json:"total"`
}
type CreateScanTaskOption struct {
	LibraryId      uint
	UserId         uint
	OnFileComplete func(task *SyncLibraryTask)
	OnFileError    func(task *SyncLibraryTask, err error)
	OnError        func(task *SyncLibraryTask, err error)
	OnComplete     func(task *SyncLibraryTask)
	ProcessOption  *ProcessImageOption
}

func (t *SyncLibraryTask) GetOutput() interface{} {
	return t.output
}

func (t *SyncLibraryTask) Done() {
	if t.option.OnComplete != nil {
		t.option.OnComplete(t)
	}
	t.BaseTask.Done()
}
func (t *SyncLibraryTask) AbortError(err error) error {
	if t.option.OnError != nil {
		t.option.OnError(t, err)
	}
	return t.BaseTask.AbortError(err)
}
func (t *SyncLibraryTask) AbortFileError(imagePath string, err error) {
	fileErr := errors.New(imagePath + " err: " + err.Error())
	t.Logger.Error(fileErr.Error())
	if t.option.OnFileError != nil {
		t.option.OnFileError(t, fileErr)
	}
}
func (t *SyncLibraryTask) OnFileComplete() {
	if t.option.OnFileComplete != nil {
		t.option.OnFileComplete(t)
	}
}
func CreateSyncLibraryTask(option CreateScanTaskOption) (*SyncLibraryTask, error) {
	newTask := &SyncLibraryTask{
		BaseTask: task.NewBaseTask(TypeScanLibrary, "-1", task.GetStatusText(nil, task.StatusDone)),
		option:   option,
	}
	for _, existedTask := range module.Task.Pool.Tasks {
		if scanOutput, ok := existedTask.(*SyncLibraryTask); ok && scanOutput.output.Id == option.LibraryId {
			if existedTask.GetStatus() == task.GetStatusText(nil, task.StatusDone) {
				return scanOutput, nil
			}
			module.Task.Pool.RemoveTaskById(existedTask.GetId())
			break
		}
	}
	library, err := service.GetLibraryWithUser(option.LibraryId, option.UserId)
	if err != nil {
		return nil, err
	}
	if library == nil {
		return nil, errors.New("library not found")
	}
	newTask.library = library
	output := ScanTaskOutput{
		Id:   library.ID,
		Path: library.Path,
		Name: library.Name,
	}
	newTask.output = &output
	newTask.Logger = plugins.DefaultYouLogPlugin.Logger.NewScope("Task").WithFields(
		youlog.Fields{
			"path":      library.Path,
			"libraryId": library.ID,
		})

	module.Task.Pool.AddTask(newTask)
	return newTask, nil
}
