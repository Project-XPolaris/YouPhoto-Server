package task

import (
	"github.com/allentom/harukap/module/task"
	"github.com/projectxpolaris/youphoto/service"
	"image"
)

type ReadImageTaskOption struct {
	Uid          string
	ParentTaskId string
	Path         string
}
type ReadImageTaskOutput struct {
}
type ReadImageTask struct {
	*task.BaseTask
	option     *ReadImageTaskOption
	TaskOutput *ReadImageTaskOutput
	Image      image.Image
}

func (t *ReadImageTask) Stop() error {
	return nil
}

func (t *ReadImageTask) Start() error {
	var err error
	t.Image, err = service.GetImageFromFilePath(t.option.Path)
	if err != nil {
		return t.AbortError(err)
	}
	t.Done()
	return nil
}

func (t *ReadImageTask) Output() (interface{}, error) {
	return t.TaskOutput, nil
}

func NewReadImageTask(option *ReadImageTaskOption) *ReadImageTask {
	t := &ReadImageTask{
		BaseTask:   task.NewBaseTask(TypeReadImage, option.Uid, task.GetStatusText(nil, task.StatusRunning)),
		TaskOutput: &ReadImageTaskOutput{},
		option:     option,
	}
	t.ParentTaskId = option.ParentTaskId
	return t
}
