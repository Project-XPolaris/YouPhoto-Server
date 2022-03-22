package httpapi

type CurrentUserTemplate struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

func NewCurrentUserTemplate(id, username string) *CurrentUserTemplate {
	return &CurrentUserTemplate{
		Id:       id,
		Username: username,
	}
}
