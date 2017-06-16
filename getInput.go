package main

import (
	"log"
	"net"
	"os"
)

func toPorts(inputStrings []string) (ports []*net.UDPAddr, err error) {
	for _, addrString := range inputStrings {
		addr, err := net.ResolveUDPAddr("udp", addrString)
		if err != nil {
			return nil, err
		}

		ports = append(ports, addr)
	}
	return
}

func checkInput(input []string) {
	if len(input) < 2 || input[1] == "-h" {
		log.Println("UDP Packet Duplicator")
		log.Println("Usage: udpPacketDuplicator [:hostPort] [:remotePort]...")
		os.Exit(0)
	}
}
