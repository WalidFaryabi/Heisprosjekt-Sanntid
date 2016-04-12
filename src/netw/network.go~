package netw

import (
	"fmt"
	"net"
	"bufio"
	"strconv"
	"os"
	"strings"
)

func checkAndPrintError(err error, description string) {
	if err != nil {
		fmt.Println(err)
		fmt.Println(description)
	}
}

func GetConnectionForListening(address string) (*net.UDPConn) {

	udpAddr, err := net.ResolveUDPAddr("udp", address)
	checkAndPrintError(err, "Unsuccessful udpAdress retrieval")
	
	conn, err := net.ListenUDP("udp4", udpAddr)
	checkAndPrintError(err, "Error occurred in establishing connection for listening")
	return conn
}


func GetConnectionForDialing(address string) (*net.UDPConn) {
	
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	checkAndPrintError(err, "Unsuccessful udpAdress retrieval \n")

	conn, err := net.DialUDP("udp", nil, udpAddr)
	checkAndPrintError(err, "Error occurred in establishing a connection for sending")
	return conn	
}

func GetLocalIP() (string) {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}

func GetPort() (int) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter listen port: ")
	portString, _ := reader.ReadString('\n')
	portString = strings.Replace(portString, "\n", "", -1)
	port,err := strconv.Atoi(portString)

	if err != nil {
		fmt.Println("ERROR IN CONVERSION")
		fmt.Printf("%s \n", err)
		return 0
	}
	return port
}

/*
func GetUDPAddress(ipAddress string, port int) (*net.UDPAddr) {
	var address string
	if port != 0 {
		p := 20000+port
		address = ipAddress+":"+strconv.Itoa(p)
	} else {
		address = ipAddress+":"
	}

	addr, err := net.ResolveUDPAddr("udp4", address)
	checkAndPrintError(err, "Unsuccessful udpAdress retrieval \n")
	return addr
}

func GetConnectionForListening(udpAddr *net.UDPAddr) (*net.UDPConn) {
	conn, err := net.ListenUDP("udp4", udpAddr)
	checkAndPrintError(err, "Error occurred in establishing connection for listening \n")
	return conn
}

func GetConnectionForDialing(udpAddr *net.UDPAddr) (*net.UDPConn) {
	conn, err := net.DialUDP("udp", nil, udpAddr)
	checkAndPrintError(err, "Error occurred in establishing a connection for sending \n")
	return conn	
} */

