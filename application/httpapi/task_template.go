package httpapi

import "github.com/projectxpolaris/youphoto/service"

type ScanLibraryDetail struct {
	Id          uint   `json:"id"`
	Path        string `json:"path"`
	Current     int64  `json:"current"`
	Total       int64  `json:"total"`
	CurrentPath string `json:"currentPath"`
	CurrentName string `json:"currentName"`
	Name        string `json:"name"`
}

func NewScanLibraryDetail(output *service.ScanTaskOutput) (*ScanLibraryDetail, error) {
	return &ScanLibraryDetail{
		Id:          output.Id,
		Path:        output.Path,
		Current:     output.Current,
		CurrentPath: output.CurrentPath,
		CurrentName: output.CurrentName,
		Name:        output.Name,
		Total:       output.Total,
	}, nil
}

type RemoveLibraryDetail struct {
	Id      uint   `json:"id"`
	Path    string `json:"path"`
	Name    string `json:"name"`
	Total   int64  `json:"total"`
	Current int64  `json:"current"`
}

func NewRemoveLibraryDetail(output *service.RemoveLibraryTaskOutput) (*RemoveLibraryDetail, error) {
	return &RemoveLibraryDetail{
		Id:      output.Id,
		Path:    output.Path,
		Name:    output.Name,
		Total:   output.Total,
		Current: output.Current,
	}, nil
}
