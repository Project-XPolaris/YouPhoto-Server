package httpapi

import (
	"github.com/allentom/haruka"
	"github.com/allentom/haruka/middleware"
	"github.com/projectxpolaris/youphoto/module"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var Logger = log.New().WithFields(log.Fields{
	"scope": "Application",
})

func GetEngine() *haruka.Engine {
	e := haruka.NewEngine()
	e.UseCors(cors.AllowAll())
	module.Auth.AuthMiddleware.OnError = func(ctx *haruka.Context, err error) {
		AbortError(ctx, err, http.StatusForbidden)
		ctx.Abort()
	}
	module.Auth.AuthMiddleware.RequestFilter = func(c *haruka.Context) bool {
		NoAuthPath := []string{
			"/oauth/youauth",
			"/oauth/youauth/password",
			"/oauth/youplus",
			"/info",
			"/state",
		}
		for _, path := range NoAuthPath {
			if c.Pattern == path {
				return false
			}
		}
		return true
	}
	e.UseMiddleware(module.Auth.AuthMiddleware)
	e.UseMiddleware(middleware.NewPaginationMiddleware("page", "pageSize", 1, 20))
	e.Router.GET("/libraries", getLibraryListHandler)
	e.Router.POST("/libraries", createLibraryHandler)
	e.Router.POST("/library/{id:[0-9]+}/scan", scanLibraryHandler)
	e.Router.DELETE("/library/{id:[0-9]+}", removeLibraryHandler)
	e.Router.GET("/images", getImageListHandler)
	e.Router.GET("/image/{id:[0-9]+}", getImageHandler)
	e.Router.GET("/image/{id:[0-9]+}/thumbnail", getImageThumbnailHandler)
	e.Router.GET("/image/{id:[0-9]+}/near", getNearImageHandler)
	e.Router.GET("/image/{id:[0-9]+}/tagger", getImageTaggerHandler)
	e.Router.GET("/image/{id:[0-9]+}/album", getAlbumDetailHandler)
	e.Router.POST("/image/{id:[0-9]+}/upscale", upscaleImageHandler)
	e.Router.DELETE("/images", deleteImageByIdsHandler)
	e.Router.POST("/image/upload", uploadImageByFileHandler)
	e.Router.POST("/album/{id:[0-9]+}/image", addImageToAlbumHandler)
	e.Router.DELETE("/album/{id:[0-9]+}/image", removeImageFromAlbumHandler)
	e.Router.POST("/album", createAlbumHandler)
	e.Router.GET("/albums", getAlbumListHandler)
	e.Router.DELETE("/album/{id:[0-9]+}", removeAlbumHandler)
	e.Router.GET("/tags", getImageTagListHandler)
	e.Router.POST("/image/base64", uploadImageByBase64Handler)
	e.Router.POST("/image/deepdanbooru", deepdanbooruHandler)
	e.Router.POST("/color/match", getColorMatchHandler)
	e.Router.GET("/thumbnail/{id}", getThumbnailHandler)
	e.Router.GET("/image/{id:[0-9]+}/raw", getImageRawHandler)
	e.Router.GET("/info", serviceInfoHandler)
	e.Router.GET("/user/current", getCurrentUserHandler)
	e.Router.GET("/readdir", readDirectoryHandler)
	e.Router.GET("/tasks", taskListHandler)
	e.Router.GET("/task/object", module.Task.GetTaskByIdHandler)
	e.Router.GET("/oauth/youauth", generateAccessCodeWithYouAuthHandler)
	e.Router.POST("/oauth/youauth/password", generateAccessCodeWithYouAuthPasswordHandler)
	e.Router.POST("/oauth/youplus", YouPlusLoginHandler)
	e.Router.GET("/user/auth", youPlusTokenHandler)
	e.Router.GET("/sdw/models", getModelsHandler)
	e.Router.POST("/sdw/currentModel", switchModelHandler)
	e.Router.GET("/sdw/info", sdwInfoHandler)
	e.Router.POST("/sdw/text2image", text2ImageHaldler)
	e.Router.GET("/sdw/samplers", getSamplerListHandler)
	e.Router.GET("/sdw/upscalers", getUpscalerListHandler)
	e.Router.GET("/sdw/progress", getProgressHandler)
	e.Router.POST("/sdw/config", saveConfigHandler)
	e.Router.GET("/sdw/config", getSDWConfigListHandler)
	e.Router.DELETE("/sdw/config", deleteSDWConfigHandler)
	e.Router.POST("/sdw/interrupt", intrruptHandler)
	e.Router.GET("/sdw/skip", skipHandler)
	e.Router.POST("/sdw/preprocess", newPreprocessHandler)
	e.Router.POST("/lora/config", saveLoraConfigHandler)
	e.Router.POST("/lora/train", loraTrainHandler)
	e.Router.GET("/lora/config", getLoraConfigListHandler)
	e.Router.DELETE("/lora/config", deleteLoraConfigHandler)
	e.Router.GET("/lora/task/interrupt", loraInterruptHandler)
	e.Router.GET("/tagger/models", getImageTaggerModelHandler)
	e.Router.GET("/state", stateHandler)
	module.Task.AddConverter(NewScanLibraryDetail)
	module.Task.AddConverter(NewRemoveLibraryDetail)
	return e
}
