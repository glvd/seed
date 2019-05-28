package seed

import (
	"github.com/godcong/go-trait"
	"go.uber.org/zap"
)

var log = trait.NewZapSugar(zap.String("package", "seed"))
