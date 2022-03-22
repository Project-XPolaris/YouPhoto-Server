package httpapi

import "github.com/projectxpolaris/youphoto/service"

type TaskTemplate struct {
	Id     string      `json:"id"`
	Type   string      `json:"type"`
	Detail interface{} `json:"detail"`
}

func NewTaskListTemplate(taskList []service.Task) []*TaskTemplate {
	taskTemplateList := make([]*TaskTemplate, 0)
	for _, task := range taskList {
		taskTemplateList = append(taskTemplateList, NewTaskTemplate(task))
	}
	return taskTemplateList
}
func NewTaskTemplate(task service.Task) *TaskTemplate {
	data := &TaskTemplate{
		Id:   task.GetId(),
		Type: service.TaskTypeNameMapping[task.GetType()],
	}
	output := task.GetOutput()
	switch output.(type) {
	case *service.ScanTaskOutput:
		data.Detail = NewScanLibraryDetail(output.(*service.ScanTaskOutput))
	case *service.RemoveLibraryTaskOutput:
		data.Detail = NewRemoveLibraryDetail(output.(*service.RemoveLibraryTaskOutput))
	}

	return data
}

type ScanLibraryDetail struct {
	Id          uint   `json:"id"`
	Path        string `json:"path"`
	Current     int64  `json:"current"`
	CurrentPath string `json:"currentPath"`
	CurrentName string `json:"currentName"`
}

func NewScanLibraryDetail(output *service.ScanTaskOutput) *ScanLibraryDetail {
	return &ScanLibraryDetail{
		Id:          output.Id,
		Path:        output.Path,
		Current:     output.Current,
		CurrentPath: output.CurrentPath,
		CurrentName: output.CurrentName,
	}
}

type RemoveLibraryDetail struct {
	Id   uint   `json:"id"`
	Path string `json:"path"`
}

func NewRemoveLibraryDetail(output *service.RemoveLibraryTaskOutput) *RemoveLibraryDetail {
	return &RemoveLibraryDetail{
		Id:   output.Id,
		Path: output.Path,
	}
}
