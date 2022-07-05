package plugins

import (
	"github.com/allentom/harukap/plugins/storage"
	"github.com/projectxpolaris/youphoto/config"
)

var StorageEnginePlugin = &storage.Engine{}

func GetDefaultStorage() storage.FileSystem {
	defaultStorageName := config.DefaultConfigProvider.Manager.GetString("storage.default")
	return StorageEnginePlugin.GetStorage(defaultStorageName)
}
