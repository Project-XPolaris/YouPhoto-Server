package httpapi

import "github.com/allentom/haruka"

var getCurrentUserHandler haruka.RequestHandler = func(context *haruka.Context) {
	data := NewCurrentUserTemplate(context.Param["uid"].(string), context.Param["username"].(string))
	MakeSuccessResponse(data, context)
}
