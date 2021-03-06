package main

import (
	"crypto/rand"
	"fmt"
	"net"
	"testing"
)

func TestBasicSingleTargetProxy(t *testing.T) {
	ports, err := toPorts([]string{":8000", ":8001"})
	if err != nil {
		t.Fatal(err)
	}

	hostConn, err := net.ListenUDP("udp", ports[0])
	remoteConn1, err := net.ListenUDP("udp", ports[1])

	// start the packet duplicator
	go proxy(hostConn, ports[1:])

	// send a test packet
	testString := "hello moto"
	hostConn.WriteToUDP([]byte(testString), ports[0])

	// check the proxy target recieved the packet
	testBuffer1 := make([]byte, 100)
	bytes, _, err := remoteConn1.ReadFromUDP(testBuffer1)
	if err != nil {
		t.Error(err)
	}
	if string(testBuffer1[:bytes]) != testString {
		t.Error("Target endpoint 1 did not recieve correct payload",
			"\nexpected:", testString,
			"\nrecieved:", string(testBuffer1[:bytes]))
	}
}

func TestBasicDoubleTargetProxy(t *testing.T) {
	ports, err := toPorts([]string{":8002", ":8003", ":8004"})
	if err != nil {
		t.Fatal(err)
	}

	hostConn, err := net.ListenUDP("udp", ports[0])
	remoteConn1, err := net.ListenUDP("udp", ports[1])
	remoteConn2, err := net.ListenUDP("udp", ports[2])

	// start the packet duplicator
	go proxy(hostConn, ports[1:])

	// send a test packet
	testString := "hello moto"
	hostConn.WriteToUDP([]byte(testString), ports[0])

	// check the proxy targets recieved the packet
	testBuffer1 := make([]byte, 100)
	bytes, _, err := remoteConn1.ReadFromUDP(testBuffer1)
	if err != nil {
		t.Error(err)
	}
	if string(testBuffer1[:bytes]) != testString {
		t.Error("Target endpoint 1 did not recieve correct payload",
			"\nexpected:", testString,
			"\nrecieved:", string(testBuffer1[:bytes]))
	}

	testBuffer2 := make([]byte, 100)
	bytes, _, err = remoteConn2.ReadFromUDP(testBuffer2)
	if err != nil {
		t.Error(err)
	}
	if string(testBuffer2[:bytes]) != testString {
		t.Error("Target endpoint 2 did not recieve correct payload",
			"\nexpected:", testString,
			"\nrecieved:", string(testBuffer1[:bytes]))
	}
}

func BenchmarkProxy(b *testing.B) {
	ports, err := toPorts([]string{":8005", ":8006", ":8007"})
	if err != nil {
		b.Fatal(err)
	}

	hostConn, err := net.ListenUDP("udp", ports[0])

	// send a test packet
	testByte := make([]byte, 1024)
	rand.Read(testByte)
	fmt.Println(testByte)

	buffer := make([]byte, bufferSize)

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		hostConn.WriteToUDP(testByte, ports[0])
		proxyPacket(hostConn, ports[1:], buffer)
	}
}
