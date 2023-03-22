package task

import (
	"github.com/allentom/harukap/module/task"
	"github.com/projectxpolaris/youphoto/service"
)

type GenerateThumbnailTaskOption struct {
	Uid          string
	ParentTaskId string
	FullPath     string
}
type GenerateThumbnailTaskOutput struct {
}
type GenerateThumbnailTask struct {
	*task.BaseTask
	option     *GenerateThumbnailTaskOption
	TaskOutput *GenerateThumbnailTaskOutput
	thumbnail  string
}

func (t *GenerateThumbnailTask) Stop() error {
	return nil
}

func (t *GenerateThumbnailTask) Start() error {
	var err error
	t.thumbnail, err = service.GenerateThumbnail(t.option.FullPath)
	if err != nil {
		return t.AbortError(err)
	}
	t.Done()
	return nil
}

func (t *GenerateThumbnailTask) Output() (interface{}, error) {
	return t.TaskOutput, nil
}

func NewGenerateThumbnailTask(option *GenerateThumbnailTaskOption) *GenerateThumbnailTask {
	t := &GenerateThumbnailTask{
		BaseTask:   task.NewBaseTask(TypeGenerateThumbnail, option.Uid, task.GetStatusText(nil, task.StatusRunning)),
		TaskOutput: &GenerateThumbnailTaskOutput{},
		option:     option,
	}
	t.ParentTaskId = option.ParentTaskId
	return t
}
