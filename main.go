package main

import (
	"github.com/allentom/harukap"
	"github.com/allentom/harukap/cli"
	"github.com/projectxpolaris/youphoto/application/httpapi"
	"github.com/projectxpolaris/youphoto/config"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/plugins"
	"github.com/projectxpolaris/youphoto/utils"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	err := config.InitConfigProvider()
	if err != nil {
		logrus.Fatal(err)
	}
	err = plugins.DefaultYouLogPlugin.OnInit(config.DefaultConfigProvider)
	if err != nil {
		logrus.Fatal(err)
	}
	bootLogger := plugins.DefaultYouLogPlugin.Logger.NewScope("boot")
	bootLogger.Info("init thumbnail path")
	isThumbnailsStoreExist := utils.CheckFileExist(config.Instance.ThumbnailStorePath)
	if !isThumbnailsStoreExist {
		bootLogger.Info("thumbnail folder not exist, create it")
		err = os.Mkdir(config.Instance.ThumbnailStorePath, os.ModePerm)
		if err != nil {
			logrus.Fatal(err)
		}
	}
	appEngine := harukap.NewHarukaAppEngine()
	appEngine.ConfigProvider = config.DefaultConfigProvider
	appEngine.LoggerPlugin = plugins.DefaultYouLogPlugin
	appEngine.UsePlugin(plugins.DefaultYouPlusPlugin)
	appEngine.UsePlugin(database.DefaultPlugin)
	appEngine.UsePlugin(plugins.DefaultThumbnailServicePlugin)
	appEngine.UsePlugin(&plugins.DefaultRegisterPlugin)
	if config.Instance.YouAuthConfig != nil {
		plugins.DefaultYouAuthOauthPlugin.ConfigPrefix = config.Instance.YouAuthConfigPrefix
		appEngine.UsePlugin(plugins.DefaultYouAuthOauthPlugin)
	}
	appEngine.HttpService = httpapi.GetEngine()
	if err != nil {
		logrus.Fatal(err)
	}
	appWrap, err := cli.NewWrapper(appEngine)
	if err != nil {
		logrus.Fatal(err)
	}
	appWrap.RunApp()
}
