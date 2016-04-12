package msg_handler

import(
	"network"
	"encoding/json"	
)

var LocalIP	string
var Localport string
//var neighourIP string
//var neighbourPORT string
var neighbourconnection network.UDPCONN

var elev_id int 			// each elevator has an unique ID
var nElevators int 	
var IPaddress string // along with port
type buttonType int


func broadcast(conn *net.UDPConn) {
	t0 := time.Now()
	run_bc := true
	msg := "Hello!"
	addr := netw.GetLocalIP()+":"+strconv.Itoa(20000+listenPort)
	
	broadcast_msg := Message{msg, addr}

	for {
		if(run_bc){
			buffer,err := json.Marshal(broadcast_msg)
	
			if err != nil {
				fmt.Println("ERROR IN MARSHAL")
				fmt.Println("%s", err)
			}

			_,_ = conn.Write(buffer)

			t1 := time.Now()
			if(t1.Sub(t0) > 2*(1000*time.Millisecond)){
				run_bc = false
			}
		} else {
			break
		}
		
	}
	conn.Close()
}



func receive_


func send_initMsg(init_msg initialization_msg){


}




///////**************************////////////////////////

func Init_newElevator(init_msg initialization_msg){
	elev_id = init_msg.New_id
	nElevators = init_msg.NumberOfElevators
	neighbourIP = init_msg.NextElevatorAddr
	neighbourPORT = init_msg.NextElevatorPort
}

func Send_requestedOrderEvaluation([]elev_score int, floor int, buttontype buttonType){
	msg := msg_orderRequestEvaluation{MsgID : OrderRequestEvaluation,Elev_id: elev_id, []Elev_score : []elev_score,Floor : floor, ButtonType : buttontype}
	buffer,err := json.Marshal(msg)
	if err != nil{
		fmt.Println("ERROR IN MARSHAL OF REQUESTED ORDER")
		fmt.Println("%s", err)
	}
	_,_ = neighbourconnection.Write(buffer)
}

func Send_newOrder(floor int, button ButtonType, chosenElevator int){	
	msg := msg_OrderRequest{MsgID : OrderRequest,Elev_targetID : chosenElevator, Floor : floor, Buttontype : button }
	buffer,err := json.Marshal(msg)
	if err != nil{
		fmt.Println("ERROR IN MARSHAL OF SENDING ORDER")
		fmt.Println("%s", err)
	}
	_,_ = neighbourconnection.Write(buffer)
}

func Send_elevInitCompleted(successfull bool){
	msg := msg_elevInit{MsgID : Elevator_initializationStatus,Elev_targetID : elev_id, SuccessfullInit : successfull}
	buffer,err := json.Marshal(msg)
	if err != nil{
		fmt.Println("ERROR IN MARSHAL OF SENDING ORDER")
		fmt.Println("%s", err)
	}
	_,_ = neighbourconnection.Write(buffer)
}
}



func GetID()(int){
	return elev_id
}

func GetNelevators()(int){
	return nElevators
}

type msg_OrderRequest struct{
	Elev_targetID int
	Floor int
	Buttontype buttonType
}