package doctron_core

import (
	"time"

	"github.com/lampnick/doctron/converter"
)

type Doctron struct {
	config         converter.DoctronConfig
	buf            []byte
	convertElapsed time.Duration
}

func (d *Doctron) Log(format string, args ...interface{}) {
	d.config.IrisCtx.Application().Logger().Infof(format, args...)
}

type DoctronI interface {
	converter.Converter
}
