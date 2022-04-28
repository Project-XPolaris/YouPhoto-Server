package httpapi

import (
	"github.com/allentom/haruka"
	"github.com/projectxpolaris/youphoto/commons"
	"github.com/projectxpolaris/youphoto/plugins"
)

func AbortError(ctx *haruka.Context, err error, status int) {
	if apiError, ok := err.(*commons.APIError); ok {
		plugins.DefaultYouLogPlugin.Logger.Error(apiError.Err.Error())
		ctx.JSONWithStatus(haruka.JSON{
			"success": false,
			"err":     apiError.Desc,
			"code":    apiError.Code,
		}, status)
		return
	}
	plugins.DefaultYouLogPlugin.Logger.Error(err.Error())
	ctx.JSONWithStatus(haruka.JSON{
		"success": false,
		"err":     err.(error).Error(),
		"code":    "9999",
	}, status)
}

func MakeSuccessResponse(data interface{}, ctx *haruka.Context) {
	ctx.JSON(haruka.JSON{
		"success": true,
		"data":    data,
	})
}

func MakeListResponse(ctx *haruka.Context, page int, pageSize int, total int64, data interface{}) {
	ctx.JSON(haruka.JSON{
		"success":  true,
		"page":     page,
		"pageSize": pageSize,
		"count":    total,
		"result":   data,
	})
}
