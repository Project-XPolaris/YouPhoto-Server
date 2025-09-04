package task

import (
	"path/filepath"
	"strings"

	"github.com/allentom/harukap/module/task"
	"github.com/projectxpolaris/youphoto/service"
)

type ScanImageFileTaskOption struct {
	Uid          string
	Path         string
	ParentTaskId string
}
type ScanImageFileTaskOutput struct {
}
type AlbumInfo struct {
	Name string
	Tags []string
}

type ScanImageFileTask struct {
	*task.BaseTask
	option     *ScanImageFileTaskOption
	TaskOutput *ScanImageFileTaskOutput
	pathList   []string
	albumMap   map[string]*AlbumInfo // 目录路径 -> 相册信息的映射
}

func (t *ScanImageFileTask) Stop() error {
	return nil
}

func (t *ScanImageFileTask) Start() error {
	scanner := service.NewImageScanner(t.option.Path)
	scanner.OnHit = func(path string) error {
		t.pathList = append(t.pathList, path)
		return nil
	}
	// 添加相册发现回调
	scanner.OnAlbumHit = func(dirPath, albumName string, tags []string) error {
		t.albumMap[dirPath] = &AlbumInfo{
			Name: albumName,
			Tags: tags,
		}
		return nil
	}
	err := scanner.Scan()
	if err != nil {
		return t.AbortError(err)
	}
	t.Done()
	return nil
}

func (t *ScanImageFileTask) Output() (interface{}, error) {
	return t.TaskOutput, nil
}

func NewScanImageFileTask(option *ScanImageFileTaskOption) *ScanImageFileTask {
	t := &ScanImageFileTask{
		BaseTask:   task.NewBaseTask(TypeScanImageFile, option.Uid, task.GetStatusText(nil, task.StatusRunning)),
		TaskOutput: &ScanImageFileTaskOutput{},
		option:     option,
		pathList:   make([]string, 0),
		albumMap:   make(map[string]*AlbumInfo),
	}
	t.ParentTaskId = option.ParentTaskId
	return t
}

// GetAlbumInfoForPath 根据图片路径获取相册信息
func (t *ScanImageFileTask) GetAlbumInfoForPath(imagePath string) *AlbumInfo {
	imageDir := filepath.Dir(imagePath)
	// 检查是否有直接匹配的目录
	if albumInfo, exists := t.albumMap[imageDir]; exists {
		return albumInfo
	}
	// 检查父目录中是否有相册配置
	for dirPath, albumInfo := range t.albumMap {
		if strings.HasPrefix(imagePath, dirPath) {
			return albumInfo
		}
	}
	return nil
}

// GetAlbumNameForPath 根据图片路径获取相册名称（保持向后兼容）
func (t *ScanImageFileTask) GetAlbumNameForPath(imagePath string) string {
	albumInfo := t.GetAlbumInfoForPath(imagePath)
	if albumInfo != nil {
		return albumInfo.Name
	}
	return ""
}
