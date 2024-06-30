package httpapi

import (
	"errors"
	"github.com/allentom/haruka"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/service"
	"net/http"
)

type CreateAlbumRequestBody struct {
	Name string `json:"name"`
}

var createAlbumHandler haruka.RequestHandler = func(context *haruka.Context) {
	var requestBody CreateAlbumRequestBody
	err := context.ParseJson(&requestBody)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	var uid = ""
	if claims, ok := context.Param["claim"]; ok {
		uid = claims.(*database.User).Uid
	} else {
		AbortError(context, errors.New("need auth to create album"), http.StatusBadRequest)
		return
	}
	album, err := service.CreateAlbum(requestBody.Name, uid)
	MakeSuccessResponse(NewBaseAlbumTemplate(album), context)
}

var addImageToAlbumHandler haruka.RequestHandler = func(context *haruka.Context) {
	var requestBody struct {
		ImageIds []uint `json:"imageIds"`
	}
	err := context.ParseJson(&requestBody)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	var uid = ""
	if claims, ok := context.Param["claim"]; ok {
		uid = claims.(*database.User).Uid
	} else {
		AbortError(context, errors.New("need auth to create album"), http.StatusBadRequest)
		return
	}
	id, err := context.GetPathParameterAsInt("id")
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	err = service.AddImageToAlbum(uint(id), uid, requestBody.ImageIds...)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	MakeSuccessResponse(nil, context)
}

var removeImageFromAlbumHandler haruka.RequestHandler = func(context *haruka.Context) {
	var requestBody struct {
		ImageIds []uint `json:"imageIds"`
	}
	err := context.ParseJson(&requestBody)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	var uid = ""
	if claims, ok := context.Param["claim"]; ok {
		uid = claims.(*database.User).Uid
	} else {
		AbortError(context, errors.New("need auth to create album"), http.StatusBadRequest)
		return
	}
	id, err := context.GetPathParameterAsInt("id")
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	err = service.RemoveImageFromAlbum(uint(id), uid, requestBody.ImageIds...)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	MakeSuccessResponse(nil, context)
}

var getAlbumListHandler haruka.RequestHandler = func(context *haruka.Context) {
	queryBuilder := service.AlbumQueryBuilder{
		Page:     context.Param["page"].(int),
		PageSize: context.Param["pageSize"].(int),
	}
	err := context.BindingInput(&queryBuilder)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	if claims, ok := context.Param["claim"]; ok {
		queryBuilder.Uid = claims.(*database.User).Uid
	}
	albumList, count, err := queryBuilder.Query()
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	data := NewBaseAlbumTemplateList(albumList)
	MakeListResponse(context, queryBuilder.Page, queryBuilder.PageSize, count, data)
}

var removeAlbumHandler haruka.RequestHandler = func(context *haruka.Context) {
	id, err := context.GetPathParameterAsInt("id")
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	deleteImage := context.GetQueryString("deleteImage")
	if id == 0 {
		AbortError(context, errors.New("need albumId"), http.StatusBadRequest)
		return
	}
	if claims, ok := context.Param["claim"]; ok {
		uid := claims.(*database.User).Uid
		err := service.RemoveAlbum(uint(id), uid, len(deleteImage) > 0)
		if err != nil {
			AbortError(context, err, http.StatusInternalServerError)
			return
		}
		MakeSuccessResponse(nil, context)
	} else {
		AbortError(context, errors.New("need auth"), http.StatusBadRequest)
		return
	}
}

var getAlbumDetailHandler haruka.RequestHandler = func(context *haruka.Context) {
	id, err := context.GetPathParameterAsInt("id")
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	if id == 0 {
		AbortError(context, errors.New("need albumId"), http.StatusBadRequest)
		return
	}
	if claims, ok := context.Param["claim"]; ok {
		uid := claims.(*database.User).Uid
		album, err := service.GetAlbumById(uint(id), uid)
		if err != nil {
			AbortError(context, err, http.StatusInternalServerError)
			return
		}
		MakeSuccessResponse(NewBaseAlbumTemplate(album), context)
	} else {
		AbortError(context, errors.New("need auth"), http.StatusBadRequest)
		return
	}
}
