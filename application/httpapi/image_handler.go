package httpapi

import (
	"bytes"
	context2 "context"
	"github.com/allentom/haruka"
	"github.com/projectxpolaris/youphoto/config"
	"github.com/projectxpolaris/youphoto/plugins"
	"github.com/projectxpolaris/youphoto/service"
	"github.com/projectxpolaris/youphoto/utils"
	"io/ioutil"
	"net/http"
	"path/filepath"
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
	imageList, count, err := queryBuilder.Query()
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	data := NewBaseImageTemplateList(imageList)
	MakeListResponse(context, queryBuilder.Page, queryBuilder.PageSize, count, data)
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
	//http.ServeFile(context.Writer, context.Request, filepath.Join(config.Instance.ThumbnailStorePath, image.Thumbnail))
}

var getImageRawHandler haruka.RequestHandler = func(context *haruka.Context) {
	id, err := context.GetPathParameterAsInt("id")
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	image, err := service.GetImageById(uint(id), "Library")
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	http.ServeFile(context.Writer, context.Request, filepath.Join(image.Library.Path, image.Path))
}
