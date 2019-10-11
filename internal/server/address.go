package server

import (
	"net"
	"strconv"
	"strings"
)

func FixHost(host string) string {
	if strings.Contains(host, "http://") {
		host = strings.Replace(host, "http://", "", 1)
	}
	if strings.Contains(host, "https://") {
		host = strings.Replace(host, "https://", "", 1)
	}
	return host
}

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

func GetIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}
