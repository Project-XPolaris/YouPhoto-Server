package httpapi

import "github.com/projectxpolaris/youphoto/database"

const TimeFormat = "2006-01-02 15:04:05"

type BaseLibraryTemplate struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func NewBaseLibraryTemplate(library *database.Library) BaseLibraryTemplate {
	return BaseLibraryTemplate{
		Id:        library.Model.ID,
		Name:      library.Name,
		Path:      library.Path,
		CreatedAt: library.CreatedAt.Format(TimeFormat),
		UpdatedAt: library.UpdatedAt.Format(TimeFormat),
	}
}

func NewBaseLibraryTemplateList(libraries []*database.Library) []BaseLibraryTemplate {
	baseLibraryTemplates := make([]BaseLibraryTemplate, 0)
	for _, library := range libraries {
		baseLibraryTemplates = append(baseLibraryTemplates, NewBaseLibraryTemplate(library))
	}
	return baseLibraryTemplates
}
