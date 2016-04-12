
package main
import ("fmt"
		"./netw"
		"./msg_handler"
)
func waitForNeighbourElevAddr() {
	for {
		if(!msg_handler.IsNeighbourElevatorAddressEmtpy()) {
			break		
		}
	}
}

func init_localAddress() {
	msg_handler.LocalIP = netw.GetLocalIP()
	msg_handler.LocalPort = netw.GetPort()
}

func main() {
	init_localAddress()
	msg_handler.Broadcast()
	
	go msg_handler.ListenForElevMessages()	
    
	waitForNeighbourElevAddr()
	
	msg_handler.SetNeighbourElevConnection()
	
	go msg_handler.SendElevMessages()
	
    var input string
    fmt.Scanln(&input)
}


