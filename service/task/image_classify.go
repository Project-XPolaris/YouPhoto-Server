package task

import (
	"github.com/allentom/harukap/module/task"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/plugins"
	"os"
)

type ImageClassifyTaskOption struct {
	Uid          string
	ParentTaskId string
	Path         string
	ImageId      uint
}
type ImageClassifyTaskOutput struct {
}
type ImageClassifyTask struct {
	*task.BaseTask
	option     *ImageClassifyTaskOption
	TaskOutput *ImageClassifyTaskOutput
	ImageId    uint
}

func (t *ImageClassifyTask) Stop() error {
	return nil
}

func (t *ImageClassifyTask) Start() error {
	rawFile, err := os.Open(t.option.Path)
	if err != nil {
		return t.AbortError(err)
	}
	if rawFile != nil {
		predictions, _ := plugins.DefaultImageClassifyPlugin.Client.Predict(rawFile)
		savePredictionList := make([]*database.Prediction, 0)
		for _, prediction := range predictions {
			savePredictionList = append(savePredictionList, &database.Prediction{
				ImageId:     t.option.ImageId,
				Label:       prediction.Label,
				Probability: prediction.Prob,
			})
		}
		err = database.Instance.Where("image_id = ?", t.option.ImageId).Delete(&database.Prediction{}).Error
		if err != nil {
			return t.AbortError(err)
		}
		err = database.Instance.Create(&savePredictionList).Error
		if err != nil {
			return t.AbortError(err)
		}

	}
	t.Done()
	return nil
}

func (t *ImageClassifyTask) Output() (interface{}, error) {
	return t.TaskOutput, nil
}

func NewImageClassifyTask(option *ImageClassifyTaskOption) *ImageClassifyTask {
	t := &ImageClassifyTask{
		BaseTask:   task.NewBaseTask(TypeImageClassify, option.Uid, task.GetStatusText(nil, task.StatusRunning)),
		TaskOutput: &ImageClassifyTaskOutput{},
		option:     option,
	}
	t.ParentTaskId = option.ParentTaskId
	return t
}
