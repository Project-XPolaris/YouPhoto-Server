package service

import (
	"fmt"
	"github.com/allentom/harukap"
	"github.com/allentom/harukap/plugins/tagger"
	"github.com/project-xpolaris/youplustoolkit/youlog"
	"github.com/projectxpolaris/youphoto/database"
	"github.com/projectxpolaris/youphoto/plugins"
)

type AutoTaggerService struct {
	Engine *harukap.HarukaAppEngine
	logger *youlog.Scope
}

func NewAutoTaggerService(e *harukap.HarukaAppEngine) *AutoTaggerService {
	serviceLogger := e.LoggerPlugin.Logger.NewScope("Auto tagger")
	return &AutoTaggerService{
		Engine: e,
		logger: serviceLogger,
	}
}

func (s *AutoTaggerService) Start() error {
	return s.Process()
}

func (s *AutoTaggerService) Process() error {
	if !plugins.DefaultImageTaggerPlugin.IsEnable() {
		s.logger.Info("image tagger plugin is disabled")
	}
	s.logger.Info("start auto tagger")
	for {
		images, err := GetImagesWithoutTags()
		if err != nil {
			return err
		}
		if len(images) == 0 {
			break
		}
		for _, image := range images {
			s.logger.Info(fmt.Sprintf("start tag image %s", image.Path))
			_, err = TagImageById(image.ID, tagger.ModelAuto, 0.7)
			if err != nil {
				s.logger.Error(err)
				s.logger.Error(fmt.Sprintf("tag image %s failed", image.Path))
			}
			image.Tagged = true
			database.Instance.Save(&image)

		}
	}
	return nil
}
