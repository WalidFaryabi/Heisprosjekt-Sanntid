package network

import (
	"fmt"
	"net"
	"os"
	"encoding/json"
	"strconv"
)

func checkAndPrintError(err error, description string) {
	if err != nil {
		fmt.Println(err)
		fmt.Printf(description)
	}
}

func GetUDPAddress(ipAddress string, port int) (*net.UDPAddr) {
	port := "20000"+strconv.Itoa(port)
	address := ipAddress+":"+port
	
	addr, err := net.ResolveUDPAddr("udp", address)
	checkAndPrintError(err, "Unsuccessful udpAdress retrieval")
	return addr
}

func GetConnectionForListening(udpAddr *UDPAddr) (*net.UDPConn) {
	conn, err := net.ListenUDP("udp", udpAddr)
	checkAndPrintError(err, "Error occurred in establishing connection for listening")
	return conn
}

func GetConnectionForDialing(udpAddr *UDPAddr) (*net.UDPConn) {
	conn, err := net.DialUDP("udp", nil, udpAddr)
	checkAndPrintError(err, "Error occurred in establishing a connection for sending")
	return conn	
}

