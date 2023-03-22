package task

import (
	"github.com/allentom/harukap/module/task"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/plugins"
	"os"
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
	option     *DeepdanbooruTaskOption
	TaskOutput *DeepdanbooruTaskOutput
	thumbnail  string
}

func (t *DeepdanbooruTask) Stop() error {
	return nil
}

func (t *DeepdanbooruTask) Start() error {
	rawFile, err := os.Open(t.option.FullPath)
	result, err := plugins.DefaultDeepDanbooruPlugin.Client.Tagging(rawFile)
	if err != nil {
		return t.AbortError(err)
	}
	savePredictions := make([]database.DeepdanbooruResult, 0)
	for _, prediction := range result {
		if prediction.Prob > 0.5 {
			savePredictions = append(savePredictions, database.DeepdanbooruResult{
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
