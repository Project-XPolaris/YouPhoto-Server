package httpapi

import (
	"github.com/allentom/haruka"
	"github.com/allentom/haruka/middleware"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

var Logger = log.New().WithFields(log.Fields{
	"scope": "Application",
})

func GetEngine() *haruka.Engine {
	e := haruka.NewEngine()
	e.UseCors(cors.AllowAll())
	e.UseMiddleware(middleware.NewLoggerMiddleware())
	e.UseMiddleware(middleware.NewPaginationMiddleware("page", "pageSize", 1, 20))
	e.UseMiddleware(&AuthMiddleware{})
	e.UseMiddleware(&ReadUserMiddleware{})
	e.Router.GET("/libraries", getLibraryListHandler)
	e.Router.POST("/libraries", createLibraryHandler)
	e.Router.POST("/library/{id:[0-9]+}/scan", scanLibraryHandler)
	e.Router.DELETE("/library/{id:[0-9]+}", removeLibraryHandler)
	e.Router.GET("/images", getImageListHandler)
	e.Router.GET("/image/{id:[0-9]+}/thumbnail", getImageThumbnailHandler)
	e.Router.GET("/image/{id:[0-9]+}/raw", getImageRawHandler)
	e.Router.GET("/info", serviceInfoHandler)
	e.Router.GET("/user/current", getCurrentUserHandler)
	e.Router.GET("/readdir", readDirectoryHandler)
	e.Router.GET("/tasks", taskListHandler)
	return e
}
