package service

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

var scanTargetExtensions = []string{
	".jpg", ".png", ".jpeg", ".bmp",
}
var StopError = errors.New("stop")

// YouPhotoConfig youphoto.json 配置文件结构
type YouPhotoConfig struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type ImageScanner struct {
	BasePath   string
	OnHit      func(string) error
	OnAlbumHit func(string, string, []string) error // 参数：目录路径，相册名称，标签列表
}

func NewImageScanner(basePath string) *ImageScanner {
	return &ImageScanner{
		BasePath: basePath,
	}
}

func (s *ImageScanner) Scan() error {
	err := afero.Walk(AppFs, s.BasePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			// 检查目录中是否有 youphoto.json 文件
			configPath := filepath.Join(path, "youphoto.json")
			if exists, _ := afero.Exists(AppFs, configPath); exists {
				// 读取配置文件
				configData, readErr := afero.ReadFile(AppFs, configPath)
				if readErr == nil {
					var config YouPhotoConfig
					if parseErr := json.Unmarshal(configData, &config); parseErr == nil {
						if config.Name != "" && s.OnAlbumHit != nil {
							// 调用相册回调函数，传递相册名称和标签
							s.OnAlbumHit(path, config.Name, config.Tags)
						}
					}
				}
			}
			return nil
		}
		match := strings.ToLower(info.Name())
		if strings.HasPrefix(match, ".") {
			return nil
		}
		for _, extension := range scanTargetExtensions {
			if strings.HasSuffix(match, extension) {
				if s.OnHit != nil {
					err = s.OnHit(path)
					if err == StopError {
						return err
					}
				}
			}
		}
		return nil
	})
	return err
}
