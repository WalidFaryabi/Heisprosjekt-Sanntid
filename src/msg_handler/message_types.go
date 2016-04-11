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
	OrderRequest
	newElevatorconnections
)	 		


type messageID int
const(

)

type initialization_msg struct {
	new_id int // Giving a new ID to a new elevator
	numberOfElevators int
	nextElevatorAddr string
	nextElevatorPort string
}



type msg_orderRequest struct {
	elev_id int
	[]elev_score int
	floor int
	buttontype buttonType
}



/*type messageType struct{
	
	msg_type msgType


	//initalization info:
	new_id int
	numberOfElevators int
	nextElevatorAddr string
	nextElevatorPORT string

	//OrderRequest
	elev_id int
	[]elev_score int
	floor int
	buttontype buttonType


}*/