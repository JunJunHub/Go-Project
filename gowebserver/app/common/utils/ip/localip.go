package ip

import (
	"fmt"
	"log"
	"net"
)

//GetOutboundIP Get preferred outbound ip of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	return localAddr.IP.String()
}

//GetLocalIP 获取本机IP地址
func GetLocalIP() (ips []string, err error) {
	ips = make([]string, 20)
	address, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, addr := range address {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}

		ips = append(ips, ipAddr.IP.String())
	}
	return
}
