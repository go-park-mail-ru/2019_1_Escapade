package server

import (
	"fmt"
	"os"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
)

// Port return port
func Port(conf *config.Configuration) (port string) {
	env := os.Getenv(conf.Server.PortURL)
	if os.Getenv(conf.Server.PortURL)[0] != ':' {
		port = ":" + env
	} else {
		port = env
	}
	fmt.Println("launched, look at us on " + conf.Server.Host + port)
	return
}
