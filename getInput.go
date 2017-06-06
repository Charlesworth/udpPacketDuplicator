package main

import (
	"log"
	"net"
	"os"
)

func getPorts(inputStrings []string) (hostPort *net.UDPAddr, remotePorts []*net.UDPAddr, err error) {
	if len(inputStrings) < 2 || inputStrings[1] == "-h" {
		help()
	}

	for i, addrString := range os.Args[1:] {
		addr, err := net.ResolveUDPAddr("udp", addrString)
		if err != nil {
			return nil, nil, err
		}

		if i == 0 {
			hostPort = addr
		} else {
			remotePorts = append(remotePorts, addr)
		}
	}

	return
}

func help() {
	log.Println("UDP Packet Duplicator")
	log.Println("Usage: udpPacketDuplicator [:hostPort] [:remotePort]...")

	os.Exit(0)
}
