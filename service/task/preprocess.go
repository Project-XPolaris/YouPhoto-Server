package task

import (
	"fmt"
	"github.com/allentom/harukap/module/task"
	"github.com/projectxpolaris/youphoto/config"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/service/sdw"
	"path/filepath"
)

type PreprocessTaskOption struct {
	Uid          string
	ParentTaskId string
	LibraryId    uint
	Param        *sdw.PreprocessParam
}
type PreprocessTaskOutput struct {
}
type PreprocessTask struct {
	*task.BaseTask
	option     *PreprocessTaskOption
	TaskOutput *PreprocessTaskOutput
	OutputPath string
}

func (t *PreprocessTask) Stop() error {
	return nil
}

func (t *PreprocessTask) Start() error {
	var library database.Library
	err := database.Instance.First(&library, t.option.LibraryId).Error
	if err != nil {
		return t.AbortError(err)

	}
	outPath, err := filepath.Abs(filepath.Join(config.Instance.PreprocessPath, fmt.Sprintf("%d", library.ID)))
	if err != nil {
		return t.AbortError(err)
	}
	t.OutputPath = outPath
	if len(t.option.Param.ProcessSrc) == 0 {
		t.option.Param.ProcessSrc = library.Path
	}
	if len(t.option.Param.ProcessDst) == 0 {
		t.option.Param.ProcessDst = outPath
	}
	err = sdw.DefaultSDWClient.Preprocess(t.option.Param)
	if err != nil {
		return t.AbortError(err)
	}
	t.Done()
	return nil
}

func (t *PreprocessTask) Output() (interface{}, error) {
	return t.TaskOutput, nil
}

func NewPreprocessTask(option *PreprocessTaskOption) *PreprocessTask {
	t := &PreprocessTask{
		BaseTask:   task.NewBaseTask(TypePreprocess, option.Uid, task.GetStatusText(nil, task.StatusRunning)),
		TaskOutput: &PreprocessTaskOutput{},
		option:     option,
	}
	t.ParentTaskId = option.ParentTaskId
	return t
}
