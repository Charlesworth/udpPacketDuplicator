package main

import (
	"errors"
	"net"
	"os"
)

func getPorts(inputStrings []string) (hostPort *net.UDPAddr, remotePorts []*net.UDPAddr, err error) {
	if len(inputStrings) < 2 {
		return nil, nil, errors.New("Require at least 2 ports as input arguements")
	}

	for i, addrString := range os.Args[1:] {
		addr, err := net.ResolveUDPAddr("udp", addrString)
		if err != nil {
			return nil, nil, err
		}

		if i == 1 {
			hostPort = addr
		} else {
			remotePorts = append(remotePorts, addr)
		}
	}

	return
}
