package main

import (
	"log"
	"net"
)

var MTU = 1000

func main() {
	sock := initSocket(":8090")
	defer sock.Close()
	remoteAddr1, err := net.ResolveUDPAddr("udp", "127.0.0.1:8091")
	remoteAddr2, err := net.ResolveUDPAddr("udp", "127.0.0.1:8092")
	if err != nil {
		log.Fatal("bad forward addresss", err)
	}

	for {
		buffer := make([]byte, MTU)
		// buffer := []byte{}

		bytes, _, err := sock.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal("unable to read UDP packet payload due to error:", err)
		}
		if bytes == 0 {
			log.Fatal("unable to read UDP packet, zero bytes in size")
		}
		buffer = buffer[:bytes]

		// send a plain UDP message back
		sock.WriteToUDP(buffer[0:bytes], remoteAddr1)
		sock.WriteToUDP(buffer[0:bytes], remoteAddr2)
	}
}

func initSocket(port string) *net.UDPConn {
	localAddr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Fatal("error resolving local IP port ", port, ":", err)
	}

	socket, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatal("Error starting listenUDP server:", err)
	}

	return socket
}
