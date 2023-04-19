package httpapi

import (
	"errors"
	"github.com/allentom/haruka"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/module"
	"github.com/projectxpolaris/youphoto/service"
	"github.com/projectxpolaris/youphoto/service/lora"
	"github.com/projectxpolaris/youphoto/service/sdw"
	"github.com/projectxpolaris/youphoto/service/task"
	"net/http"
)

func checkLoraClientAndResult(c *haruka.Context) bool {
	if sdw.DefaultSDWClient == nil {
		AbortError(c, errors.New("sdw client not init"), 404)
		return false
	}
	return true
}

var saveLoraConfigHandler = func(c *haruka.Context) {
	user := c.Param["claim"].(*database.User)
	var requestBody ConfigRequestBody
	err := c.ParseJson(&requestBody)
	if err != nil {
		AbortError(c, err, http.StatusBadRequest)
		return
	}
	config, err := service.SaveLoraConfig(requestBody.Name, requestBody.Config, user.ID)
	if err != nil {
		AbortError(c, err, http.StatusBadRequest)
		return
	}
	data := NewLoraConfigTemplate(config)
	MakeSuccessResponse(data, c)
}

var getLoraConfigListHandler = func(c *haruka.Context) {
	user := c.Param["claim"].(*database.User)
	configs, err := service.GetLoraConfigList(user.ID)
	if err != nil {
		AbortError(c, err, http.StatusBadRequest)
		return
	}
	data := NewLoraConfigTemplateList(configs)
	MakeSuccessResponse(data, c)
}

var deleteLoraConfigHandler = func(c *haruka.Context) {
	id, err := c.GetQueryInt("id")
	if err != nil {
		AbortError(c, err, http.StatusBadRequest)
		return
	}
	user := c.Param["claim"].(*database.User)
	err = service.DeleteLoraConfig(uint(id), user.ID)
	if err != nil {
		AbortError(c, err, http.StatusBadRequest)
		return
	}
	MakeSuccessResponse(nil, c)
}

var loraTrainHandler = func(c *haruka.Context) {
	if !checkLoraClientAndResult(c) {
		return
	}
	configId, err := c.GetQueryInt("configid")
	if err != nil {
		AbortError(c, err, http.StatusBadRequest)
		return
	}
	libraryId, err := c.GetQueryInt("libraryid")
	if err != nil {
		AbortError(c, err, http.StatusBadRequest)
		return
	}
	trainTask, err := task.NewLoraTrainTask(&task.LoraTrainTaskOption{
		Uid:       "-1",
		LibraryId: uint(libraryId),
		ConfigId:  uint(configId),
	})
	if err != nil {
		AbortError(c, err, http.StatusBadRequest)
		return
	}
	module.Task.Pool.AddTask(trainTask)
	go trainTask.Start()
	MakeSuccessResponse(nil, c)
}

var loraInterruptHandler = func(c *haruka.Context) {
	taskId := c.GetQueryString("id")
	trainTask := module.Task.Pool.GetTaskById(taskId)
	if trainTask == nil {
		AbortError(c, errors.New("task not found"), http.StatusBadRequest)
		return
	}
	if loraTrainTask, ok := trainTask.(*task.LoraTrainTask); ok {
		trainId := loraTrainTask.TaskOutput.TrainId
		if trainId != "" {
			err := lora.DefaultLoraTrainClient.InterruptTask(trainId)
			if err != nil {
				AbortError(c, err, http.StatusBadRequest)
				return
			}
			MakeSuccessResponse(nil, c)
		}

	} else {
		AbortError(c, errors.New("task not found"), http.StatusBadRequest)
		return
	}
}
