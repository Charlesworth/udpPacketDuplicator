package main

import (
	"errors"
	"net"
	// _ "net/http/pprof"
	"os"

	log "github.com/sirupsen/logrus"
)

var bufferSize = 1490

func main() {
	// add the following to provide profiling
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	log.Infoln("UDP Packet Duplicator")
	checkInput(os.Args[1:])

	ports, err := toPorts(os.Args[1:])
	checkErr(err)

	host := ports[0]
	remotes := ports[1:]

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
		err := proxyPacket(sock, forwards, buffer)
		if err != nil {
			log.Errorln(err)
		}
	}
}

func proxyPacket(sock *net.UDPConn, forwards []*net.UDPAddr, buffer []byte) error {
	bytes, _, err := sock.ReadFromUDP(buffer)
	if err != nil {
		return err
	}
	if bytes == 0 {
		return errors.New("zero byte packet")
	}

	log.Infoln(bytes, "byte packet recieved, forwarding")

	for _, addr := range forwards {
		_, err := sock.WriteToUDP(buffer[0:bytes], addr)
		if err != nil {
			return err
		}
	}

	return nil
}
