package seed

import (
	"github.com/godcong/go-trait"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

// InitLogger ...
func InitLogger() *zap.SugaredLogger {
	trait.InitGlobalZapSugar()
	log = trait.ZapSugar()
	return log
}
