

package main

import (
	"net"
	//"os"
	."fmt"
	//"strconv"
)

func LookUpHost(name string){
	addrs, err := net.LookupHost(name)
	if err !=nil{
		Println("Error ", err.Error())
		return		
		}
	for _, s := range addrs{
		Println(s)
	}
	
}

/*func createServerConnectionTCP(port string)(connection *net.Conn){
	ln, err := net.Listen("tcp", port)
	if err != nil{
		fmt.Println("Error", err.Error())	
		return	
		}
	for{
		connection, err := ln.Accept()
		if err != nil{
			fmt.Println("ERROR", err.Error())
		}
	}

} */

func ResolveUDP(port string)(udpaddr *net.UDPAddr){
	resolve, err := net.ResolveUDPAddr("udp", port)
	if err != nil{
		Println("error", err.Error())
		return
		}
	return resolve
	
}

func listenUDP(addr *net.UDPAddr)(Conn *net.UDPConn){
	Conn, err := net.ListenUDP("udp4", addr)
	if err != nil{
		Println("anothah error", err.Error())
		return
	}
	return Conn
}
func DialTheUDP(localAddr *net.UDPAddr,addr *net.UDPAddr)(conn *net.UDPConn){
	conn,err := net.DialUDP("udp4", localAddr, addr)
	if(err != nil){
		Println("Sorry you got an error", err.Error())
		return	
	}
	return conn
}

func sendMessageToServer(addr *net.UDPAddr, conn *net.UDPConn){
	a := make([]byte,1024)
	i :=0
	for j := 'a'; j<'z'; j++{
		a[i] = byte(j)
		i++ 
	}
	_,err := conn.Write(a)
	if err != nil{
		Println("ERROR", err.Error())	
	}
	
}


/*func sender(c chan string, ip string, port int){
	//addr, err := net.ResolveUDPAddr("udp4", ip+":"+strconv.Itoa(port))
	if err != nil { Println(err) }
	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil { Println(err) }

	for {
		conn.Write([]byte(<-c))
	}
}*/

func readMessageFromServer(addr *net.UDPAddr/*conn *net.UDPConn*/){

	//addr, err := net.ResolveUDPAddr("udp4", ":"+strconv.Itoa(port))
	//if err != nil { Println(err) }
	//conn, err := net.ListenUDP("udp4", addr)
	//if err != nil { Println(err) }
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil { Println(err) }


	for {
		b := make([]byte,1024)	
		n,_,_ := conn.ReadFromUDP(b)
		s := string(b[:n])
		Println("Received:", s)
	}
}

/*
func readMessageFromServer(port int){

	addr, err := net.ResolveUDPAddr("udp4", ":"+strconv.Itoa(port))
	if err != nil { Println(err) }
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil { Println(err) }

	for {
		b := make([]byte,1024)	
		n,_,_ := conn.ReadFromUDP(b)
		s := string(b[:n])
		Println("Received:", s)
	}
}*/



func main() {

	var localAddress *net.UDPAddr
	var serverAddress *net.UDPAddr
	var udpConn *net.UDPConn

	localAddress = ResolveUDP("129.241.187.255:0")
	serverAddress = ResolveUDP("129.241.187.23:20021")
	udpConn = DialTheUDP(localAddress, serverAddress)
	sendMessageToServer(serverAddress, udpConn)
	Println("READING MESSAGE")
	readMessageFromServer(serverAddress)

	/*c := make(chan string)
	go sender(c, "129.241.187.255", 20021)
	go readMessageFromServer(20021)
	go readMessageFromServer(30000)

	for j := 'a'; j<'z'; j++ {
		c <- string(j)
	}*/ 

	/*a := make([]byte, 1024)
	var udpaddr *net.UDPAddr
	var udpConn *net.UDPConn
	udpaddr = ResolveUDP(":30000")
	udpConn = listenUDP(udpaddr)
	var size int
	size, udpaddr, _ = udpConn.ReadFromUDP(a)
	Println(size)
	s := string(a)
	Println(s)
	Println(udpaddr.Network())*/

	select {}

}

























