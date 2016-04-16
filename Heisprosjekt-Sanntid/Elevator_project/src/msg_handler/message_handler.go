package msg_handler


import(
	"fmt"
	"net"
	"../netw"
	"encoding/json"	
	"strconv"
	"time"
	"bufio"
	"os"
	"strings"
)	



const broadcastIP string = "129.241.187.255"
const broadcastPort string = "20022"

var LocalIP	string
var LocalPort int
var newElevatorAddress string
var neighbourElevatorAddress string
var neighbourConnection *net.UDPConn
var localAddress string = ""
var listenConnection *net.UDPConn

var prevElevConnected bool = false 	//hmm not 100% sure on usage yet
var singleStateElevator bool = true

var elev_id int 			// each elevator has an unique ID
var nElevators int	
var IPaddress string // along with port

var SemaphoreMessage chan int = make(chan int, 1)

var receiverConn *net.UDPConn

func AssignElevatorVariables() { //lol calling this singleElevator, and then having it assign a new elev_id if there are multiple elevators. smh.
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
	fmt.Println("ADDRESS: ", addr, newElevatorAddress)
	/*if elev_id == 0 {
		elev_id = nElevators
	} else if elev_id == 1 && newElevatorAddress != addr{
		nElevators++
		newElevatorAddress = addr
	}*/
	//if(elev_id)
}

func send_msg(msg Message){
	buffer,err := json.Marshal(msg)
	if err != nil{
		fmt.Println("ERROR IN MARSHAL OF SENDING message")
		fmt.Println("%s", err)
	}
	<-SemaphoreMessage
	_,_ = neighbourConnection.Write(buffer)
	SemaphoreMessage <- 1
	fmt.Println("semaphore done")
}


func SendElevMessages( C_message chan Message) {
	
	//msg := "Hello!"
	addr := GetLocalAddress()
	broadcastLastAddress := "" //This variable will be used to make sure we do not process the same info multiple times.
	var msgRecv Message
	var newMsg bool
	
	for {
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
					case NewElevatorConnectionEstablished: // maybe use msg ID ?
						//msg := Message{MsgID : NewElevatorConnection, StringMsg : msg, LocalAddr : addr, NumElev : numElev, NewElevAddr : newElevatorAddress} // send newelevator through channel
						//newElevAddr is the address of the new elevator
						if(msgRecv.Elev_id == nElevators){
							fmt.Println("Successfull connection")
						}
			
					case NewElevatorConnection:
						newNeighbourLocaladdr := msgRecv.LocalAddr//  Localaddress of the new elevator
						if(broadcastLastAddress != newNeighbourLocaladdr){
							broadcastLastAddress= msgRecv.LocalAddr
							/*if(elev_id == 0){
								elev_id == 1
							}/*/
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
								SetNeighbourElevatorAddress(msgRecv.NewElevatorLocalAddress) // set neighbour address up
								setNeighbourElevConnection()				//establish new connection por favor
								configMsg := Message{MsgID : NewElevatorInitConfig, NewElevatorLocalAddress : msgRecv.LocalAddr, NumElev : nElevators} 
								send_msg(configMsg)		//sending message to the new node such that it can become a part of the system.
								//neighbourElevatorAddress = msgRecv.newElevatorAddress
								//this is the previous last node on our circle.
							}else if(elev_id == 0) {		//unitialized elevator
								elev_id = msgRecv.NumElev
								nElevators = elev_id
								singleStateElevator = false
								SetNeighbourElevatorAddress(msgRecv.NewElevatorLocalAddress)
								setNeighbourElevConnection()
								//send confirmation to elev_id == 1
								msg := Message{MsgID : NewElevatorConnectionEstablished, Elev_id : elev_id }
								send_msg(msg)
							}
						}
					case BroadcastMsg:	
						msg := Message{MsgID : BroadcastAcknowledged, NewElevatorLocalAddress : addr, NumElev : nElevators}
						fmt.Println("Broadcast received, acknowledging it, sending nElevators to the other guy")
						fmt.Printf("Sending elevators : %i \n", nElevators)
						elev_id = 1
						send_msg(msg)
				}

			}
			/*
			numElev := nElevators
			newElevAddr := newElevatorAddress
			dial_msg := Message{MsgID : NewElevatorConnection,StringMsg : msg, LocalAddr : addr, NumElev : numElev, NewElevAddr : newElevAddr}
			
			buffer,err := json.Marshal(dial_msg)
	
			if err != nil {
				fmt.Println("ERROR IN MARSHAL")
				fmt.Println("%s", err)
			}
			_,err = neighbourConnection.Write(buffer)

			if err != nil {
				fmt.Println("ERROR IN DIALING")
				fmt.Println(err)		
			}/*/
		}
}
var destroyReceive bool = false //just to kill one of the read threads
//modify this to be a standard listen function
var SemaphoreRead chan int = make(chan int, 1)

func ReceiveElevMessages(C_message chan Message) {
	buffer := make([]byte, 1024)

	//addr := "" currently not using this.
	//numElev := 0	currentlyn not using this.
	broadcastLastLocalAddress := ""
	//t0 := time.Now()
	//time.Sleep(time.Second * 10)
	var message Message
	//deadlineTime := time.Second * 10
	//deadtest := time.Time
	//deadtest.Duration = deadlinetime
	//deadlineTimer := time.Now().Add(5 * time.Second)
	conn := receiverConn
	//conn.SetReadDeadline(deadlineTimer) // CURRENTLY TIMER IS SET TO 10 SEC
	//timerbla := time.NewTimer(time.Second * 5)
	for {	if(destroyReceive){break}
	//		timer_beforeRead := time.Now()
		//	fmt.Println("waiting for message")
			//fmt.Println("localaddress: ",conn.RemoteAddr().String())
			//<- SemaphoreRead	// need deadline here.
			n, _ := conn.Read(buffer) //does this even have a deadline? Ask joey. Set deadline on this for 10 sec.	currently not using error.
			//SemaphoreRead <- 1
	//		timer_afterRead := time.Now()
	/*		if(timer_afterRead.Sub(timer_beforeRead) >(2*time.Second) && singleStateElevator){// && !prevElevConnected ){	//2 sec after listening it will be set to single elevator by default
				AssignElevatorVariables()
				if(nElevators == 1){
					singleStateElevator = true
				}*/
		//	}
		//	fmt.Println("Received something")
			//Theres...actually no listen commands? I receive data..might implement some sort of command system in order to listen for some specifics..
			if(n != 0){
				_ = json.Unmarshal(buffer[:n], &message)
				switch(message.MsgID){
				case NewElevatorConnection:
					if(broadcastLastLocalAddress != message.NewElevatorLocalAddress){
						broadcastLastLocalAddress = message.NewElevatorLocalAddress
						C_message <- message
					}
					fmt.Println("NewElevatorConnection")
					//newNeighbourLocaladdr := message.LocalAddr //local address of the new elevator // not using this yet.
					
				case BroadcastMsg:
					//established connection with someone else. 
					//msg := Message{MsgID : BroadcastAcknowledged}
					//since this was the elevator acknowleding a connection, it will send acknowledgement of connection.
					//AssignElevatorVariables()
					if(broadcastLastLocalAddress == message.LocalAddr){break}
					broadcastLastLocalAddress = message.LocalAddr
					nElevators++
					if(nElevators == 1){	//if this is the only elevator standing, then it will connect to the other elevator immediately
						nElevators++
						SetNeighbourElevatorAddress(message.LocalAddr)
						setNeighbourElevConnection() 	//connection established.
						fmt.Println("broadcast msg received, nElevators = 1, ", message.StringMsg, neighbourElevatorAddress)

						C_message<-message
					}else{//commented out to test a special situation.
						fmt.Println("something else")
						//message.MsgID = NewElevatorConnection
						//C_message <-message
						//fmt.Println("Broadcast msg received nelevators != 1")
					}

				
				case BroadcastAcknowledged:
					//someone acknowledged your message. Now you wait for proper connection
					SetNeighbourElevatorAddress(message.NewElevatorLocalAddress )
					setNeighbourElevConnection()
					fmt.Println(neighbourElevatorAddress)
					nElevators = message.NumElev
					elev_id = nElevators
					fmt.Println("Broadcast acknowledged")
				}
		}	
//		for{}
	}
	for{}
}
/*
func setupElevatorListen(){
	listenIP := ""
	listenAddr := listenIP + ":" + strconv.Itoa(20000 + LocalPort)
	listen_conn := netw.GetConnectionForListening(listenAddr)
	listenConnection = listen_conn

	//listenW(listen_conn) commented out in order to avoid error
}
/*func ListenForElevMessages() { AVOIDING THIS ATM SUCH THATI  CAN AVOID ERROR MESSAGE.
	listenIP := ""
	listenAddr := listenIP + ":" + strconv.Itoa(20000+LocalPort) 
	listen_conn := netw.GetConnectionForListening(listenAddr)
	
//	listen(listen_conn)	commented out to avoid errors
}*/



func broadcast(conn *net.UDPConn) {
	//t0 := time.Now()
	//run_bc := true
	msg := "Hello I am broadcasting!"
	addr := GetLocalAddress()
	broadcast_msg := Message{MsgID : BroadcastMsg,StringMsg : msg, LocalAddr : addr}
	broadcastTimer := time.NewTimer(time.Second * 3)	//possibly make a text file where u can store values?
	LOOP:
	for{
		//if(!run_bc){break}
		select{
			case <-broadcastTimer.C:
				break LOOP
			default:
			buffer,err := json.Marshal(broadcast_msg)
	
			if err != nil {
				fmt.Println("ERROR IN MARSHAL")
				fmt.Println("%s", err)
			}
			_,_ = conn.Write(buffer)
			
			//fmt.Println("Message sent")
		}
	}
	fmt.Println("broadcast socket closed")
	conn.Close()
}
	/*broadcast_msg := Message{MsgID : BroadcastMsg,StringMsg : msg, LocalAddr : addr}

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
     }*/


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
	send_msg(msg)
	/*buffer,err := json.Marshal(msg)
	if err != nil{
		fmt.Println("ERROR IN MARSHAL OF REQUESTED ORDER")
		fmt.Println("%s", err)
	}*/
	
	//_,_ = neighbourconnection.Write(buffer)
}

func Send_newOrder(floor int, button ButtonType, chosenElevator int){//, conn *net.UDPConn){	
	msg := Message{MsgID : OrderRequest,Elev_targetID: chosenElevator, Floor : floor, Buttontype : button }
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

func setNeighbourElevConnection() {
	if neighbourConnection == nil && neighbourElevatorAddress != "" {
		neighbourConnection = netw.GetConnectionForDialing(neighbourElevatorAddress)
	}
//	fmt.Println("bla: ", neighbourConnection.LocalAddr().String())
}

func send_broadcastData(){}

func GetNeighbourElevAddress()(string) {
	return neighbourElevatorAddress
}



func TestSetNeighBourElevConnection(conn *net.UDPConn){
	neighbourConnection = conn
}


func GetNeighbourElevConnection()(*net.UDPConn) {
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
	//for{if(GetNeighbourElevConnection)}
	//SemaphoreMessage <- 1 
//	conn := GetNeighbourElevConnection()
	
	buffer := make([]byte, 1024)
	var message Message
	elev_id = 1 //testing purposes..remove this
	//listenIP := ""
	//listenAddr := listenIP + ":" + strconv.Itoa(20000 + LocalPort)
	//listen_conn := netw.GetConnectionForListening(listenAddr)
	conn := receiverConn
	//<-SemaphoreRead // take the semaphore nadn ever give it back
	destroyReceive = true
	for{
		fmt.Println("reading elevator messages now")
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
				//C_order.Elev_score <- message.Elevbroadcast(broadcast_conn)_score //ask studass
			case OrderRequest:
				fmt.Println("Right type")
				if(message.Elev_targetID == elev_id){
					fmt.Println("Right ID")
					C_elevatorCommand <- 2
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
			}
		}
	}
}


func InitElevatorNetwork(){
	nElevators = 0
	elev_id = 0
	init_localAddress()
	broadcastPortInput := numberUserInput("Broadcast port:")

	broadcastAddr := broadcastIP + ":" + strconv.Itoa(20000 + broadcastPortInput) //broadcastIP should not be a constant string. It should be dependent on which local area network we are in. make this possible joey.
	broadcast_conn := netw.GetConnectionForDialing(broadcastAddr)
	broadcast(broadcast_conn)

	receiverAddress := 	"" + ":" + strconv.Itoa(20000+LocalPort)
	setReceiverConn(receiverAddress) // Walid: I put broadcast() above setReceiverConn() because we were receiving messages from our own broadcast after it expires
}

func setReceiverConn(localaddress string){
	receiverConn = netw.GetConnectionForListening(localaddress)
}


func init_localAddress() {
	LocalIP = netw.GetLocalIP()
	LocalPort = netw.GetPort()
}

func numberUserInput(typeOfInput string) (int) {
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
