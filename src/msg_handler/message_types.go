package msg_handler

import(
	//"fmt"
)

var elev_id int 			// each elevator has an unique ID
var nElevators int 	
var IPaddress string // along with port
type buttonType int
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

	//initalization_msg
	New_id int // Giving a new ID to a new elevator
	NumberOfElevators int
	NextElevatorAddr string
	NextElevatorPort string

	//general Elev msg variables:
	Floor int
	Buttontype buttonType
	//order Request Evaluation
	[]Elev_score int
	//Floor int
	//Buttontype buttonType
	//msg order request

	Elev_targetID int
	//Floor int
	Buttontype buttonType

	//elev_init
	Elev_targetID int
	SuccesfullInit bool
	Elev_failedID int

}


type ch_elevOrder struct{
	Floor int
	Button buttonType


}

