package msg_handler


import(
	"fmt"
	"net"
	"../network"
	"encoding/json"	
	"strconv"
	"time"
	"bufio"
	"os"
	"strings"
)	
/***************************************************************************************************************************\\\\\\\\
This package deals with communication between different elevators, handling messages. There are goroutines made for 
receiving messages, sending messages and broadcasting for connection. Please refer to system documentation for more information.

In this file, functions are defined for the msg_handler package.

*****************************************************************************************************************************//////


const broadcastIP string = "129.241.187.255" //using this as a basis for broadcasting.

/******************************************private variables*****************************************/


//Conn/Socket pointers
var receiverConn *net.UDPConn
var neighbourConnection *net.UDPConn
//varbroadcast_conn *net.UDPConn


//Address containers
var neighbourElevatorAddress string
var localAddress string = ""
var broadcastAddr string = ""

var localIP	string
var localPort int

//Elevator variables

var elev_id int 			// each elevator has an unique ID
var nElevators int	
var prevElevConnected bool = false 	//hmm not 100% sure on usage yet
var singleStateElevator bool = true

//System bools

var circleConnection bool = false		//This bool verifies that we have full circular connection
var watchdoge bool = false				//This bool verifies that the watchdog timer has started.
// timer used in handling I am alive messages
var watchdogTimer *time.Timer

//semaphores
var SemaphoreMessage chan int = make(chan int, 1)
//var destroyReceive bool = false //just to kill one of the read threads
//modify this to be a standard listen function
var SemaphoreRead chan int = make(chan int, 1)
var SemaphoreNewConnection chan int = make(chan int, 1)
var BroadcastMutex chan int = make(chan int, 1)

var aliveMessageContainer chan int = make(chan int,1)


/******************************************Public Functions*****************************************/


/**********************threads**********************/

func Task_networkSupervisor() {
	
	
	for{if (circleConnection == true){break} }
	 //do not send alive messages before theres a full circle connection
	fmt.Println("We broke loop")
	//watchdogTimer = time.AfterFunc(time.Duration(3)*time.Second, watchdogCallback)
	//watchdoge = true
	for {
		if(circleConnection == true && neighbourConnection != nil){
			IAmAliveMsg := Message{MsgID : IAmAlive, Elev_id : GetID()}
			send_msg(IAmAliveMsg)
		}
		select {
		case alive_id :=<-aliveMessageContainer:
			fmt.Printf("Received alive from elev id : %d \n", alive_id)
			break

		case <-time.After(5 * time.Second):
			watchdogCallback() // we done guufed
			return
		}
	time.Sleep(time.Second * 1)
	}

}


		/*
		if (circleConnection == true && watchdoge == true) { //do not send alive messages before theres a full circle connection
			fmt.Println("Sending alive")
			send_msg(IAmAliveMsg)
			time.Sleep(time.Millisecond * 40)
			//if(!circleConnection){
			//	msg := Message{MsgID : CheckNeighbourConnection, LocalAddr : GetLocalAddress()}
			//	send_msg(msg)
			//}			
		}else{
			//if we establish conenction between two elevators, and then add another one, the timer must stop temporarily.
			watchdogTimer.Stop()
			watchdoge = false
		}*/
	

/*
for {
select {
	case alivemsg<-aliveMessageReceived:
		break

	case <-After(10 * Millisecond):
		return
}*/

func Task_sendElevMessages( C_message chan Message) {
	
	//msg := "Hello!"
	//addr := GetLocalAddress()
	//	broadcastLastLocalAddress := "" //This variable will be used to make sure we do not process the same info multiple times.

	var msgRecv Message
	var newMsg bool

	
	for {
		time.Sleep(time.Millisecond * 40)
			if(singleStateElevator){
				//broadcast data every once in a while untill it gets connected?

			}			
			newMsg = false
			select{
				case msgRecv = <-C_message:
					newMsg = true
				default:
					newMsg = false
					//no current commando

			}
			if(newMsg){
				switch(msgRecv.MsgID){
					/*case NewElevatorConnectionEstablished: // maybe use msg ID ?
						//msg := Message{MsgID : NewElevatorConnection, StringMsg : msg, LocalAddr : addr, NumElev : numElev, NewElevAddr : newElevatorAddress} // send newelevator through channel
						//newElevAddr is the address of the new elevator
						if(msgRecv.Elev_id == nElevators){
							fmt.Println("Successfull connection")
						} */
			
					case NewElevatorConnection:
						fmt.Println("received new elevator connection")
						// we simply just send the message to our neighbour since our elev_id does not match the target, which is checked in receiveElevMessages
					//	if(broadcastLastLocalAddress != msgRecv.LocalAddr){
					//		broadcastLastLocalAddress = msgRecv.LocalAddr
							fmt.Println("inside NewElevatorConnection")
							send_msg(msgRecv)	
					//	}

						/*newNeighbourLocaladdr := msgRecv.LocalAddr//  Localaddress of the new elevator
						if(broadcastLastAddress != newNeighbourLocaladdr){
							broadcastLastAddress= msgRecv.LocalAddr
						
							if(elev_id == 1){	//elevator with unique ID 1 will be first to deal with broadcastes messages
								//numElev = message.NumElev
								//if(broadcastLastLocalAddress != newNeighbourLocaladdr ){
									//this address have not been sent before
								broadcastMsg := Message{MsgID : NewElevatorConnection, NewElevatorLocalAddress : newNeighbourLocaladdr, LocalAddr : addr, TargetID : nElevators}//, numElev : nElevators}
								nElevators++
								send_msg(broadcastMsg) 
							}else if(elev_id != nElevators){
								nElevators++ // should maybe have this in broadcastLocaladdress
								send_msg(msgRecv) // forward the message to the correct node
							}else if(elev_id == msgRecv.TargetID){ //can also use message.TargetID == elev_id or TargetID == nElevators
								nElevators++
								//establish connection.
								neighbourConnection.Close() 
								neighbourConnection = nil
								setNeighbourElevatorAddress(msgRecv.NewElevatorLocalAddress) // set neighbour address up
								setNeighbourElevConnection()				//establish new connection por favor
								configMsg := Message{MsgID : NewElevatorInitConfig, NewElevatorLocalAddress : msgRecv.LocalAddr, NumElev : nElevators} 
								send_msg(configMsg)		//sending message to the new node such that it can become a part of the system.
								//neighbourElevatorAddress = msgRecv.newElevatorAddress
								//this is the previous last node on our circle.
							}else if(elev_id == 0) {		//unitialized elevator
								elev_id = msgRecv.NumElev
								nElevators = elev_id
								singleStateElevator = false
								setNeighbourElevatorAddress(msgRecv.NewElevatorLocalAddress)
								setNeighbourElevConnection()
								//send confirmation to elev_id == 1
								msg := Message{MsgID : NewElevatorConnectionEstablished, Elev_id : elev_id }
								send_msg(msg)
							}
						}*/
						case BroadcastMsg:	
						   // <-SemaphoreNewConnection

							
							//msg := Message{MsgID : BroadcastAcknowledged, NewElevatorLocalAddress : msgRecv.LocalAddr, NumElev : nElevators} // I removed "NewElevatorLocalAddress : addr" from the struct since NewElevatorAddress will be responsible for sending the localAddress to the new elevator
							//fmt.Println("Broadcast received, acknowledging it, sending nElevators to the other guy")
							//fmt.Printf("Sending elevators : %i \n", nElevators)
							//send_msg(msg)
							//SemaphoreNewConnection <-1
							fmt.Println("Received a message, performing correct procedure")

						case BroadcastAcknowledged:
							send_msg(msgRecv)

						case NewElevatorConnectionEstablished:
							send_msg(msgRecv) // We just want the first elevator to know that we have successfully connected to it
						//	broadcastLastLocalAddress = ""

							//SemaphoreNewConnection <-1
						case CheckNeighbourConnection:
							send_msg(msgRecv)

						case ExternalOrderComplete:
							send_msg(msgRecv)
						case Debug:
							if(msgRecv.TargetID == GetID()){
							fmt.Println(msgRecv.StringMsg)
							fmt.Println("Yeah nuggah i got u")
							}else{
								send_msg(msgRecv)
							}
						
				}

			}
		}
}
func Task_receiveElevMessages(C_message chan Message, C_elevatorCommand chan int, C_order chan Ch_elevOrder) {
	buffer := make([]byte, 1024)
	time.Sleep(1 * time.Second)
	broadcastLastLocalAddress := ""
	var message Message
	conn := receiverConn
	for {	
			time.Sleep(time.Millisecond * 10)
			n, _ := conn.Read(buffer) //does this even have a deadline? Ask joey. Set deadline on this for 10 sec.	currently not using error.
			if(n != 0){
				_ = json.Unmarshal(buffer[:n], &message)
				switch(message.MsgID){

				case NewElevatorConnection:
					fmt.Println("New elevator connection recv")
					nElevators++
					circleConnection = false
										
						fmt.Println("broadcast != last broadcast")
						
						if(elev_id == 0){
							fmt.Println("WE HAVE REACHED THE LAST ELEVATOR")
							setNeighbourElevatorAddress(message.LocalAddr)
							setNeighbourElevConnection()
							elev_id = message.TargetID
							nElevators = message.NumElev
							fmt.Println("elev_id: ", elev_id)
							fmt.Println("nElevators: ", nElevators)
							singleStateElevator = false
							establishedConnectionMsg := Message{MsgID : NewElevatorConnectionEstablished, Elev_id : elev_id}
							C_message <- establishedConnectionMsg
						}else if(elev_id == (message.TargetID - 1) && nElevators > 2){
							fmt.Println("Connected")
							GetNeighbourElevConnection().Close()
							neighbourConnection = nil
							setNeighbourElevatorAddress(message.NewElevatorLocalAddress)
							setNeighbourElevConnection()
							disableAliveMsg()
							C_message <-message

							//this elevator was the previous last one.

						}else{
							disableAliveMsg()
							C_message <- message
							fmt.Println("Forwarding new elevator connection message.")

						}

					
				case BroadcastMsg:
					//fmt.Println("WE REACH INSIDE BROADCASTMSG")

					if(broadcastLastLocalAddress == message.LocalAddr){break} //all broadcastlastlocal addresses are the same lol.
					circleConnection = false
					broadcastLastLocalAddress = message.LocalAddr
					nElevators++
					if (elev_id == 0) {elev_id = nElevators} // a way to set the elev_id of the first elevator
					
					if (elev_id == 1) { // Only the first elevator should handle incoming broadcasts from newly initialized elevators
						if(nElevators == 1){	//if this is the only elevator standing, then it will connect to the other elevator immediately
							nElevators++
							setNeighbourElevatorAddress(message.LocalAddr)
							setNeighbourElevConnection() 
							fmt.Println("broadcast msg received: ", message.StringMsg, neighbourElevatorAddress) 
							C_message<-message 	// we let the current "message" struct be passed as usual and let the BroadcastMsg in the send task handle it.
							
							newElevatorMessage := Message{MsgID : NewElevatorConnection, LocalAddr: GetLocalAddress(), TargetID : nElevators,
							NewElevatorLocalAddress : message.LocalAddr, NumElev : nElevators} // Create a message that should be sent to the new elevator
							<-SemaphoreNewConnection	//this goes all the way til the messagei s sent
							C_message<-newElevatorMessage
							SemaphoreNewConnection <-1	//this goes all the way til the messagei s sent
							establishedConnectionMsg := Message{MsgID : NewElevatorConnectionEstablished, Elev_id : elev_id}
							C_message<-establishedConnectionMsg

						} else {//try and do nothing when u receive a commit message

							numberOfElev := nElevators
							fmt.Println(" broadcast message receivedw hen nElevators >= 2")
							fmt.Println("Elevators : ")
							fmt.Println(nElevators)
							disableAliveMsg()
							newElevatorMessage := Message{MsgID : NewElevatorConnection, LocalAddr : GetLocalAddress(), TargetID : numberOfElev, NumElev : numberOfElev,
						    NewElevatorLocalAddress : message.LocalAddr, StringMsg : "Broadcast acknowledged"}


							 // Create a message that should be sent to the new elevator

							// we let the current "message" struct be passed as usual and let the BroadcastMsg in the send task handle it.
							//C_message<-message // no point in doing this.
							//<-SemaphoreNewConnection
							C_message<-newElevatorMessage
							//SemaphoreNewConnection <-1
						}
					}

				case BroadcastAcknowledged:
					//someone acknowledged your message.
					if (elev_id == 0) { // Since the elev_id is zero, this is the elevator who sent the broadcast, and should therefor receive the acknoledge
						fmt.Println("Broadcast acknowledged")
						elev_id = message.NumElev // we set the elev_id such that the message from NewElevatorConnection, which is behind, will correctly identify this as the Target_ID	
						<-	SemaphoreNewConnection 						
						nElevators = elev_id
						SemaphoreNewConnection <- 1 //must make sure if we apply a new elev ID before we send the newElevatorconnection message.
					//	singleStateElevator = false
					} else { // pass the acknowledge message to the neighbour and update the the current nElevators variable
						fmt.Println("inside broadcst acknowledge")
						nElevators = message.NumElev

						// need to check if the current elevator is the second to last elevator, such that we can rewire the connection to the new elevator
						if((nElevators - elev_id) == 1) {
							GetNeighbourElevConnection().Close()
							neighbourConnection = nil
							setNeighbourElevatorAddress(message.NewElevatorLocalAddress)
							setNeighbourElevConnection()
						}

						C_message<-message
					}
				

				case NewElevatorConnectionEstablished:
					if(message.Elev_id == GetID()){
						fmt.Println("Full circle connection")
						broadcastLastLocalAddress = ""
						singleStateElevator = false
					    enableAliveMsg()
						//watchdogTimer.Reset(time.Duration(3)*time.Second)
					}else{
						fmt.Println("The newest elevator has successfully connected to us!")
						fmt.Println("forwarding the message")
						C_message <- message
						//circleConnection = true
						broadcastLastLocalAddress = ""
					}
				case IAmAlive:
					if (!circleConnection) {
						enableAliveMsg()
						
					}
					aliveMessageContainer <- message.Elev_id



				case CheckNeighbourConnection:
					if (GetNeighbourElevConnection() == nil) {
						// We have found the elevator with nil neighbourconnection
						fmt.Println("Founc elevator with nil neighbourconnection: ", elev_id)
						setNeighbourElevatorAddress(message.LocalAddr)
						setNeighbourElevConnection()
						// for now i will set the circleConn here since when this is called we have a full circle again
						circleConnection = true
					} else {
						// pass the message to the next neighbour
						C_message<-message
					}		
					
						
				case OrderRequestEvaluation:
					channel_message := Ch_elevOrder{Floor : message.Floor, Button : message.Buttontype, Elev_score : message.Elev_score, Elev_id : message.Elev_id} 
					C_elevatorCommand <- OrderRequestEvaluation
					C_order <- channel_message
					//C_order.Floor <- message.Floor
					//C_order.Buttontype <-message.Buttontype
					//C_order.Elev_score <- message.Elevbroadcast(broadcast_conn)_score //ask studass
				case OrderRequest:
					fmt.Println("Right type")
					if(message.Elev_targetID == elev_id){
						fmt.Println("Right ID")
						C_elevatorCommand <- OrderRequest
						fmt.Println("Sent to elevator channel correctly")
						channel_message := Ch_elevOrder{Floor : message.Floor, Button : message.Buttontype}
						C_order <- channel_message
						fmt.Println("Sent struct further to C_order correctly")
						//C_order.Floor <- message.Floor
						//C_order.Buttontype <- message.Buttontype //once again ask studass
					}else{
						Send_newOrder(message.Floor, message.Buttontype, message.Elev_targetID)//,nil)	
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
				case ExternalOrderComplete:
					if(message.TargetID != GetID()){
						C_message <-message
						C_elevatorCommand <- ExternalOrderComplete
						channel_message := Ch_elevOrder{Floor : message.Floor, Button : message.Buttontype}
						C_order <-channel_message
					}
				case Debug:

					if(message.TargetID == elev_id){
					fmt.Println(message.StringMsg)	
					C_message <- message
					}else{
						C_message<-message
					}
					
				}
			}
		}
	for{}
}

func Task_broadcastSupervisor(){
	//time.Sleep(60 * time.Second)	//wait 30 sec after init before you run this.
	//broadcast_msg := Message{MsgID : BroadcastMsg,StringMsg : msg, LocalAddr : addr}
	//broadcast_Tim
	for{

		if(singleStateElevator){//if it is not connected to any elevator, it will try broadcasting for 3 sec every 3 min.
			broadcast_conn := netw.GetConnectionForDialing(broadcastAddr)
			broadcast(broadcast_conn)
			fmt.Println("We broadcastet yup")
			fmt.Println(singleStateElevator)
		}else{
			fmt.Println("skipped broadcasting due to connection")
		}

		time.Sleep(120 * time.Second)
	}
}


func broadcast(conn *net.UDPConn) {
	//t0 := time.Now()
	//run_bc := true
	fmt.Println("Broadcast initialising")
	msg := "Hello I am broadcasting!"
	addr := GetLocalAddress()
	broadcast_msg := Message{MsgID : BroadcastMsg,StringMsg : msg, LocalAddr : addr}
	//broadcast_endMsg := Message{MsgID : BroadcastMsg, StringMsg : "Broadcasting done"}
	broadcastTimer := time.NewTimer(time.Duration(3)*time.Second)	//possibly make a text file where u can store values? //stfu seriosuly stfu
LOOP:
	for{
		//if(!run_bc){break}
		select{
			case <-broadcastTimer.C:
				fmt.Println("this gets called")
				//buffer,_ := json.Marshal(broadcast_endMsg)
				//_,_ = conn.Write(buffer)
				break LOOP
			default:
			fmt.Println("STILL BROADCASTING")
			if(!singleStateElevator){
				break LOOP
			}

			buffer,err := json.Marshal(broadcast_msg)
			
			if err != nil {
				fmt.Println("ERROR IN MARSHAL")
				fmt.Println("%s", err)
			}
			_,err = conn.Write(buffer)
			if err != nil {
				fmt.Println("WE GOT ERROR!")
			}
			time.Sleep(time.Millisecond * 50)

			//fmt.Println("Message sent")
		}
	}
	fmt.Println("broadcast socket closed")
	conn.Close()
}



/**********************Send functions**********************/

func Send_requestedOrderEvaluation(elev_score []float64, floor int, buttontype ButtonType,elev_id_original int ){
	msg := Message{MsgID : OrderRequestEvaluation,Elev_id : elev_id_original, Elev_score : elev_score,Floor : floor, Buttontype : buttontype}
	send_msg(msg)
	
}

func Send_newOrder(floor int, button ButtonType, chosenElevator int){//, conn *net.UDPConn){	
	msg := Message{MsgID : OrderRequest,Elev_targetID: chosenElevator, Floor : floor, Buttontype : button }
	send_msg(msg)

}

func Send_elevInitCompleted(successfull bool){
	msg := Message{MsgID : Elevator_initializationStatus,Elev_targetID : elev_id, SuccessfullInit : successfull}
	send_msg(msg)
	
}

func send_broadcastData(){}
//Private send function, which ultimately sends the message.
func send_msg(msg Message){
	buffer,err := json.Marshal(msg)
	if err != nil{
		fmt.Println("ERROR IN MARSHAL OF SENDING message")
		fmt.Println("%s", err)
	}
	if(neighbourConnection == nil){
		return
	}
	<-SemaphoreMessage
	_,err = neighbourConnection.Write(buffer)
	SemaphoreMessage <- 1
/*	if err != nil {
		fmt.Println("Error in writing to neighbourconnection")
		neighbourConnection = nil
		neighbourElevatorAddress = ""
		circleConnection = false
		if(nElevators == 2) { // we need to treat this case by itself, ask walid whether this is an alrgiht place to put it
			nElevators--
			singleStateElevator = true
		}
	}*/

}

func Send_debug(text string, targetId int){
	msg := Message{MsgID : Debug, StringMsg : text,TargetID : targetId}
	send_msg(msg)
}

func Send_externalCommandComplete(floor,button int, targetID int){
	externalOrderCompleteMsg := Message{MsgID : ExternalOrderComplete, Floor : floor, Buttontype : ButtonType(button), TargetID : targetID}
	send_msg(externalOrderCompleteMsg)
}

/**********************Get / Set functions**********************/

func setNeighbourElevatorAddress(address string) {
	if neighbourElevatorAddress != address {
		neighbourElevatorAddress = address
	}
}

func setNeighbourElevConnection() {
	if neighbourConnection == nil && neighbourElevatorAddress != "" {
		neighbourConnection = netw.GetConnectionForDialing(neighbourElevatorAddress)
	}
}

func setReceiverConn(localaddress string){
	receiverConn = netw.GetConnectionForListening(localaddress)
}


func GetNeighbourElevAddress()(string) {
	return neighbourElevatorAddress
}
func GetNeighbourElevConnection()(*net.UDPConn) {
	return neighbourConnection
}

func GetLocalAddress() (string) {
	return localIP+":"+ strconv.Itoa(20000+localPort)
}

func GetID()(int){
	return elev_id
}

func GetNelevators()(int){
	return nElevators
}

func GetSingleElevatorState()(bool){
	return singleStateElevator
}

/**********************Initialize function**********************/

func InitElevatorNetwork(){
	nElevators = 0
	elev_id = 0
	init_localAddress()
	broadcastPortInput := NumberUserInput("Broadcast port:")

	broadcastAddr = broadcastIP + ":" + strconv.Itoa(20000 + broadcastPortInput) //broadcastIP should not be a constant string. It should be dependent on which local area network we are in. make this possible joey.
	//broadcast_conn := netw.GetConnectionForDialing(broadcastAddr)
	//broadcast(broadcast_conn)

	receiverAddress := 	"" + ":" + strconv.Itoa(20000+localPort)
	setReceiverConn(receiverAddress) // Walid: I put broadcast() above setReceiverConn() because we were receiving messages from our own broadcast after it expires
	//broadcast(broadcast_conn)
}

func init_localAddress() {
	localIP = netw.GetLocalIP()
	localPort = netw.GetPort()
}




/**********************Elevator identifier functions**********************/


func AssignElevatorVariables() { 
	if nElevators == 0 && neighbourElevatorAddress == "" {
		nElevators = 1;
	}
	if elev_id == 0 {
		elev_id = nElevators;
	}else if(elev_id > 0){
		return
	}
}

func newElevDetected(addr string) {
	fmt.Println("ADDRESS: ", addr, neighbourElevatorAddress)
	/*if elev_id == 0 {
		elev_id = nElevators
	} else if elev_id == 1 && newElevatorAddress != addr{
		nElevators++
		newElevatorAddress = addr
	}*/
	//if(elev_id)
}


/**********************Helping functions**********************/

func isNeighbourElevatorAddressEmtpy()(bool) {
	if (neighbourElevatorAddress != "") {
		return false
	}
	return true
}


func NumberUserInput(typeOfInput string) (int) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(typeOfInput)
	inputString, _ := reader.ReadString('\n')
	inputString = strings.Replace(inputString, "\n", "", -1)
	input,err := strconv.Atoi(inputString)

	if err != nil {
		fmt.Println("ERROR IN CONVERSION")
		fmt.Printf("%s \n", err)
		return 0
	}
	return input
}

func enableAliveMsg(){
	circleConnection = true
}

func disableAliveMsg(){
	circleConnection = false
}

func watchdogCallback() {
	//Nothing here atm.
	fmt.Println("Watchdog called")
	msg := Message{MsgID : ConnectionLost, LocalAddr : GetLocalAddress(),Elev_id :  GetID()}
	send_msg(msg)
}

