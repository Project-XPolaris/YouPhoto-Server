package httpapi

import (
	"errors"
	"github.com/allentom/haruka"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/module"
	"github.com/projectxpolaris/youphoto/service"
	"github.com/projectxpolaris/youphoto/service/sdw"
	"github.com/projectxpolaris/youphoto/service/task"
	"net/http"
)

func checkClientAndResult(c *haruka.Context) bool {
	if sdw.DefaultSDWClient == nil {
		AbortError(c, errors.New("sdw client not init"), 404)
		return false
	}
	return true
}

var getModelsHandler = func(c *haruka.Context) {
	if !checkClientAndResult(c) {
		return
	}
	models, err := sdw.DefaultSDWClient.GetModels()
	if err != nil {
		AbortError(c, err, http.StatusInternalServerError)
		return
	}
	MakeSuccessResponse(models, c)
}

type SwitchModelRequestBody struct {
	ModelName string `json:"name"`
}

var switchModelHandler = func(c *haruka.Context) {
	if !checkClientAndResult(c) {
		return
	}
	var requestBody SwitchModelRequestBody
	err := c.ParseJson(&requestBody)
	if err != nil {
		AbortError(c, err, http.StatusBadRequest)
		return
	}
	err = sdw.DefaultSDWClient.UpdateOption(map[string]interface{}{
		"sd_model_checkpoint": requestBody.ModelName,
	})
	if err != nil {
		AbortError(c, err, http.StatusInternalServerError)
		return
	}
	MakeSuccessResponse(nil, c)
}

var sdwInfoHandler = func(c *haruka.Context) {
	if !checkClientAndResult(c) {
		return
	}
	info, err := sdw.DefaultSDWClient.GetOption()
	if err != nil {
		AbortError(c, err, http.StatusInternalServerError)
		return
	}
	MakeSuccessResponse(info, c)
}

var text2ImageHaldler = func(c *haruka.Context) {
	if !checkClientAndResult(c) {
		return
	}
	requestBody := sdw.TextToImageParam{}
	err := c.ParseJson(&requestBody)
	if err != nil {
		AbortError(c, err, http.StatusBadRequest)
		return
	}
	result, err := sdw.DefaultSDWClient.TextToImage(&requestBody)
	if err != nil {
		AbortError(c, err, http.StatusInternalServerError)
		return
	}
	MakeSuccessResponse(result, c)
}
var getSamplerListHandler = func(c *haruka.Context) {
	if !checkClientAndResult(c) {
		return
	}
	list, err := sdw.DefaultSDWClient.GetSamplers()
	if err != nil {
		AbortError(c, err, http.StatusInternalServerError)
		return
	}
	MakeSuccessResponse(list, c)
}

var getUpscalerListHandler = func(c *haruka.Context) {
	if !checkClientAndResult(c) {
		return
	}
	list, err := sdw.DefaultSDWClient.GetUpscaler()
	if err != nil {
		AbortError(c, err, http.StatusInternalServerError)
		return
	}
	MakeSuccessResponse(list, c)
}

var getProgressHandler = func(c *haruka.Context) {
	if !checkClientAndResult(c) {
		return
	}
	progress, err := sdw.DefaultSDWClient.GetProgress()
	if err != nil {
		AbortError(c, err, http.StatusInternalServerError)
		return
	}
	MakeSuccessResponse(progress, c)
}

type ConfigRequestBody struct {
	Name   string `json:"name"`
	Config string `json:"config"`
}

var saveConfigHandler = func(c *haruka.Context) {
	if !checkClientAndResult(c) {
		return
	}
	user := c.Param["claim"].(*database.User)
	var requestBody ConfigRequestBody
	err := c.ParseJson(&requestBody)
	if err != nil {
		AbortError(c, err, http.StatusBadRequest)
		return
	}
	config, err := service.SaveSDWConfig(requestBody.Name, requestBody.Config, user.ID)
	if err != nil {
		AbortError(c, err, http.StatusBadRequest)
		return
	}
	data := NewSDWConfigTemplate(config)
	MakeSuccessResponse(data, c)
}

var getSDWConfigListHandler = func(c *haruka.Context) {
	user := c.Param["claim"].(*database.User)
	configs, err := service.GetSDWConfigList(user.ID)
	if err != nil {
		AbortError(c, err, http.StatusBadRequest)
		return
	}
	data := NewSDWConfigTemplateList(configs)
	MakeSuccessResponse(data, c)
}

var deleteSDWConfigHandler = func(c *haruka.Context) {
	id, err := c.GetQueryInt("id")
	if err != nil {
		AbortError(c, err, http.StatusBadRequest)
		return
	}
	user := c.Param["claim"].(*database.User)
	err = service.DeleteSDWConfig(uint(id), user.ID)
	if err != nil {
		AbortError(c, err, http.StatusBadRequest)
		return
	}
	MakeSuccessResponse(nil, c)
}

var intrruptHandler = func(c *haruka.Context) {
	if !checkClientAndResult(c) {
		return
	}
	err := sdw.DefaultSDWClient.Interrupt()
	if err != nil {
		AbortError(c, err, http.StatusInternalServerError)
		return
	}
	MakeSuccessResponse(nil, c)
}

var skipHandler = func(c *haruka.Context) {
	if !checkClientAndResult(c) {
		return
	}
	err := sdw.DefaultSDWClient.Skip()
	if err != nil {
		AbortError(c, err, http.StatusInternalServerError)
		return
	}
	MakeSuccessResponse(nil, c)
}

var newPreprocessHandler = func(c *haruka.Context) {
	if !checkClientAndResult(c) {
		return
	}
	libraryId, err := c.GetQueryInt("id")
	if err != nil {
		AbortError(c, err, http.StatusBadRequest)
		return
	}

	param := sdw.PreprocessParam{}
	err = c.ParseJson(&param)
	if err != nil {
		AbortError(c, err, http.StatusBadRequest)
		return
	}
	newTask := task.NewPreprocessTask(&task.PreprocessTaskOption{
		Uid:       "-1",
		LibraryId: uint(libraryId),
		Param:     &param,
	})
	module.Task.Pool.AddTask(newTask)
	go newTask.Start()
	if err != nil {
		AbortError(c, err, http.StatusInternalServerError)
		return
	}
	MakeSuccessResponse(nil, c)
}
