package httpapi

import (
	"github.com/allentom/haruka"
	"github.com/projectxpolaris/youphoto/service"
)

var taskListHandler haruka.RequestHandler = func(context *haruka.Context) {
	tasks := service.DefaultTaskPool.Tasks
	data := NewTaskListTemplate(tasks)
	MakeSuccessResponse(data, context)
}
