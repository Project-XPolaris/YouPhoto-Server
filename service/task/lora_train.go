package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/allentom/harukap/module/task"
	"github.com/projectxpolaris/youphoto/config"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/service/lora"
	"github.com/projectxpolaris/youphoto/service/sdw"
	"os"
	"path/filepath"
	"time"
)

type TrainOptions struct {
	PreprocessConfig *sdw.PreprocessParam    `json:"preprocess_config"`
	Step             int                     `json:"step"`
	TranParam        *lora.TrainConfigValues `json:"train_param"`
}
type LoraTrainTaskOption struct {
	Uid          string
	ParentTaskId string
	LibraryId    uint
	ConfigId     uint
}
type LoraTrainTaskOutput struct {
	TrainId      string  `json:"trainId"`
	Step         int     `json:"step"`
	Epoch        int     `json:"epoch"`
	TotalEpoch   int     `json:"totalEpoch"`
	TotalStep    int     `json:"totalStep"`
	AllStep      int     `json:"allStep"`
	AllTotalStep int     `json:"allTotalStep"`
	Progress     float64 `json:"progress"`
	LibraryName  string  `json:"libraryName"`
	LibraryId    uint    `json:"libraryId"`
}
type LoraTrainTask struct {
	*task.BaseTask
	option     *LoraTrainTaskOption
	TaskOutput *LoraTrainTaskOutput
	OutputPath string
}

func (t *LoraTrainTask) Stop() error {
	return nil
}

func (t *LoraTrainTask) Start() error {
	var library database.Library
	err := database.Instance.First(&library, t.option.LibraryId).Error
	if err != nil {
		return t.AbortError(err)
	}
	var savedConfig database.LoraConfig
	err = database.Instance.First(&savedConfig, t.option.ConfigId).Error
	if err != nil {
		return t.AbortError(err)
	}
	// load config
	var trainConfig TrainOptions
	err = json.Unmarshal([]byte(savedConfig.Config), &trainConfig)
	if err != nil {
		return t.AbortError(err)
	}

	// preprocess
	outPath, err := filepath.Abs(filepath.Join(config.Instance.PreprocessPath, "lora_train", t.Id, "train", fmt.Sprintf("%d_%d", trainConfig.Step, library.ID)))
	if err != nil {
		return t.AbortError(err)
	}
	err = os.MkdirAll(outPath, 0755)
	if err != nil {
		return t.AbortError(err)
	}
	t.OutputPath = outPath
	trainConfig.PreprocessConfig.ProcessSrc = library.Path
	trainConfig.PreprocessConfig.ProcessDst = outPath
	preProcessTask := NewPreprocessTask(&PreprocessTaskOption{
		Uid:          t.option.Uid,
		ParentTaskId: t.option.ParentTaskId,
		LibraryId:    library.ID,
		Param:        trainConfig.PreprocessConfig,
	})
	t.SubTaskList = append(t.SubTaskList, preProcessTask)
	err = task.RunTask(preProcessTask)
	if err != nil {
		return t.AbortError(err)
	}
	// use default config
	trainConfig.TranParam.ModelList = trainConfig.TranParam.PretrainedModelNameOrPath
	trainConfig.TranParam.NetworkModule = "networks.lora"
	// create output path
	modelOutputPath := filepath.Join(config.Instance.ModelOutPath)
	err = os.MkdirAll(modelOutputPath, 0755)
	if err != nil {
		return t.AbortError(err)
	}
	outputPathAbs, err := filepath.Abs(modelOutputPath)
	if err != nil {
		return t.AbortError(err)
	}
	trainConfig.TranParam.OutputDir = outputPathAbs
	trainConfig.TranParam.TrainDataDir = filepath.Dir(outPath)
	// print config in json
	configJson, _ := json.MarshalIndent(trainConfig.TranParam, "", "  ")
	fmt.Println("train config", string(configJson))
	// train

	createdTask, err := lora.DefaultLoraTrainClient.Train(trainConfig.TranParam)
	if err != nil {
		return t.AbortError(err)
	}
	var status *lora.TrainTask
	t.TaskOutput.TrainId = createdTask.Id
	// wait for complete
	for {
		status, err = lora.DefaultLoraTrainClient.FetchTask(createdTask.Id)
		if err != nil {
			return t.AbortError(err)
		}
		t.TaskOutput.Epoch = status.Epoch
		t.TaskOutput.Step = status.Steps
		t.TaskOutput.TotalEpoch = status.TotalEpochs
		t.TaskOutput.TotalStep = status.TotalSteps
		t.TaskOutput.AllTotalStep = status.TotalEpochs * status.TotalSteps
		t.TaskOutput.AllStep = ((status.Epoch - 1) * status.TotalSteps) + status.Steps
		if t.TaskOutput.AllTotalStep != 0 {
			t.TaskOutput.Progress = float64(t.TaskOutput.AllStep) / float64(t.TaskOutput.AllTotalStep)
		}
		if status.Status != "running" {
			break
		}
		<-time.After(1 * time.Second)
	}
	if status.Status == "error" {
		return t.AbortError(errors.New(fmt.Sprintf("train error: %s", status.Error)))
	}
	if status.Status == "interrupted" {
		t.Status = "interrupted"
	}
	t.Done()
	return nil
}

func (t *LoraTrainTask) Output() (interface{}, error) {
	return t.TaskOutput, nil
}

func NewLoraTrainTask(option *LoraTrainTaskOption) (*LoraTrainTask, error) {
	var library database.Library
	err := database.Instance.First(&library, option.LibraryId).Error
	if err != nil {
		return nil, err
	}
	t := &LoraTrainTask{
		BaseTask: task.NewBaseTask(TypeLoraTrain, option.Uid, task.GetStatusText(nil, task.StatusRunning)),
		TaskOutput: &LoraTrainTaskOutput{
			LibraryId:   library.ID,
			LibraryName: library.Name,
		},
		option: option,
	}

	t.ParentTaskId = option.ParentTaskId
	return t, nil
}
