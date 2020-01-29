package factory

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	uc "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/usecase/middleware"
)

func NewMiddlewareAuth(c config.Configuration,
	auth infrastructure.AuthService,
	log infrastructure.LoggerI,
	metrics infrastructure.MetricsI) []infrastructure.MiddlewareI {
	return []infrastructure.MiddlewareI{
		uc.NewRecover(),
		uc.NewCORS(c.Cors),
		uc.NewAuth(auth, log)}
}

func NewMiddlewareNonAuth(c config.Configuration, subnet string,
	log infrastructure.LoggerI,
	metrics infrastructure.MetricsI) []infrastructure.MiddlewareI {
	return []infrastructure.MiddlewareI{
		uc.NewRecover(),
		uc.NewCORS(c.Cors),
		uc.NewLogger(log),
		uc.NewMetrics(metrics, subnet)}
}

//  api.Use(mi.CORS(cors), mi.Metrics(subnet))
// apiWithAuth.Use(mi.Auth(config.Cookie, config.Auth, config.AuthClient))
