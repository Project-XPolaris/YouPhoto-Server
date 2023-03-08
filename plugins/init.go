package plugins

import (
	"github.com/allentom/harukap"
)

type InitPlugin struct {
}

func (p *InitPlugin) OnInit(e *harukap.HarukaAppEngine) error {
	return nil
}
