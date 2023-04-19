package plugins

import (
	"fmt"
	"github.com/allentom/harukap"
	"github.com/projectxpolaris/youphoto/service/lora"
	"github.com/projectxpolaris/youphoto/service/sdw"
	"os"
)

type InitPlugin struct {
}

func (p *InitPlugin) OnInit(e *harukap.HarukaAppEngine) error {
	logger := e.LoggerPlugin.Logger.NewScope("Init plugin")
	logger.Info("deepdanbooru init")
	DefaultDeepdanbooruLauncher.Start()
	// init sdw client
	logger.Info("sdw init")
	sdwEnable := e.ConfigProvider.Manager.GetBool("sdw.enable")
	if sdwEnable {
		url := e.ConfigProvider.Manager.GetString("sdw.url")
		logger.Info(fmt.Sprintf("use sdw url %s", url))
		sdw.DefaultSDWClient = sdw.NewSDWClient(&sdw.Conf{
			Url: url,
		})
		_, err := sdw.DefaultSDWClient.GetModels()
		if err != nil {
			sdw.DefaultSDWClient = nil
			logger.Error(err)
			logger.Error("sdw init failed")
		}
	} else {
		logger.Info("sdw disabled")
	}

	// init preprocess dir
	preprocessPath := e.ConfigProvider.Manager.GetString("preprocess.outputpath")
	if preprocessPath == "" {
		logger.Error("preprocess output path is empty")
		return nil
	}
	logger.Info(fmt.Sprintf("preprocess output path %s", preprocessPath))
	stat, _ := os.Stat(preprocessPath)
	if stat != nil {
		err := os.MkdirAll(preprocessPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// init sdw client
	logger.Info("lora init")
	loraEnable := e.ConfigProvider.Manager.GetBool("lora.enable")
	if loraEnable {
		url := e.ConfigProvider.Manager.GetString("lora.url")
		logger.Info(fmt.Sprintf("use lora url %s", url))
		lora.DefaultLoraTrainClient = lora.NewLoraTrainClient(&lora.Conf{
			Url: url,
		})
		_, err := lora.DefaultLoraTrainClient.FetchInfo()
		if err != nil {
			lora.DefaultLoraTrainClient = nil
			logger.Error(err)
			logger.Error("lora init failed")
		}
	} else {
		logger.Info("lora disabled")
	}
	return nil
}
