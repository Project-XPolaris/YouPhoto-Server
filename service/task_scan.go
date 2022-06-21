package service

import (
	"errors"
	"github.com/allentom/harukap/module/task"
	"github.com/project-xpolaris/youplustoolkit/youlog"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/module"
	"github.com/projectxpolaris/youphoto/plugins"
	"github.com/projectxpolaris/youphoto/utils"
	"gorm.io/gorm"
	"path/filepath"
)

type SyncLibraryTask struct {
	*task.BaseTask
	option  CreateScanTaskOption
	output  *ScanTaskOutput
	library *database.Library
	Logger  *youlog.Scope
}

func (t *SyncLibraryTask) Stop() error {
	return nil
}

func (t *SyncLibraryTask) Start() error {
	t.Logger.Info("start sync remove missing images")
	var existImageCount int64
	err := database.Instance.Model(database.Image{}).
		Where("library_id = ?", t.library.ID).Count(&existImageCount).Error
	if err != nil {
		t.AbortError(err)
		return err
	}
	for idx := 0; idx < int(existImageCount); idx += 20 {
		var images []database.Image
		err = database.Instance.Model(database.Image{}).
			Where("library_id = ?", t.library.ID).
			Offset(idx).
			Limit(20).
			Find(&images).Error
		if err != nil {
			t.AbortError(err)
			return err
		}
		err = database.Instance.Transaction(func(tx *gorm.DB) error {
			for _, image := range images {
				if !utils.CheckFileExist(filepath.Join(t.library.Path, image.Path)) {
					err := tx.Unscoped().Model(&database.Image{}).Where("id = ?", image.ID).Delete(database.Image{}).Error
					if err != nil {
						return err
					}
				}
			}
			return nil
		})
		if err != nil {
			t.AbortError(err)
			return err
		}
	}

	//count total

	scanner := NewImageScanner(t.library.Path)
	idx := 0
	scanner.OnHit = func(path string) {
		t.output.Total += 1
	}
	err = scanner.Scan()
	if err != nil {
		t.AbortError(err)
		return err
	}
	scanner.OnHit = func(path string) {
		t.output.Current += int64(idx + 1)
		t.output.CurrentPath = path
		t.output.CurrentName = filepath.Base(path)
		imagePath, err := filepath.Rel(t.library.Path, path)
		if err != nil {
			t.AbortFileError(path, err)
			return
		}
		_, err = CreateImage(imagePath, t.library.ID, path)
		if err != nil {
			if err != nil {
				t.AbortFileError(path, err)
				return
			}
		} else {
			t.OnFileComplete()
		}
	}
	t.Logger.Info("start scan library")
	err = scanner.Scan()
	if err != nil {
		t.AbortError(err)
		return err
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
}

func (t *SyncLibraryTask) GetOutput() interface{} {
	return t.output
}

func (t *SyncLibraryTask) Done() {
	t.Status = TaskStatusDone
	if t.option.OnComplete != nil {
		t.option.OnComplete(t)
	}
}
func (t *SyncLibraryTask) AbortError(err error) {
	t.BaseTask.Err = err
	if t.option.OnError != nil {
		t.option.OnError(t, err)
	}
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
	task := &SyncLibraryTask{
		BaseTask: task.NewBaseTask(TaskTypeScanLibrary, "-1", TaskStatusRunning),
		option:   option,
	}
	for _, task := range module.Task.Pool.Tasks {
		if scanOutput, ok := task.(*SyncLibraryTask); ok && scanOutput.output.Id == option.LibraryId {
			if task.GetStatus() == TaskStatusRunning {
				return scanOutput, nil
			}
			module.Task.Pool.RemoveTaskById(task.GetId())
			break
		}
	}
	library, err := GetLibraryWithUser(option.LibraryId, option.UserId)
	if err != nil {
		return nil, err
	}
	if library == nil {
		return nil, errors.New("library not found")
	}
	task.library = library
	output := ScanTaskOutput{
		Id:   library.ID,
		Path: library.Path,
		Name: library.Name,
	}
	task.output = &output
	task.Logger = plugins.DefaultYouLogPlugin.Logger.NewScope("Task").WithFields(
		youlog.Fields{
			"path":      library.Path,
			"libraryId": library.ID,
		})

	module.Task.Pool.AddTask(task)
	return task, nil
}
