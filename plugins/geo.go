package plugins

import (
	"github.com/allentom/harukap"
	"googlemaps.github.io/maps"
)

var DefaultGeoPlugin = &GeoPlugin{}

type GeoPlugin struct {
	Client *maps.Client
}

func (p *GeoPlugin) OnInit(e *harukap.HarukaAppEngine) error {
	logger := e.LoggerPlugin.Logger.NewScope("GeoPlugin")
	confManager := e.ConfigProvider.Manager
	enableGeo := confManager.GetBool("geo.enable")
	if !enableGeo {
		logger.Info("geo plugin is disabled")
		return nil
	}
	apikey := confManager.GetString("geo.apikey")
	if apikey == "" {
		logger.Error("geo apikey is empty")
		return nil
	}
	client, err := maps.NewClient(maps.WithAPIKey(apikey))
	if err != nil {
		logger.Error(err)
		return err
	}
	p.Client = client
	logger.Info("geo plugin init success")
	return nil

}

func (p *GeoPlugin) IsEnable() bool {
	return p.Client != nil
}
