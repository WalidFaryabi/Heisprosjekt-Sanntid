package msg_handler


import(
	"fmt"
	"net"
	"../netw"
	"encoding/json"	
	"strconv"
	"time"
)	



const broadcastIP string = "129.241.187.255"
const broadcastPort string = "20022"

var LocalIP	string
var LocalPort int
var neighbourElevatorAddress string
var neighbourConnection *net.UDPConn

var elev_id int 			// each elevator has an unique ID
var nElevators int 	
var IPaddress string // along with port



func send_msg(msg Message){
	buffer,err := json.Marshal(msg)
	if err != nil{
		fmt.Println("ERROR IN MARSHAL OF SENDING message")
		fmt.Println("%s", err)
	}
	_,_ = neighbourConnection.Write(buffer)

	/*switch(msg.MsgID){
		case Initialization:
			//joey pls insert
		//case Elevator joey pls
		case OrderRequestEvaluation:

	}*/
}


func SendElevMessages() {

	msg := "Hello!"
	addr := GetLocalAddress()
	
	for {
			dial_msg := Message{StringMsg : msg, LocalAddr : addr}
			buffer,err := json.Marshal(dial_msg)
	
			if err != nil {
				fmt.Println("ERROR IN MARSHAL")
				fmt.Println("%s", err)
			}
			_,err = neighbourConnection.Write(buffer)

			if err != nil {
				fmt.Println("ERROR IN DIALING")
				fmt.Println(err)		
			}
		}
}

func listen(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	addr := ""
	t0 := time.Now()
	var message Message

	for {	
			if (neighbourElevatorAddress == ""){
				n, err := conn.Read(buffer)
				t1 := time.Now()
				if(t1.Sub(t0) > 2*(1000*time.Millisecond)){
					nElevators = 1
					fmt.Println("NUMBER OF ELEVATORS: 1")
				}
				
				if err != nil {
					fmt.Println("ERROR IN READING MESSAGE")
					fmt.Println(err)
				}

				if n != 0 {
					_ = json.Unmarshal(buffer[:n], &message)
					fmt.Println(message.stringMsg, message.LocalAddr[len(message.LocalAddr)-2:])
					SetNeighbourElevatorAddress(message.LocalAddr)
				}

			} else {
				n, err := conn.Read(buffer)
				if err != nil {
					fmt.Println("ERROR IN READING MESSAGE")
					fmt.Println(err)
				}
				if n != 0 {
					_ = json.Unmarshal(buffer[:n], &message)
					
					fmt.Println(message.StringMsg, message.LocalAddr[len(message.LocalAddr)-2:])
					addr = message.LocalAddr
					if (addr != "" && neighbourElevatorAddress != addr) {
						fmt.Println("WE HAVE DETECTED A NEW ELEVATOR WITH ADDRESS: ", addr)
					}
				}		
			}
		}
}

func ListenForElevMessages() {
	listenIP := ""
	listenAddr := listenIP + ":" + strconv.Itoa(20000+LocalPort) 
	listen_conn := netw.GetConnectionForListening(listenAddr)
	
	listen(listen_conn)	
}

func broadcast(conn *net.UDPConn) {
	t0 := time.Now()
	run_bc := true
	msg := "Hello!"
	addr := GetLocalAddress()
	
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

func Broadcast() {
	broadcastAddr := broadcastIP+":"+broadcastPort
	brdcast_conn := netw.GetConnectionForDialing(broadcastAddr)
	
	broadcast(brdcast_conn)
}

func IsNeighbourElevatorAddressEmtpy()(bool) {
	if (neighbourElevatorAddress != "") {
		return false
	}
	return true
}


///////**************************////////////////////////

func Init_newElevator(init_msg initialization_msg){
	elev_id = init_msg.New_id
	nElevators = init_msg.NumberOfElevators
	neighbourIP = init_msg.NextElevatorAddr
	neighbourPORT = init_msg.NextElevatorPort
}

func Send_requestedOrderEvaluation(elev_score []int, floor int, buttontype buttonType){
	msg := Message{MsgID : OrderRequestEvaluation,Elev_id : elev_id, Elev_score : elev_score,Floor : floor, ButtonType : buttontype}
	send_msg(msg)
	/*buffer,err := json.Marshal(msg)
	if err != nil{
		fmt.Println("ERROR IN MARSHAL OF REQUESTED ORDER")
		fmt.Println("%s", err)
	}*/
	
	//_,_ = neighbourconnection.Write(buffer)
}

func Send_newOrder(floor int, button ButtonType, chosenElevator int){	
	msg := Message{MsgID : OrderRequest,Elev_targetID : chosenElevator, Floor : floor, Buttontype : button }
	send_msg(msg)
	/*buffer,err := json.Marshal(msg)
	if err != nil{
		fmt.Println("ERROR IN MARSHAL OF SENDING ORDER")
		fmt.Println("%s", err)
	}
	_,_ = neighbourconnection.Write(buffer)*/
}

func Send_elevInitCompleted(successfull bool){
	msg := Message{MsgID : Elevator_initializationStatus,Elev_targetID : elev_id, SuccessfullInit : successfull}
	send_msg(msg)
	/*buffer,err := json.Marshal(msg)
	if err != nil{
		fmt.Println("ERROR IN MARSHAL OF SENDING ORDER")
		fmt.Println("%s", err)
	}
	_,_ = neighbourconnection.Write(buffer)*/
}

func SetNeighbourElevatorAddress(address string) {
	if neighbourElevatorAddress != address {
		neighbourElevatorAddress = address
	}
}


func GetNeighbourElevAddress()(string) {
	return neighbourElevatorAddress
}

func SetNeighbourElevConnection() {
	if neighbourConnection == nil && neighbourElevatorAddress != "" {
		neighbourConnection = netw.GetConnectionForDialing(neighbourElevatorAddress)
	}
}

func GetNeighbourElevConnection() (*net.UDPConn) {
	return neighbourConnection
}

func GetLocalAddress() (string) {
	return LocalIP+":"+strconv.Itoa(20000+LocalPort)
}

func GetID()(int){
	return elev_id
}

func GetNelevators()(int){
	return nElevators
}



