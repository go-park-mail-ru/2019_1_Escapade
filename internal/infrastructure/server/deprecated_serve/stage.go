package serve

/*
import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
	"google.golang.org/grpc"
)

func StageRunServerHTTP(configServer config.Server, handler http.Handler, port string) server.Stage {
	return func() error {
		var (
			exitFunc = func() { utils.Debug(false, "✗✗✗ Exit ✗✗✗") }
			srv      = ConfigureServer(handler, configServer, port)
		)
		return ServeHTTP(srv, configServer, exitFunc)
	}
}

// run server
func StageRunServerGRPC(configServer config.Server, GRPC *grpc.Server, port string) server.Stage {
	return func() error {
		var (
			exitFunc = func() { utils.Debug(false, "✗✗✗ Exit ✗✗✗") }
		)
		return ServeGRPC(GRPC, configServer, port, exitFunc)
	}
}
*/
