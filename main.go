package main

import (
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

var bufferSize = 1490

func main() {
	log.Infoln("UDP Packet Duplicator")

	host, remotes, err := getPorts(os.Args[1:])
	checkErr(err)

	log.Infoln("Host Port:", host)
	log.Infoln("Remote Ports:", remotes)

	hostConn, err := net.ListenUDP("udp", host)
	checkErr(err)
	defer hostConn.Close()

	proxy(hostConn, remotes)
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func proxy(sock *net.UDPConn, forwards []*net.UDPAddr) {
	buffer := make([]byte, bufferSize)

	for {
		bytes, _, err := sock.ReadFromUDP(buffer)
		if err != nil {
			log.Errorln(err)
			continue
		}
		if bytes == 0 {
			log.Errorln("zero byte packet")
			continue
		}

		log.Infoln(bytes, "byte packet recieved, forwarding")

		for _, addr := range forwards {
			_, err := sock.WriteToUDP(buffer[0:bytes], addr)
			if err != nil {
				log.Errorln(err)
				continue
			}
		}
	}
}
