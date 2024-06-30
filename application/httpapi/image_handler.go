package httpapi

import (
	"bytes"
	context2 "context"
	"errors"
	"github.com/allentom/haruka"
	task2 "github.com/allentom/harukap/module/task"
	"github.com/projectxpolaris/youphoto/config"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/plugins"
	"github.com/projectxpolaris/youphoto/service"
	"github.com/projectxpolaris/youphoto/service/task"
	"github.com/projectxpolaris/youphoto/utils"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

var getImageListHandler haruka.RequestHandler = func(context *haruka.Context) {
	queryBuilder := service.ImagesQueryBuilder{
		Page:     context.Param["page"].(int),
		PageSize: context.Param["pageSize"].(int),
	}
	err := context.BindingInput(&queryBuilder)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	md5List := context.GetQueryStrings("md5")
	if len(md5List) > 0 {
		queryBuilder.Md5 = md5List
	}
	if claims, ok := context.Param["claim"]; ok {
		queryBuilder.UserId = claims.(*database.User).ID
	} else {
		AbortError(context, errors.New("need auth"), http.StatusBadRequest)
		return
	}
	imageList, count, err := queryBuilder.Query()
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	data := NewBaseImageTemplateList(imageList)
	MakeListResponse(context, queryBuilder.Page, queryBuilder.PageSize, count, data)
}
var getImageHandler haruka.RequestHandler = func(context *haruka.Context) {
	id, err := context.GetPathParameterAsInt("id")
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	image, err := service.GetImageById(uint(id), "ImageColor", "Prediction", "DeepdanbooruResult", "Tags")
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	context.JSON(haruka.JSON{
		"success": true,
		"data":    NewBaseImageTemplate(image),
	})
}
var getImageThumbnailHandler haruka.RequestHandler = func(context *haruka.Context) {
	id, err := context.GetPathParameterAsInt("id")
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	image, err := service.GetImageById(uint(id))
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	storageKey := filepath.Join(config.Instance.ThumbnailStorePath, image.Thumbnail)
	storage := plugins.GetDefaultStorage()
	buf, err := storage.Get(context2.Background(), utils.DefaultBucket, storageKey)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	data, _ := ioutil.ReadAll(buf)
	http.ServeContent(context.Writer, context.Request, image.Thumbnail, time.Now(), bytes.NewReader(data))
}
var getThumbnailHandler haruka.RequestHandler = func(context *haruka.Context) {
	id := context.GetPathParameterAsString("id")
	storageKey := filepath.Join(config.Instance.ThumbnailStorePath, id)
	storage := plugins.GetDefaultStorage()
	buf, err := storage.Get(context2.Background(), utils.DefaultBucket, storageKey)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	data, _ := ioutil.ReadAll(buf)
	http.ServeContent(context.Writer, context.Request, id, time.Now(), bytes.NewReader(data))
}
var getImageRawHandler haruka.RequestHandler = func(context *haruka.Context) {
	id, err := context.GetPathParameterAsInt("id")
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	isDownload := context.GetQueryString("download")
	image, err := service.GetImageById(uint(id), "Library")
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	if isDownload == "1" {
		context.Writer.Header().Set("Content-Disposition", "attachment; filename="+image.Name)
	}
	http.ServeFile(context.Writer, context.Request, filepath.Join(image.Library.Path, image.Path))
}

var getNearImageHandler haruka.RequestHandler = func(context *haruka.Context) {
	id, err := context.GetPathParameterAsInt("id")
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	maxDistance, err := context.GetQueryInt("maxDistance")
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	images, err := service.GetNearImage(uint(id), maxDistance)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	data := NewNearImageTemplateList(images)
	context.JSON(haruka.JSON{
		"success": true,
		"data":    data,
	})

}

var getColorMatchHandler haruka.RequestHandler = func(context *haruka.Context) {
	var option service.MatchColorOption
	err := context.ParseJson(&option)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	images, err := service.MatchColor(option)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}

	context.JSON(haruka.JSON{
		"success": true,
		"data":    NewColorMatchTemplateList(images),
	})

}

var deepdanbooruHandler haruka.RequestHandler = func(context *haruka.Context) {
	if !plugins.DefaultDeepDanbooruPlugin.Enable {
		AbortError(context, nil, http.StatusForbidden)
		return
	}
	id, err := context.GetQueryInt("id")
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	image, err := service.GetImageById(uint(id), "Library")
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(image.Library.Path, image.Path)
	dbrtask := task.NewDeepdanbooruTask(&task.DeepdanbooruTaskOption{
		Uid:      "-1",
		FullPath: filePath,
		ImageId:  image.ID,
	})
	err = task2.RunTask(dbrtask)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	context.JSON(haruka.JSON{
		"success": true,
		"data":    NewDeepdanbooruTemplateList(dbrtask.Predictions),
	})
}

type UploadImageByBase64RequestBody struct {
	Base64    string `json:"base64"`
	Filename  string `json:"filename"`
	LibraryId int    `json:"libraryId"`
}

var uploadImageByBase64Handler haruka.RequestHandler = func(context *haruka.Context) {
	var requestBody UploadImageByBase64RequestBody
	err := context.ParseJson(&requestBody)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	image, err := task.SaveImageByBase64(requestBody.Base64, requestBody.Filename, uint(requestBody.LibraryId))
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	context.JSON(haruka.JSON{
		"success": true,
		"data":    NewBaseImageTemplate(image),
	})
}

var getImageTaggerHandler haruka.RequestHandler = func(context *haruka.Context) {
	id, err := context.GetPathParameterAsInt("id")
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	model := context.GetQueryString("model")
	result, err := service.TagImageById(uint(id), model)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	context.JSON(haruka.JSON{
		"success": true,
		"data":    NewImageTagTemplateList(result),
	})
}

var getImageTagListHandler haruka.RequestHandler = func(context *haruka.Context) {
	queryBuilder := service.TagQueryBuilder{
		Page:     context.Param["page"].(int),
		PageSize: context.Param["pageSize"].(int),
	}
	err := context.BindingInput(&queryBuilder)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	tagList, count, err := queryBuilder.Query()
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	data := NewImageTagTemplateList(tagList)
	MakeListResponse(context, queryBuilder.Page, queryBuilder.PageSize, count, data)
}

var getImageTaggerModelHandler haruka.RequestHandler = func(context *haruka.Context) {
	models, err := service.GetTaggerList()
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	context.JSON(haruka.JSON{
		"success": true,
		"data":    models,
	})
}

var uploadImageByFileHandler haruka.RequestHandler = func(context *haruka.Context) {
	file, _, err := context.Request.FormFile("file")
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	libraryId, err := context.GetQueryInt("libraryId")
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	var uid uint = 0
	if claims, ok := context.Param["claim"]; ok {
		uid = claims.(*database.User).ID
	} else {
		AbortError(context, errors.New("need auth"), http.StatusBadRequest)
		return
	}
	albumName := context.GetQueryString("albumName")
	albumId, _ := context.GetQueryInt("albumId")

	filename := context.GetQueryString("filename")
	image, err := task.SaveImageByFile(file, filename, uint(libraryId), albumName, uint(albumId), uid)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	context.JSON(haruka.JSON{
		"success": true,
		"data":    NewBaseImageTemplate(image),
	})
}

var deleteImageByIdsHandler haruka.RequestHandler = func(context *haruka.Context) {
	rawIds := context.GetQueryStrings("ids")
	rawDeleteImage := context.GetQueryString("deleteImage")
	ids := make([]uint, 0)
	for _, id := range rawIds {
		idInt, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			AbortError(context, err, http.StatusBadRequest)
			return
		}
		ids = append(ids, uint(idInt))
	}
	err := service.DeleteImageByIds(ids, len(rawDeleteImage) > 0)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	context.JSON(haruka.JSON{
		"success": true,
	})
}
