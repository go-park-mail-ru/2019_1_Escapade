package server

import (
	"net"
	"strconv"
)

// TODO deprecated - delete it
func FixPort(port string) string {
	if port[0] != ':' {
		return ":" + port
	}
	return port
}

func Port(port string) (string, int, error) {
	intPort, err := strconv.Atoi(port)
	if err != nil {
		return port, 0, err
	}
	return ":" + port, intPort, err
}

// GetIP return host of the service
func GetIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}
