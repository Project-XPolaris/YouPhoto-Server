package task

import (
	"github.com/allentom/harukap/module/task"
	"github.com/projectxpolaris/youphoto/service"
)

type TaggerTaskOption struct {
	Uid          string
	ParentTaskId string
	FullPath     string
	ImageId      uint
	TaggerModel  string
}
type TaggerTaskOutput struct {
}
type TaggerTask struct {
	*task.BaseTask
	option     *TaggerTaskOption
	TaskOutput *TaggerTaskOutput
	thumbnail  string
}

func (t *TaggerTask) Stop() error {
	return nil
}

func (t *TaggerTask) Start() error {
	_, err := service.TagImageById(t.option.ImageId, t.option.TaggerModel, 0.7)
	if err != nil {
		return t.AbortError(err)
	}
	t.Done()
	return nil
}

func (t *TaggerTask) Output() (interface{}, error) {
	return t.TaskOutput, nil
}

func NewTaggerTask(option *TaggerTaskOption) *TaggerTask {
	t := &TaggerTask{
		BaseTask:   task.NewBaseTask(TypeTagger, option.Uid, task.GetStatusText(nil, task.StatusRunning)),
		TaskOutput: &TaggerTaskOutput{},
		option:     option,
	}
	t.ParentTaskId = option.ParentTaskId
	return t
}
