package main

import (
	"fmt"
	"net"
	"os"
)


func checkError(err error, errorpart string) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func get_serverAddress(udpType, ipAddress,port string)(addr *net.UDPAddr){
	addr, err := net.ResolveUDPAddr(udpType, ipAddress + ":" + port)
	checkError(err, ": ServerAddress receival")
	return addr		
}

func createRemoteSocket(addr *net.UDPAddr, udpType string)(*net.UDPConn){
	remoteSocket, err := net.DialUDP(udpType, nil, addr)
	checkError(err, ": Unsuccessfull socket connection")
	return remoteSocket	 
}

func createLocalSocket(udpType, ipAddress,port string)(*net.UDPConn){
	addr, err := net.ResolveUDPAddr(udpType, ipAddress +":" + port)
	checkError(err,": Unsuccessfull retrieval of address")
	conn, err := net.ListenUDP(udpType, addr)
	checkError(err, "Unsuccessfull connection to local socket")
	return conn
}

func send_message(message string, remoteSocket *net.UDPConn){
	buffer := make([]byte, 1024)
	for i:=0; i<len(message); i++ {
		buffer[i] = byte(message[i])
	}
	for {
		_, err := remoteSocket.Write(buffer)
		checkError(err, "error in send message")
	}
}



func recv_message(port string, localSocket *net.UDPConn)(int, []byte){
	buffer := make([]byte,1024)	
	size, _ := localSocket.Read(buffer)
	return size, buffer
}


func read_message(buffer[] byte, size int){
	for i := 0; i<size; i++{
		s := string(buffer[size])
		fmt.Println(s)
	}
}

func main() {
	msg := string("Hello, I am from the other computer. I found out about our existence")

	addr := get_serverAddress("udp4", "129.241.187.157", "20019")
	remote_socket := createRemoteSocket(addr, "udp4")
	
	go send_message(msg, remote_socket)
	
	select {}
}
