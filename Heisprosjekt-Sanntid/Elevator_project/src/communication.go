
package main
import ( "fmt"
	//	"./netw"
		"./msg_handler"
		"./FSM"
		"time"
)

func main() {
	msg_handler.SemaphoreMessage <- 1
	msg_handler.SemaphoreRead <-1
	msg_handler.SemaphoreNewConnection <- 1
	
	C_elevatorInfoContainer := make(chan msg_handler.Ch_elevOrder,10)
	
	C_elevatorOrders := make(chan int,10)

	C_messages := make(chan msg_handler.Message,10)
	
	FSM.InitElevator()
	msg_handler.InitElevatorNetwork()	

	go msg_handler.Task_broadcastSupervisor()
	go msg_handler.Task_receiveElevMessages(C_messages,C_elevatorOrders, C_elevatorInfoContainer)

	go msg_handler.Task_sendElevMessages(C_messages)	
	go msg_handler.Task_networkSupervisor()
	go FSM.Thread_elevatorStateMachine(C_elevatorOrders, C_elevatorInfoContainer)
	time.Sleep(10 * time.Second)
	
	fmt.Println(msg_handler.GetID())
	for{
		elev_id := msg_handler.NumberUserInput("elev id: ")
		msg_handler.Send_debug("ye hear me nuggah?",elev_id)
		fmt.Printf("ELEV ID: %i \n", msg_handler.GetID())
	}
	select{

	}
	for{}

    var input string
    fmt.Scanln(&input)
}


