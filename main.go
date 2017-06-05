package main

import (
	"log"
	"net"
	"os"
)

var MTU = 1000

func main() {

	host, remotes, err := getPorts(os.Args)
	checkErr(err)

	hostConn, err := net.ListenUDP("udp", host)
	checkErr(err)
	defer hostConn.Close()

	proxy(hostConn, remotes)
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln("Fatal Error:", err)
	}
}

func proxy(sock *net.UDPConn, forwards []*net.UDPAddr) {
	buffer := make([]byte, MTU)

	for {
		bytes, _, err := sock.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Error: unable to read UDP packet payload:", err)
			continue
		}
		if bytes == 0 {
			log.Println("Error: unable to read UDP packet, zero bytes in size")
			continue
		}

		buffer = buffer[:bytes]

		for _, addr := range forwards {
			sock.WriteToUDP(buffer[0:bytes], addr)
		}
	}
}
