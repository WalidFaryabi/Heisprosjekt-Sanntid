package network

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

func createRemoteSocket(addr *net.UDPAddr udpType)(*net.UDPConn){
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

func send_messageBuffer(message string, remoteSocket *net.UDPConn){
	buffer := make([]byte, 1024)
	for i:=0; i<len(message); i++ {
		buffer[i] = byte(message[i])
	}
	for {
		_, err := remoteSocket.Write(buffer)
		checkError(err, "error write")
	}
}


func send_message(message string, remoteSocket *net.UDPConn){
	buffer := make([]byte, 1024)
	for i:=0; i<len(message); i++ {
		buffer[i] = byte(message[i])
	}
	for {
		_, err := remoteSocket.Write(buffer)
		checkError(err, "error write")
	}
}



func recv_message(port string, localSocket *net.UPDConn)(size, buffer[] byte){
	buffer := make([]byte,1024)	
	size, _ := udpConn.Read(buffer)
	return size, buffer
}



func read(localSocket *net.UDPConn) {
	for {
		b := make([]byte,1024)	
		n,_ := localSocket.Read(b)
		s := string(b[:n])
		fmt.Println(s)
	}
}

