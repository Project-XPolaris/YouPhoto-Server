package httpapi

import (
	"github.com/allentom/haruka"
	"github.com/projectxpolaris/youphoto/service"
	"net/http"
)

type CreateLibraryRequestBody struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

var createLibraryHandler haruka.RequestHandler = func(context *haruka.Context) {
	var requestBody CreateLibraryRequestBody
	err := context.ParseJson(&requestBody)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	library, err := service.CreateLibrary(requestBody.Name, requestBody.Path)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	data := NewBaseLibraryTemplate(library)
	MakeSuccessResponse(data, context)
}

var getLibraryListHandler haruka.RequestHandler = func(context *haruka.Context) {
	queryBuilder := service.LibraryQueryBuilder{
		Page:     context.Param["page"].(int),
		PageSize: context.Param["pageSize"].(int),
	}
	err := context.BindingInput(&queryBuilder)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	libraryList, count, err := queryBuilder.Query()
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	data := NewBaseLibraryTemplateList(libraryList)
	MakeListResponse(context, queryBuilder.Page, queryBuilder.PageSize, count, data)
}

var scanLibraryHandler haruka.RequestHandler = func(context *haruka.Context) {
	id, err := context.GetPathParameterAsInt("id")
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	option := service.CreateScanTaskOption{
		LibraryId: uint(id),
	}
	_, err = service.CreateSyncLibraryTask(option)
	MakeSuccessResponse(nil, context)
}

var removeLibraryHandler haruka.RequestHandler = func(context *haruka.Context) {
	id, err := context.GetPathParameterAsInt("id")
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	option := service.RemoveLibraryTaskOption{
		LibraryId: uint(id),
	}
	_, err = service.CreateRemoveLibraryTask(option)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	MakeSuccessResponse(nil, context)
}
