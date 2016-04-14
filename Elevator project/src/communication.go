
package main
import ( "fmt"
		"./netw"
		"./msg_handler"
		"./FSM"
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

type gos int

func blablo(test gos){
	fmt.Println("TEST")
	fmt.Println(test)
}

func main() {
	
	//FSM.Event_init()
	//for{}
	stfu_joey := make(chan msg_handler.Ch_elevOrder)
	stfu_joey_pls := make(chan int)
	go FSM.Thread_elevatorStateMachine(stfu_joey_pls, stfu_joey)
	for{
	}

	init_localAddress()
	msg_handler.Broadcast()
	
	go msg_handler.ListenForElevMessages()	
    
	waitForNeighbourElevAddr()
	
	msg_handler.SetNeighbourElevConnection()
	
	go msg_handler.SendElevMessages()
	
    var input string
    fmt.Scanln(&input)
}


