package task

import (
	"github.com/allentom/harukap/module/task"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/plugins"
)

type DeepdanbooruTaskOption struct {
	Uid          string
	ParentTaskId string
	FullPath     string
	ImageId      uint
}
type DeepdanbooruTaskOutput struct {
}
type DeepdanbooruTask struct {
	*task.BaseTask
	option      *DeepdanbooruTaskOption
	TaskOutput  *DeepdanbooruTaskOutput
	thumbnail   string
	Predictions []*database.DeepdanbooruResult
}

func (t *DeepdanbooruTask) Stop() error {
	return nil
}

func (t *DeepdanbooruTask) Start() error {
	request := plugins.DefaultDeepdanbooruLauncher.Launch(t.option.FullPath)
	result, err := request.Wait()
	if err != nil {
		return t.AbortError(err)
	}
	savePredictions := make([]*database.DeepdanbooruResult, 0)
	for _, prediction := range result {
		if prediction.Prob > 0.5 {
			savePredictions = append(savePredictions, &database.DeepdanbooruResult{
				ImageId: t.option.ImageId,
				Tag:     prediction.Tag,
				Prob:    prediction.Prob,
			})
		}
	}
	err = database.Instance.Where("image_id = ?", t.option.ImageId).Delete(&database.DeepdanbooruResult{}).Error
	if err != nil {
		return t.AbortError(err)
	}
	err = database.Instance.Create(&savePredictions).Error
	if err != nil {
		return t.AbortError(err)
	}
	t.Predictions = savePredictions
	t.Done()
	return nil
}

func (t *DeepdanbooruTask) Output() (interface{}, error) {
	return t.TaskOutput, nil
}

func NewDeepdanbooruTask(option *DeepdanbooruTaskOption) *DeepdanbooruTask {
	t := &DeepdanbooruTask{
		BaseTask:   task.NewBaseTask(TypeDeepdanbooru, option.Uid, task.GetStatusText(nil, task.StatusRunning)),
		TaskOutput: &DeepdanbooruTaskOutput{},
		option:     option,
	}
	t.ParentTaskId = option.ParentTaskId
	return t
}
