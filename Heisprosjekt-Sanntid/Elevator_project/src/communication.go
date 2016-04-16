
package main
import ( "fmt"
		"./netw"
		"./msg_handler"
		"./FSM"
		"time"
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
	//FSM.SetNFloors()
	//FSM.Event_init()
	//for{}
	msg_handler.SemaphoreMessage <- 1
	msg_handler.SemaphoreRead <-1
	stfu_joey := make(chan msg_handler.Ch_elevOrder,10)
	
	
	stfu_joey_pls := make(chan int,10)
	go FSM.Thread_elevatorStateMachine(stfu_joey_pls,stfu_joey)
	for{}
	
	C_messages := make(chan msg_handler.Message,10)

	msg_handler.InitElevatorNetwork()	
		//C_sendCommando chan int, C_message chan Message,C_elevatorCommand chan int,C_order chan Ch_elevOrder 
	go msg_handler.ReceiveElevMessages(C_messages)

	time.Sleep(10 * time.Second)
	fmt.Println("ready for sending")
	go msg_handler.SendElevMessages(C_messages)	//SendElevMessages(C_listenCommando chan int, C_message chan Message, C_elevatorCommand chan int,C_order chan Ch_elevOrder)
	go FSM.Thread_elevatorStateMachine(stfu_joey_pls,stfu_joey)
	time.Sleep(10 * time.Second)
	fmt.Println("Elevator initialized.")

	go msg_handler.Thread_elevatorCommRecv(stfu_joey_pls, stfu_joey)
	go FSM.Thread_elevatorStateMachine(stfu_joey_pls, stfu_joey)
	fmt.Printf("elev id %i \n", msg_handler.GetID())
	for{}
	go FSM.TestElevator()
	select{

	}
	for{}
	//go msg_handler.Thread_elevatorCommRecv(stfu_joey_pls, stfu_joey)

	

	init_localAddress()
	//msg_handler.Broadcast()
	go msg_handler.Thread_elevatorCommRecv(stfu_joey_pls, stfu_joey)
	go FSM.Thread_elevatorStateMachine(stfu_joey_pls, stfu_joey)
	//go msg_handler.ListenForElevMessages()	
    
	waitForNeighbourElevAddr()
	
	//msg_handler.SetNeighbourElevConnection() PRIVATE function not global.
	

	//go msg_handler.SendElevMessages()
	
    var input string
    fmt.Scanln(&input)
}


