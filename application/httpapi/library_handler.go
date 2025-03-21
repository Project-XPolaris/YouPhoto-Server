package httpapi

import (
	"errors"
	"net/http"

	"github.com/allentom/haruka"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/module"
	"github.com/projectxpolaris/youphoto/service"
	"github.com/projectxpolaris/youphoto/service/task"
)

// CreateLibraryRequestBody 创建图库请求的数据结构
type CreateLibraryRequestBody struct {
	Name    string `json:"name"`    // 图库名称
	Path    string `json:"path"`    // 图库路径
	Private bool   `json:"private"` // 是否私有
}

// createLibraryHandler 创建新的图库处理函数
var createLibraryHandler haruka.RequestHandler = func(context *haruka.Context) {
	// 解析请求体
	var requestBody CreateLibraryRequestBody
	err := context.ParseJson(&requestBody)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	// 获取用户ID
	var uid uint = 0
	if claims, ok := context.Param["claim"]; ok {
		uid = claims.(*database.User).ID
	} else {
		AbortError(context, errors.New("need auth"), http.StatusBadRequest)
		return
	}
	// 创建图库
	library, err := service.CreateLibrary(requestBody.Name, requestBody.Path, uid, !requestBody.Private)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	data := NewBaseLibraryTemplate(library)
	MakeSuccessResponse(data, context)
}

// getLibraryListHandler 获取图库列表处理函数
var getLibraryListHandler haruka.RequestHandler = func(context *haruka.Context) {
	// 构建查询参数
	queryBuilder := service.LibraryQueryBuilder{
		Page:     context.Param["page"].(int),
		PageSize: context.Param["pageSize"].(int),
	}
	err := context.BindingInput(&queryBuilder)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	// 获取用户ID
	var userId uint = 0
	if claims, ok := context.Param["claim"]; ok {
		userId = claims.(*database.User).ID
	}
	queryBuilder.UserId = userId
	// 执行查询
	libraryList, count, err := queryBuilder.Query()
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	data := NewBaseLibraryTemplateList(libraryList)
	MakeListResponse(context, queryBuilder.Page, queryBuilder.PageSize, count, data)
}

// scanLibraryHandler 扫描图库处理函数
var scanLibraryHandler haruka.RequestHandler = func(context *haruka.Context) {
	// 获取图库ID
	id, err := context.GetPathParameterAsInt("id")
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	// 解析处理选项
	precessOption := task.ProcessImageOption{}
	err = context.ParseJson(&precessOption)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}

	// 创建扫描任务选项
	option := task.CreateScanTaskOption{
		LibraryId:     uint(id),
		ProcessOption: &precessOption,
	}
	// 获取用户ID
	if claims, ok := context.Param["claim"]; ok {
		option.UserId = claims.(*database.User).ID
	} else {
		AbortError(context, errors.New("need auth"), http.StatusBadRequest)
		return
	}
	// 创建并启动扫描任务
	task, err := task.CreateSyncLibraryTask(option)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	go task.Start()
	data, err := module.Task.SerializerTemplate(task)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	MakeSuccessResponse(data, context)
}

// removeLibraryHandler 删除图库处理函数
var removeLibraryHandler haruka.RequestHandler = func(context *haruka.Context) {
	// 获取图库ID
	id, err := context.GetPathParameterAsInt("id")
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	// 创建删除任务选项
	option := task.RemoveLibraryTaskOption{
		LibraryId: uint(id),
	}
	// 创建并启动删除任务
	task, err := task.CreateRemoveLibraryTask(option)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	go task.Start()
	data, err := module.Task.SerializerTemplate(task)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	MakeSuccessResponse(data, context)
}
