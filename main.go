package main

import (
	"errors"
	"net"
	// _ "net/http/pprof"
	"os"

	log "github.com/sirupsen/logrus"
)

var bufferSize = 1500

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
	t1 := make(chan []byte, 100)
	t2 := make(chan []byte, 100)
	chanSlice := []chan []byte{t1, t2}
	go sender(forwards[0], t1)
	go sender(forwards[1], t2)
	for {
		err := proxyPacket(sock, chanSlice, buffer)
		if err != nil {
			log.Errorln(err)
		}
	}
}

func proxyPacket(sock *net.UDPConn, forwards []chan []byte, buffer []byte) error {
	bytes, _, err := sock.ReadFromUDP(buffer)
	if err != nil {
		return err
	}
	if bytes == 0 {
		return errors.New("zero byte packet")
	}

	log.Infoln(bytes, "byte packet recieved, forwarding")

	for _, addr := range forwards {
		addr <- buffer
	}

	return nil
}

func sender(to *net.UDPAddr, get chan []byte) {
	buffer := make([]byte, bufferSize)
	conn, err := net.Dial("udp", "127.0.0.1:1234")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		buffer = <-get
		_, err := conn.Write(buffer)
		if err != nil {
			log.Errorln(err)
		}
	}
}
