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
var newElevatorAddress string
var neighbourElevatorAddress string
var neighbourConnection *net.UDPConn

var elev_id int 			// each elevator has an unique ID
var nElevators int	
var IPaddress string // along with port


func singleElevator() {
	if nElevators == 0 && neighbourElevatorAddress == "" {
		nElevators = 1;
	}
	if elev_id == 0 {
		elev_id = nElevators;
	}
}

func newElevDetected(addr string) {
	fmt.Println("ADDRESS: ", addr, newElevatorAddress)
	if elev_id == 0 {
		elev_id = nElevators
	} else if elev_id == 1 && newElevatorAddress != addr{
		nElevators++
		newElevatorAddress = addr
	}
}

func send_msg(msg Message, conn *net.UDPConn){
	buffer,err := json.Marshal(msg)
	if err != nil{
		fmt.Println("ERROR IN MARSHAL OF SENDING message")
		fmt.Println("%s", err)
	}
	
	_,_ = conn.Write(buffer)
}


func SendElevMessages() {
	
	msg := "Hello!"
	addr := GetLocalAddress()
	
	for {
			numElev := nElevators
			newElevAddr := newElevatorAddress
			dial_msg := Message{StringMsg : msg, LocalAddr : addr, NumElev : numElev, NewElevAddr : newElevAddr}
		
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
	numElev := 0
	b := ""
	t0 := time.Now()
	var message Message

	for {	
			n, err := conn.Read(buffer)
			t1 := time.Now()
			if(t1.Sub(t0) > 2*(1000*time.Millisecond)){
				singleElevator()
			}
			if err != nil {
				fmt.Println("ERROR IN READING MESSAGE")
				fmt.Println(err)
			}

			if n != 0 {		
				_ = json.Unmarshal(buffer[:n], &message)
				addr = message.LocalAddr
				numElev = message.NumElev
				if b != message.NewElevAddr {
				   b = message.NewElevAddr
				}
				if numElev > nElevators {
						nElevators = numElev			
				}
				
				if(neighbourElevatorAddress == "") {
					neighbourElevatorAddress = addr
					newElevDetected(neighbourElevatorAddress)
				}
				if(addr != "" && neighbourElevatorAddress != addr) {
					//fmt.Println("WE HAVE DETECTED A NEW ELEVATOR WITH ADDRESS: ", addr)
					newElevDetected(addr)
				}
				if((nElevators-elev_id) == 1 && elev_id != 1) {
					// WE ARE THE LAST ELEVATOR
					fmt.Println("WE ARE THE LAST ELEVATOR")
					neighbourElevatorAddress = newElevatorAddress
					GetNeighbourElevConnection().Close()
					SetNeighbourElevConnection()
				}
				if newElevatorAddress != b && b != "" && b != neighbourElevatorAddress {
						newElevatorAddress = b
				}
				fmt.Println(message.StringMsg, neighbourElevatorAddress, nElevators, elev_id, newElevatorAddress)
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
	
	broadcast_msg := Message{StringMsg : msg, LocalAddr : addr}

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
/*
func Init_newElevator(init_msg initialization_msg){
	elev_id = init_msg.New_id
	nElevators = init_msg.NumberOfElevators
	neighbourIP = init_msg.NextElevatorAddr
	neighbourPORT = init_msg.NextElevatorPort
}*/



func Send_requestedOrderEvaluation(elev_score []float64, floor int, buttontype ButtonType){
	msg := Message{MsgID : OrderRequestEvaluation,Elev_id : elev_id, Elev_score : elev_score,Floor : floor, Buttontype : buttontype}
	send_msg(msg,nil)
	/*buffer,err := json.Marshal(msg)
	if err != nil{
		fmt.Println("ERROR IN MARSHAL OF REQUESTED ORDER")
		fmt.Println("%s", err)
	}*/
	
	//_,_ = neighbourconnection.Write(buffer)
}

func Send_newOrder(floor int, button ButtonType, chosenElevator int, conn *net.UDPConn){	
	msg := Message{MsgID : OrderRequest,Elev_targetID: chosenElevator, Floor : floor, Buttontype : button }
	send_msg(msg, conn)
	/*buffer,err := json.Marshal(msg)
	if err != nil{
		fmt.Println("ERROR IN MARSHAL OF SENDING ORDER")
		fmt.Println("%s", err)
	}
	_,_ = neighbourconnection.Write(buffer)*/
}

func Send_elevInitCompleted(successfull bool){
	msg := Message{MsgID : Elevator_initializationStatus,Elev_targetID : elev_id, SuccessfullInit : successfull}
	send_msg(msg,nil)
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


func Thread_elevatorCommRecv(C_elevatorCommand chan int,C_order chan Ch_elevOrder){//assume it has made a connection for listening. This requires a connection first.
	//PAPPA is joey a faggot?
	conn := GetNeighbourElevConnection()
	buffer := make([]byte, 1024)
	var message Message

	for{
		n,err := conn.Read(buffer)
		if (err != nil){
			fmt.Println("ERROR IN READING ELEVATOR MESSAGE")
		fmt.Println(err)
		}
		if(n != 0){
			fmt.Println("Message received ! ")
			_ = json.Unmarshal(buffer[:n],&message)
		
			switch(message.MsgID){	
			case OrderRequestEvaluation:
				channel_message := Ch_elevOrder{Floor : message.Floor, Button : message.Buttontype, Elev_score : message.Elev_score} 
				C_elevatorCommand <- 1
				C_order <- channel_message
				//C_order.Floor <- message.Floor
				//C_order.Buttontype <-message.Buttontype
				//C_order.Elev_score <- message.Elev_score //ask studass
			case OrderRequest:
				if(message.Elev_targetID == elev_id){
					C_elevatorCommand <- 2
					channel_message := Ch_elevOrder{Floor : message.Floor, Button : message.Buttontype}
					C_order <- channel_message
					//C_order.Floor <- message.Floor
					//C_order.Buttontype <- message.Buttontype //once again ask studass
				}else{
					Send_newOrder(message.Floor, message.Buttontype, message.Elev_targetID,nil)	
				}
			case Elevator_initializationStatus:
					//something.
				if(message.Elev_id ==GetID()){
					//This elevator was the one initializing the command for status.
					if(message.SuccessfullInit == true){
						//All elevators are in working order fuck yeah! send nothing.
					}else{
						return
						//Not all elevators are in working order. Must send a message to every elevator with the failed target id. If it doesnt work
						//after a second
						//message.Elev_failedID 

					}
				}					
			}
		}
	}
}


