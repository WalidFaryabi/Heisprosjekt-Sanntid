package msg_handler

import(
	//"fmt"
)


type ButtonType int
const(
	BUTTON_UP = iota
	BUTTON_DOWN
	BUTTON_COMMAND
)

type msgType int
const(
	Initialization = iota
	InitNeighbourElevConn
	OrderRequestEvaluation
	OrderRequest
	NewElevatorConnection
	Elevator_initializationStatus
	BroadcastMsg
	NewElevatorConnectionEstablished
	NewElevatorInitConfig
	BroadcastAcknowledged
	Debug
	ExternalOrderComplete
)	 		

type Message struct{
	//First define the type and elev_id which is always relevant
	MsgID msgType
	Elev_id int
	StringMsg string
	LocalAddr string	//Local address of the elevator sending this message
	NumElev int 		//Number of elevators in the system.
	LocalAddrOfFirstElevator string // the address of the elevator with elev_id = 1	// Johann this is implicitily given already by setting Elev_id = 1, dont see the point

	//initalization_msg
	New_id int // Giving a new ID to a new elevator
	NumberOfElevators int
	NextElevatorAddr string 	//for the previous last elevator, such that it connects correctly
	NextElevatorPort string

	//new elevator connection
	NewElevatorLocalAddress string


	//general Elev msg variables:
	Floor int
	Buttontype ButtonType
	//order Request Evaluation
	Elev_score []float64

	//Floor int
	//Buttontype buttonType
	//msg order request

	Elev_targetID int
	TargetID int //no difference just want a better variable
	//Floor int

	//elev_init
	
	SuccessfullInit bool
	Elev_failedID int
}

type Ch_elevOrder struct{
	Floor int
	Button ButtonType
	Elev_score []float64
	Elev_id int
}

var commandTypes int
const(
	NEW_ELEVATOR_CONFIG int = iota

)