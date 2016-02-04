package main

import (
	"fmt"
	"net"
	"os"
	"encoding/json" 
)

type Message struct {
	Msg string
    	Identifier int64
}


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

func send_message(message Message, remoteSocket *net.UDPConn){
	buffer, err := json.Marshal(message)
	checkError(err, "error in json")
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
func read(socket *net.UDPConn) {

	//addr, err := net.ResolveUDPAddr("udp4", ":20019")
	//checkError(err,"reading error")
	//udpConn, err := net.ListenUDP("udp4", addr)
	//checkError(err, "no read sir")
	for {
		b := make([]byte,1024)	
		n,_ := socket.Read(b)
		s := string(b[:n])
		fmt.Println(s)
	}
}

func main(){
	ms := Message{"Alice", 1241241444}
	var m Message
	b, _ := json.Marshal(ms)
	_ = json.Unmarshal(b,&m)
	fmt.Println(m.Identifier)
	/*ms := Message{420, "Moren til Johann er hans mor"}
	fmt.Println(ms.identifier, ms.msg)
	buffer, err := json.Marshal(ms)
	checkError(err, "error in json")
	fmt.Println(buffer)
	var msgd Message
	err11 := json.Unmarshal(buffer, &msgd)
	checkError(err11,"u done fucke up")
	fmt.Println(msgd.identifier, msgd.msg)*/


}

/*func main() {
	msg := string("what do you mean?")
	
	addr := get_serverAddress("udp4", "129.241.187.255", "20019")
	remote_socket := createRemoteSocket(addr, "udp4")
	local_socket := createLocalSocket("udp4", "", "20019")
	
	go read(local_socket)
	go send_message(msg, remote_socket)
 
	
	select {}
}*/
