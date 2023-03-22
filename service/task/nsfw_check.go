package task

import (
	"github.com/allentom/harukap/module/task"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/plugins"
	"os"
)

type NSFWCheckTaskOption struct {
	Uid          string
	ParentTaskId string
	Path         string
	Image        database.Image
}
type NSFWCheckTaskOutput struct {
}
type NSFWCheckTask struct {
	*task.BaseTask
	option     *NSFWCheckTaskOption
	TaskOutput *NSFWCheckTaskOutput
	ImageId    uint
}

func (t *NSFWCheckTask) Stop() error {
	return nil
}

func (t *NSFWCheckTask) Start() error {
	rawFile, err := os.Open(t.option.Path)
	if err != nil {
		return t.AbortError(err)
	}
	if rawFile != nil {
		predictions, _ := plugins.DefaultNSFWCheckPlugin.Client.Predict(rawFile)
		for _, prediction := range predictions {
			switch prediction.Classname {
			case "Sexy":
				t.option.Image.Sexy = prediction.Probability
			case "Neutral":
				t.option.Image.Neutral = prediction.Probability
			case "Drawing":
				t.option.Image.Drawings = prediction.Probability
			case "Hentai":
				t.option.Image.Hentai = prediction.Probability
			case "Porn":
				t.option.Image.Porn = prediction.Probability
			}
		}
	}
	t.Done()
	return nil
}

func (t *NSFWCheckTask) Output() (interface{}, error) {
	return t.TaskOutput, nil
}

func NewNSFWCheckTask(option *NSFWCheckTaskOption) *NSFWCheckTask {
	t := &NSFWCheckTask{
		BaseTask:   task.NewBaseTask(TypeNSFWCheck, option.Uid, task.GetStatusText(nil, task.StatusRunning)),
		TaskOutput: &NSFWCheckTaskOutput{},
		option:     option,
	}
	t.ParentTaskId = option.ParentTaskId
	return t
}
