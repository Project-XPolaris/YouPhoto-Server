package task

import (
	"github.com/allentom/harukap/module/task"
	"github.com/projectxpolaris/youphoto/service"
)

type ScanImageFileTaskOption struct {
	Uid          string
	Path         string
	ParentTaskId string
}
type ScanImageFileTaskOutput struct {
}
type ScanImageFileTask struct {
	*task.BaseTask
	option     *ScanImageFileTaskOption
	TaskOutput *ScanImageFileTaskOutput
	pathList   []string
}

func (t *ScanImageFileTask) Stop() error {
	return nil
}

func (t *ScanImageFileTask) Start() error {
	scanner := service.NewImageScanner(t.option.Path)
	scanner.OnHit = func(path string) error {
		t.pathList = append(t.pathList, path)
		return nil
	}
	err := scanner.Scan()
	if err != nil {
		return t.AbortError(err)
	}
	t.Done()
	return nil
}

func (t *ScanImageFileTask) Output() (interface{}, error) {
	return t.TaskOutput, nil
}

func NewScanImageFileTask(option *ScanImageFileTaskOption) *ScanImageFileTask {
	t := &ScanImageFileTask{
		BaseTask:   task.NewBaseTask(TypeScanImageFile, option.Uid, task.GetStatusText(nil, task.StatusRunning)),
		TaskOutput: &ScanImageFileTaskOutput{},
		option:     option,
		pathList:   make([]string, 0),
	}
	t.ParentTaskId = option.ParentTaskId
	return t
}
