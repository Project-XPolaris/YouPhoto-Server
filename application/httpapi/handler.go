package httpapi

import (
	"fmt"
	"github.com/allentom/haruka"
	"github.com/projectxpolaris/youphoto/config"
	"github.com/projectxpolaris/youphoto/module"
	"github.com/projectxpolaris/youphoto/plugins"
	"github.com/projectxpolaris/youphoto/service"
	"github.com/projectxpolaris/youphoto/service/sdw"
	"net/http"
	"os"
	"path/filepath"
)

var serviceInfoHandler haruka.RequestHandler = func(context *haruka.Context) {
	authMaps, err := module.Auth.GetAuthConfig()
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	context.JSON(haruka.JSON{
		"success":            true,
		"name":               "YouPhoto service",
		"authEnable":         config.Instance.EnableAuth,
		"authUrl":            fmt.Sprintf("%s/%s", config.Instance.YouPlusUrl, "user/auth"),
		"oauth":              true,
		"auth":               authMaps,
		"deepdanbooruEnable": plugins.DefaultDeepDanbooruPlugin.Enable,
		"sdwEnable":          sdw.DefaultSDWClient != nil,
	})
}

var readDirectoryHandler haruka.RequestHandler = func(context *haruka.Context) {
	rootPath := context.GetQueryString("path")
	if config.Instance.YouPlusPathEnable {
		token := context.Param["token"].(string)
		items, err := plugins.DefaultYouPlusPlugin.Client.ReadDir(rootPath, token)
		if err != nil {
			AbortError(context, err, http.StatusInternalServerError)
			return
		}
		data := make([]BaseFileItemTemplate, 0)
		for _, item := range items {
			template := BaseFileItemTemplate{}
			template.AssignWithYouPlusItem(item)
			data = append(data, template)
		}
		MakeSuccessResponse(haruka.JSON{
			"path":     rootPath,
			"sep":      "/",
			"files":    data,
			"backPath": filepath.Dir(rootPath),
		}, context)
		return
	} else {
		if len(rootPath) == 0 {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				AbortError(context, err, http.StatusInternalServerError)
				return
			}
			rootPath = homeDir
		}
		infos, err := service.ReadDirectory(rootPath)
		if err != nil {
			AbortError(context, err, http.StatusInternalServerError)
			return
		}
		data := make([]BaseFileItemTemplate, 0)
		for _, info := range infos {
			template := BaseFileItemTemplate{}
			template.Assign(info, rootPath)
			data = append(data, template)
		}
		MakeSuccessResponse(
			haruka.JSON{
				"path":     rootPath,
				"sep":      string(os.PathSeparator),
				"files":    data,
				"backPath": filepath.Dir(rootPath),
			}, context)
	}
}

var stateHandler haruka.RequestHandler = func(context *haruka.Context) {
	context.JSON(haruka.JSON{
		"success": true,
		"data": haruka.JSON{
			"nsfwEnable":                plugins.DefaultNSFWCheckPlugin.Enable,
			"deepdanbooruEnable":        plugins.DefaultDeepDanbooruPlugin.Enable,
			"imageClassificationEnable": plugins.DefaultImageClassifyPlugin.Enable,
			"sdwEnable":                 sdw.DefaultSDWClient != nil,
			"imageTaggerEnable":         plugins.DefaultImageTaggerPlugin.IsEnable(),
		},
	})
}
