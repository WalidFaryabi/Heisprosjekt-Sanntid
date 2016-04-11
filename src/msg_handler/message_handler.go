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

func Send_requestedOrder([]elev_score int, floor int, buttontype buttonType){
	msg := msg_orderRequest{elev_id: elev_id, []elev_score : []elev_score,floor : floor, buttonType : buttontype }
	buffer,err := json.Marshal(msg)
	if err != nil{
		fmt.Println("ERROR IN MARSHAL OF REQUESTED ORDER")
		fmt.Println("%s", err)
	}
	_,_ = neighbourconnection.Write(buffer)
}

func receive_


func send_initMsg(init_msg initialization_msg){


}


func Init_newElevator(init_msg initialization_msg){
	elev_id = init_msg.new_id
	nElevators = init_msg.numberOfElevators
	neighbourIP = init_msg.nextElevatorAddr
	neighbourPORT = init_msg.nextElevatorPort
}



type initialization_msg struct {
	new_id int // Giving a new ID to a new elevator
	numberOfElevators int
	nextElevatorAddr string
	nextElevatorPORT string
}

type msg_orderRequest struct {
	elev_id int
	[]elev_score int
	floor int
	buttontype buttonType
}

func GetID()(int){
	return elev_id
}

func GetNelevators()(int){
	return nElevators
}
