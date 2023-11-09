package main

import (
	"github.com/allentom/haruka"
	"github.com/allentom/harukap"
	"github.com/allentom/harukap/cli"
	"github.com/projectxpolaris/youphoto/application/httpapi"
	"github.com/projectxpolaris/youphoto/config"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/module"
	"github.com/projectxpolaris/youphoto/plugins"
	"github.com/sirupsen/logrus"
	"net/http"
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
	//bootLogger := plugins.DefaultYouLogPlugin.Logger.NewScope("boot")
	appEngine := harukap.NewHarukaAppEngine()
	appEngine.ConfigProvider = config.DefaultConfigProvider
	appEngine.LoggerPlugin = plugins.DefaultYouLogPlugin
	plugins.CreateDefaultYouPlusPlugin()
	appEngine.UsePlugin(plugins.DefaultYouPlusPlugin)
	appEngine.UsePlugin(database.DefaultPlugin)
	appEngine.UsePlugin(plugins.DefaultThumbnailServicePlugin)
	appEngine.UsePlugin(&plugins.DefaultRegisterPlugin)
	appEngine.UsePlugin(plugins.StorageEnginePlugin)
	appEngine.UsePlugin(plugins.DefaultImageClassifyPlugin)
	appEngine.UsePlugin(plugins.DefaultNSFWCheckPlugin)
	appEngine.UsePlugin(plugins.DefaultDeepDanbooruPlugin)
	appEngine.UsePlugin(plugins.DefaultImageTaggerPlugin)
	appEngine.UsePlugin(&plugins.InitPlugin{})
	if config.Instance.YouAuthConfig != nil {
		plugins.CreateYouAuthPlugin()
		plugins.DefaultYouAuthOauthPlugin.ConfigPrefix = config.Instance.YouAuthConfigPrefix
		appEngine.UsePlugin(plugins.DefaultYouAuthOauthPlugin)
	}
	// init module
	module.CreateAuthModule()
	module.Auth.AuthMiddleware.OnError = func(ctx *haruka.Context, err error) {
		httpapi.AbortError(ctx, err, http.StatusForbidden)
		ctx.Abort()
	}
	module.CreateTaskModule()
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
