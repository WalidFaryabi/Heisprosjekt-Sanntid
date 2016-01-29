package main

import (
	"fmt"
	"net"
	"os"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func addMessageToBuffer(msg string, buffer []byte) {
	for i:=0; i< len(msg); i++ {
		buffer[i] = byte(msg[i])
	}
}

func readBuffer(size int, buffer []byte) {
	fmt.Println(string(buffer[:size]))
}

func write(buffer []byte) {
	addr, err := net.ResolveUDPAddr("udp4", "129.241.187.153:20021")
	checkError(err)
	conn, err := net.DialUDP("udp4", nil, addr)
	for {
		_, err := conn.Write(buffer)
		checkError(err)
	}
}
func read() {

	addr, err := net.ResolveUDPAddr("udp4", ":20021")
	checkError(err)
	udpConn, err := net.ListenUDP("udp4", addr)
	checkError(err)
	for {
		b := make([]byte,1024)	
		n,_ := udpConn.Read(b)
		s := string(b[:n])
		fmt.Println(s)
	}
}


func  main() {
	//find the ip address and use this instead
	
	buffer := make([]byte, 1024)
//	bufferListen := make([]byte, 1024)
	
	msg := string("Hello!")
	addMessageToBuffer(msg, buffer)

	go write(buffer)
	go read()

	select {}
}

/*Receiving from the server */
/*
func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", ":30000")
	checkError(err, "FIRST")
	udpConn, err := net.ListenUDP("udp", udpAddr)
	checkError(err, "SECOND")
	
	buffer := make([]byte, 1024)

	for {
		n, err := udpConn.Read(buffer)
		fmt.Print("Received: ", string(buffer[0:n]))
		checkError(err, "THIRD")
	}
}
 */

