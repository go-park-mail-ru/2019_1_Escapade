package server

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

func Port(port string) (string, int, error) {
	intPort, err := strconv.Atoi(port)
	if err != nil {
		return port, 0, err
	}
	return ":" + port, intPort, err
}

func PortString(port int) string {
	return ":" + utils.String(port)
}

func GetIP(subnet *string) string {
	var ips string
	ifaces, err := net.Interfaces()
	if err != nil {
		return err.Error()
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return err.Error()
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			ipsting := ip.String()
			fmt.Println("ips:", ipsting)
			if subnet == nil {
				ips += " " + ipsting
			} else if strings.HasPrefix(ipsting, *subnet) {
				return ipsting
			}
		}
	}
	if subnet == nil {
		return ips
	}
	return "error: no networks. Change subnet!"
}
