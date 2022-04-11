package service

import (
	"errors"
	"github.com/project-xpolaris/youplustoolkit/youlog"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/utils"
	"gorm.io/gorm"
	"path/filepath"
)

type SyncLibraryTask struct {
	BaseTask
	option CreateScanTaskOption
	output *ScanTaskOutput
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
	OnFileComplete func(task Task)
	OnFileError    func(task Task, err error)
	OnError        func(task Task, err error)
	OnComplete     func(task Task)
}

func (t *SyncLibraryTask) GetOutput() interface{} {
	return t.output
}

func (t *SyncLibraryTask) Done() {
	t.UpdateDoneStatus()
	if t.option.OnComplete != nil {
		t.option.OnComplete(t)
	}
}
func (t *SyncLibraryTask) AbortError(err error) {
	t.BaseTask.AbortError(err)
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
func CreateSyncLibraryTask(option CreateScanTaskOption) (Task, error) {
	task := &SyncLibraryTask{
		BaseTask: NewBaseTask(TaskTypeScanLibrary),
		option:   option,
	}
	for _, task := range DefaultTaskPool.Tasks {
		if scanOutput, ok := task.GetOutput().(*ScanTaskOutput); ok && scanOutput.Id == option.LibraryId {
			if task.GetStatus() == TaskStatusRunning {
				return task, nil
			}
			DefaultTaskPool.RemoveTaskById(task.GetId())
			break
		}
	}
	var library database.Library
	err := database.Instance.Find(&library, option.LibraryId).Error
	if err != nil {
		return nil, err
	}
	output := ScanTaskOutput{
		Id:   library.ID,
		Path: library.Path,
		Name: library.Name,
	}
	task.output = &output
	task.Logger = task.Logger.WithFields(
		youlog.Fields{
			"path":      library.Path,
			"libraryId": library.ID,
		})
	go func() {
		//database.Instance.Unscoped().Model(&database.Image{}).Where("library_id = ?", library.ID).Delete(database.Image{})
		// sync remove missing images
		task.Logger.Info("start sync remove missing images")
		var existImageCount int64
		err := database.Instance.Model(database.Image{}).
			Where("library_id = ?", library.ID).Count(&existImageCount).Error
		if err != nil {
			task.AbortError(err)
			return
		}
		for idx := 0; idx < int(existImageCount); idx += 20 {
			var images []database.Image
			err = database.Instance.Model(database.Image{}).
				Where("library_id = ?", library.ID).
				Offset(idx).
				Limit(20).
				Find(&images).Error
			if err != nil {
				task.AbortError(err)
				return
			}
			err = database.Instance.Transaction(func(tx *gorm.DB) error {
				for _, image := range images {
					if !utils.CheckFileExist(filepath.Join(library.Path, image.Path)) {
						err := tx.Unscoped().Model(&database.Image{}).Where("id = ?", image.ID).Delete(database.Image{}).Error
						if err != nil {
							return err
						}
					}
				}
				return nil
			})
			if err != nil {
				task.AbortError(err)
				return
			}
		}

		//count total

		scanner := NewImageScanner(library.Path)
		idx := 0
		scanner.OnHit = func(path string) {
			output.Total += 1
		}
		err = scanner.Scan()
		if err != nil {
			task.AbortError(err)
			return
		}
		scanner.OnHit = func(path string) {
			output.Current += int64(idx + 1)
			output.CurrentPath = path
			output.CurrentName = filepath.Base(path)
			imagePath, err := filepath.Rel(library.Path, path)
			if err != nil {
				task.AbortFileError(path, err)
				return
			}
			_, err = CreateImage(imagePath, library.ID, path)
			if err != nil {
				if err != nil {
					task.AbortFileError(path, err)
					return
				}
			} else {
				task.OnFileComplete()
			}
		}
		task.Logger.Info("start scan library")
		err = scanner.Scan()
		if err != nil {
			task.AbortError(err)
			return
		}
		task.Done()

	}()
	DefaultTaskPool.AddTask(task)
	return task, nil
}
