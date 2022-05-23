package httpapi

import (
	"github.com/allentom/haruka"
	"github.com/projectxpolaris/youphoto/config"
	"github.com/projectxpolaris/youphoto/service"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AuthMiddleware struct {
}

var NoAuthPaths = []string{
	"/info",
	"/oauth/youauth",
	"/oauth/youplus",
}

func (m AuthMiddleware) OnRequest(c *haruka.Context) {
	if config.Instance.EnableAnonymous {
		return
	}
	for _, path := range NoAuthPaths {
		if c.Request.URL.Path == path {
			return
		}
	}
	claim, err := service.ParseAuthHeader(c)
	if err != nil {
		c.Interrupt()
		logrus.Error(err)
		AbortError(c, err, http.StatusForbidden)
		return
	}
	c.Param["claim"] = claim
}
