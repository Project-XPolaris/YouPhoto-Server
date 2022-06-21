package httpapi

import (
	"github.com/allentom/haruka"
	"github.com/projectxpolaris/youphoto/module"
)

var taskListHandler haruka.RequestHandler = func(context *haruka.Context) {
	data, _ := module.Task.SerializerTemplateList()
	MakeSuccessResponse(data, context)
}
