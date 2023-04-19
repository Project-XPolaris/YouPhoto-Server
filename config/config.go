package config

import (
	"fmt"
	"github.com/allentom/harukap/config"
	"github.com/mitchellh/mapstructure"
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

type AuthConfig struct {
	Name   string
	Enable bool
	AppId  string
	Secret string
	Url    string
	Type   string
}
type Config struct {
	ThumbnailStorePath  string
	EnableAuth          bool
	YouPlusUrl          string
	Datasource          string
	YouPlusPathEnable   bool
	Auths               []*AuthConfig
	YouAuthConfig       *AuthConfig
	YouAuthConfigPrefix string
	EnableAnonymous     bool
	PreprocessPath      string
	ModelOutPath        string
}

func ReadConfig(provider *config.Provider) {
	configer := provider.Manager
	configer.SetDefault("addr", ":8000")
	configer.SetDefault("application", "YouPhoto Service")
	configer.SetDefault("instance", "main")
	configer.SetDefault("preprocess.outputpath", "./data/preprocess")
	configer.SetDefault("lora.outputpath", "./data/model_out")
	Instance = Config{
		ThumbnailStorePath: configer.GetString("thumbnails.store_path"),
		EnableAuth:         configer.GetBool("youplus.auth"),
		Datasource:         configer.GetString("datasource"),
		YouPlusPathEnable:  configer.GetBool("youplus.enablepath"),
		Auths:              make([]*AuthConfig, 0),
	}
	// read auth config
	rawAuth := configer.GetStringMap("auth")
	for key := range rawAuth {
		authConfig := &AuthConfig{}
		err := mapstructure.Decode(rawAuth[key], authConfig)
		if err != nil {
			panic(err)
		}
		Instance.Auths = append(Instance.Auths, authConfig)
		if authConfig.Type == "youauth" {
			Instance.YouAuthConfig = authConfig
			Instance.YouAuthConfigPrefix = fmt.Sprintf("auth.%s", key)
		}
		if authConfig.Type == "anonymous" {
			Instance.EnableAnonymous = configer.GetBool(fmt.Sprintf("auth.%s.enable", key))
		}
	}
	Instance.PreprocessPath = configer.GetString("preprocess.outputpath")
	Instance.ModelOutPath = configer.GetString("lora.outputpath")
}
