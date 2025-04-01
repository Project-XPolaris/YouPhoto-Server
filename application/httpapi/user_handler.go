package httpapi

import (
	"strconv"

	"github.com/allentom/haruka"
	"github.com/projectxpolaris/youphoto/database"
)

var getCurrentUserHandler haruka.RequestHandler = func(context *haruka.Context) {
	user := context.Param["claim"].(*database.User)
	data := NewCurrentUserTemplate(strconv.FormatUint(uint64(user.ID), 10), user.Username)
	MakeSuccessResponse(data, context)
}
