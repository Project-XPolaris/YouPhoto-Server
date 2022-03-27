package config

import (
	"github.com/allentom/harukap/config"
	"os"
)

var DefaultConfigProvider *config.Provider

func InitConfigProvider() error {
	var err error
	customConfigPath := os.Getenv("YOUPHOTO_CONFIG_PATH")
	DefaultConfigProvider, err = config.NewProvider(func(provider *config.Provider) {
		ReadConfig(provider)
	}, customConfigPath)
	return err
}

var Instance Config

type Config struct {
	ThumbnailStorePath  string
	ThumbnailServiceUrl string
	ThumbnailProvider   string
	EnableAuth          bool
	YouPlusUrl          string
	Datasource          string
	YouPlusPathEnable   bool
}

func ReadConfig(provider *config.Provider) {
	configer := provider.Manager
	configer.SetDefault("addr", ":8000")
	configer.SetDefault("application", "YouPhoto Service")
	configer.SetDefault("instance", "main")

	Instance = Config{
		ThumbnailStorePath:  configer.GetString("thumbnails.store_path"),
		ThumbnailServiceUrl: configer.GetString("thumbnails.service_url"),
		ThumbnailProvider:   configer.GetString("thumbnails.provider"),
		EnableAuth:          configer.GetBool("youplus.auth"),
		Datasource:          configer.GetString("datasource"),
		YouPlusPathEnable:   configer.GetBool("youplus.enablepath"),
	}
}
