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
	NewElevatorconnection
	Elevator_initializationStatus
)	 		

type Message struct{
	//First define the type and elev_id which is always relevant
	MsgID msgType
	Elev_id int
	StringMsg string
	LocalAddr string
	NumElev int

	//initalization_msg
	New_id int // Giving a new ID to a new elevator
	NumberOfElevators int
	NextElevatorAddr string
	NextElevatorPort string

	//general Elev msg variables:
	Floor int
	Buttontype ButtonType
	//order Request Evaluation
	Elev_score []float64
	//Floor int
	//Buttontype buttonType
	//msg order request

	Elev_targetID int
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

