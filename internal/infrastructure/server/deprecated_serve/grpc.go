package serve

/*
import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"golang.org/x/net/netutil"
	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

func ServeGRPC(grpcServer *grpc.Server, serverConfig config.Server, port string, lastFunc func()) error {
	errChan := make(chan error)
	stopChan := make(chan os.Signal)

	defer func() {
		close(stopChan)
		close(errChan)
		lastFunc()
	}()

	l, err := net.Listen("tcp", port)

	if err != nil {
		utils.Debug(true, "Listen error", err.Error())
		return err
	}

	defer l.Close()

	l = netutil.LimitListener(l, serverConfig.MaxConn)

	signal.Notify(stopChan, os.Interrupt)

	go func() {
		if err := grpcServer.Serve(l); err != nil {
			errChan <- err
			close(errChan)
			utils.Debug(false, "Serving error:", err.Error())
		}
	}()

	// block until either OS signal, or server fatal error
	select {
	case err := <-errChan:
		utils.Debug(false, "Fatal error: ", err.Error())
		return err
	case <-stopChan:
		grpcServer.GracefulStop()
	}
	return nil
}*/
