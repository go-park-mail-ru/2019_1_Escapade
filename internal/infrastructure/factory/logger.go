package factory

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	uc "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/usecase/logger"
)

func NewLogger() infrastructure.LoggerI {
	return uc.NewLog()
}
