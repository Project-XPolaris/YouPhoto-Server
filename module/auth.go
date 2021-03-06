package module

import (
	"github.com/allentom/harukap"
	"github.com/allentom/harukap/module/auth"
	"github.com/projectxpolaris/youphoto/config"
)

var Auth = &auth.AuthModule{
	Plugins: []harukap.AuthPlugin{},
}

func CreateAuthModule() {
	Auth.ConfigProvider = config.DefaultConfigProvider
	Auth.InitModule()
}
