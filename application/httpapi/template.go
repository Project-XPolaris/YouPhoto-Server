package httpapi

import (
	"github.com/project-xpolaris/youplustoolkit/youplus"
	"github.com/projectxpolaris/youphoto/database"
	"os"
	"path/filepath"
)

type BaseFileItemTemplate struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Path string `json:"path"`
}

func (t *BaseFileItemTemplate) Assign(info os.FileInfo, rootPath string) {
	if info.IsDir() {
		t.Type = "Directory"
	} else {
		t.Type = "File"
	}
	t.Name = info.Name()
	t.Path = filepath.Join(rootPath, info.Name())
}

func (t *BaseFileItemTemplate) AssignWithYouPlusItem(item youplus.ReadDirItem) {
	t.Type = item.Type
	t.Path = item.Path
	t.Name = filepath.Base(item.Path)
}

type SDWConfigTemplate struct {
	Id         uint   `json:"id"`
	Name       string `json:"name"`
	Config     string `json:"config"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}

func NewSDWConfigTemplate(config *database.SdwConfig) *SDWConfigTemplate {
	return &SDWConfigTemplate{
		Id:         config.ID,
		Name:       config.Name,
		Config:     config.Config,
		CreateTime: config.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdateTime: config.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func NewSDWConfigTemplateList(configList []*database.SdwConfig) []*SDWConfigTemplate {
	result := make([]*SDWConfigTemplate, 0)
	for _, config := range configList {
		result = append(result, NewSDWConfigTemplate(config))
	}
	return result
}

type LoraConfigTemplate struct {
	Id         uint   `json:"id"`
	Name       string `json:"name"`
	Config     string `json:"config"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}

func NewLoraConfigTemplate(config *database.LoraConfig) *LoraConfigTemplate {
	return &LoraConfigTemplate{
		Id:         config.ID,
		Name:       config.Name,
		Config:     config.Config,
		CreateTime: config.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdateTime: config.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func NewLoraConfigTemplateList(configList []*database.LoraConfig) []*LoraConfigTemplate {
	result := make([]*LoraConfigTemplate, 0)
	for _, config := range configList {
		result = append(result, NewLoraConfigTemplate(config))
	}
	return result
}
