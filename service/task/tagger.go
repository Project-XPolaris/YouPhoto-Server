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
	_, err := service.TagImageById(t.option.ImageId)
	if err != nil {
		return t.AbortError(err)
	}
	//savePredictions := make([]*database.TaggerResult, 0)
	//for _, prediction := range result {
	//	if prediction.Prob > 0.5 {
	//		savePredictions = append(savePredictions, &database.TaggerResult{
	//			ImageId: t.option.ImageId,
	//			Tag:     prediction.Tag,
	//			Prob:    prediction.Prob,
	//		})
	//	}
	//}
	//err = database.Instance.Where("image_id = ?", t.option.ImageId).Delete(&database.TaggerResult{}).Error
	//if err != nil {
	//	return t.AbortError(err)
	//}
	//err = database.Instance.Create(&savePredictions).Error
	//if err != nil {
	//	return t.AbortError(err)
	//}
	//t.Predictions = savePredictions
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
