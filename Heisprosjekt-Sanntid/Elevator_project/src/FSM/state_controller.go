package FSM

import(
	"fmt"
	"../elev_driver"
	//"../queue"
	"time"
	//"../netw"
	"../msg_handler"
	//"encoding/json"
	//"bufio"

	//"os"
	
	//"../netw"	//using it temporarily to test
)

/***************************************************************************************************************************\\\\\\\\
This package contains the elevator system.

In this file, the state machine is written. 

Stop button is not taking into account.

*****************************************************************************************************************************//////

/******************************************Variables*****************************************/

type state_slaveElevator int
const (
	INITIALIZATION state_slaveElevator = iota
	IDLE  
	MOVING 
	DOOR_OPEN 
	DOOR_CLOSED
	NOT_INITIALIZATED
)


func InitElevator(){
	nFloors := setNFloors()
	event_init(nFloors)
	InitDone = true	

}

var state_slave state_slaveElevator
var InitDone bool = false				//used to acknowledge initialization
var prevfloor,prevbuttontype int = -1,-1	//Using this to avoid avoid spam after holding a button. 
var prevExternalFloor, prevExternalButton int = -1,-1	//same with external buttons. Not sending same request multiple times.

func Thread_elevatorStateMachine(C_elevatorCommand chan int,C_order chan msg_handler.Ch_elevOrder){		//Make a channel int in main
	//miliseconds := 40
//	timer_miliseconds := time.Duration(miliseconds)*time.Millisecond
	msgType := -1
	//notSingleElevator := true
	n_floors := GetNFloors() 
	fmt.Println(orders)
	var notSingleElevator bool = false
	//var floorb, buttontypeb int 
	for{
		time.Sleep(time.Millisecond * 100)
		/*if(queue.Orders[floorb][buttontypeb] != 1 ){
			floorb,buttontypeb =setOrder()	
			
			queue.Orders[floorb][buttontypeb] = 1
			fmt.Println(queue.Orders)
		}*/
		
		
		//time.Sleep(timer_miliseconds)
		//running elevator
		/*if(msg_handler.GetNelevators() > 1){
			notSingleElevator = true
		}else{
			notSingleElevator = false
		}*/
		if(msg_handler.GetSingleElevatorState()){
			notSingleElevator = false

		}else{
			notSingleElevator = true
		}

	//fmt.Println("we get here ")
	//	msg_handler.GetNelevators not used
		for floor := 0; floor< n_floors; floor++{
			for buttontype:=elev_driver.BUTTON_CALL_UP;buttontype <= elev_driver.BUTTON_COMMAND; buttontype++{
				if(floor == 0 && buttontype == elev_driver.BUTTON_CALL_DOWN){
					continue
				}else if(floor == n_floors - 1 && buttontype == elev_driver.BUTTON_CALL_UP){
					continue	
				}
				if((buttontype != elev_driver.BUTTON_COMMAND && elev_driver.Elev_get_button_signal(buttontype,floor) == 1) && notSingleElevator && (prevExternalButton == -1 && prevExternalFloor == -1)){// && notSingleElevator ){			//outside button pressed					
					fmt.Println("This is where u fucked up walid")
					fmt.Printf(" %i		%i 	\n", floor,buttontype)
					prevExternalFloor = floor
					prevExternalButton = buttontype
					elev_driver.Elev_set_button_lamp(buttontype,floor,1)
					event_outsideButtonPressed(floor,buttontype)
				}else if(buttontype == elev_driver.BUTTON_COMMAND && elev_driver.Elev_get_button_signal(buttontype,floor) == 1 && notSingleElevator){	
						event_newQueueRequest(floor,buttontype)	
						fmt.Println("this fuck face got activated")
				}else if(elev_driver.Elev_get_button_signal(buttontype,floor) == 1 && orders[floor][buttontype] == 0 && !notSingleElevator){							//inside button pressed
					event_newQueueRequest(floor,buttontype)					
						fmt.Println("This fuck face got acctivated maybe. maybe its called fuck you?")
				}
				if( (orders[floor][buttontype] == 1)&& (prevfloor != floor) && (prevbuttontype != buttontype) ) {
					prevfloor,prevbuttontype = floor,buttontype //to avoid excessive button mashing order...
				//	fmt.Printf("This is called even though we are not in an idle state yet. Floor : %i, buttontype : %i", floor,buttontype)
					event_queueNotEmpty()
				}
			}	
		}
		currentfloor := currentFloor()
		

		//elev_driver.Elev_set_floor_indicator(2)
		
		if(currentfloor >= 0 && currentfloor<4){
		
			elev_driver.Elev_set_floor_indicator(currentfloor)
			last_direction := GetLastDirection()
			SetLastFloor(currentfloor)
			
			for button:= elev_driver.BUTTON_CALL_UP; button <=elev_driver.BUTTON_COMMAND;button++{
				switch(button){
					case elev_driver.BUTTON_CALL_UP: 
						if( (lastDirection == 1) && (orders[currentfloor][button] == 1)){
					

							event_floorInQueue(currentfloor,button)
						}else if(last_direction == -1 && orders[currentfloor][button] == 1 && (checkQueueList(ALL_ORDERS,0,1) == 1) || 	(orders[0][button] == 1 && currentfloor == 0) ){
							event_floorInQueue(currentfloor,button)
						}
					case elev_driver.BUTTON_CALL_DOWN:
						if(lastDirection == -1 && orders[currentfloor][button] == 1){
							event_floorInQueue(currentfloor,1)
						}else if( (lastDirection == 1) && (orders[currentfloor][button] == 1) && (checkQueueList(ALL_ORDERS,0,1) == 1) || ( (orders[n_floors - 1][button] == 1) && currentfloor == n_floors-1 )) {
							
							event_floorInQueue(currentfloor,button)
						}
					case elev_driver.BUTTON_COMMAND:
						if(orders[currentfloor][button] == 1){
							event_floorInQueue(currentfloor,button)
						}
						

				}
			/*	if(queue.Orders[currentfloor][i] == 1){
					Event_floorInQueue(currentfloor)
				//	fmt.Println(currentfloor)
					//fmt.Println(queue.Orders)
					//fmt.Printf("Queue at: %i", currentfloor)
				}*/
			}
		}
		if(!notSingleElevator){
			continue
			//skip the last part because it isonly for multiple elevators.
		}
		select{
			case msgrecv := <-C_elevatorCommand:
				msgType = msgrecv
			default:
				msgType = -1 
		}
		if(msgType != -1){
			switch(msgType){
			case msg_handler.OrderRequestEvaluation:
				order := <-C_order
				event_evaluateRequest(order.Floor, int(order.Button),order.Elev_id, order.Elev_score)
				//Recv Order request evaluation
			case msg_handler.OrderRequest:
				fmt.Println("Have received order")
				order := <-C_order
				if(orders[order.Floor][order.Button] == 1){
						
				}else{
			 		event_newQueueRequest(order.Floor,int(order.Button))
				}
				//Recv Order request
			case msg_handler.ExternalOrderComplete:
				order_info := <-C_order
				elev_driver.Elev_set_button_lamp(int(order_info.Button),order_info.Floor,0)
			case 10:
				//Recv elevator status from all elevators
			}
		}
	}
}


func currentFloor()(int){
	for i:= 0; i<4;i++{
		if(elev_driver.Elev_get_floor_sensor_signal() == i){
			return i
		}
	}
	return -1
}

													




func event_init(nTotalFloors int){
	init_success := elev_driver.Elev_init() // this should be tested during phase 1 initiliazation possibly. Also send a value indicating it was not properly initiated?
	if(init_success == 0){
		fmt.Println("Initialization failed")
		state_slave = NOT_INITIALIZATED
		//msg_handler.Send_elevInitCompleted(false)
		return
	}
	queue_init(nTotalFloors) // 
	
	
	if(elev_driver.Elev_get_floor_sensor_signal() != 0){
		elev_driver.Elev_set_motor_direction(-1)
		for{
		        if(elev_driver.Elev_get_floor_sensor_signal() == 0){
		     		   break
		        }
		}
		elev_driver.Elev_set_motor_direction(0) // c
	}
	
	state_slave = IDLE
	
	//msg_handler.Send_elevInitCompleted(true) crashes the program.
	
	fmt.Printf("Init done with %i floors \n", nTotalFloors)

}

func event_queueNotEmpty(){

        switch state_slave{
               case IDLE:

               	    fmt.Println("Queue not empty")
               	    //queue.MainFloor = 3
               	
                    setNextMainFloor()
               	    //queue.MainFloor = 2
                    fmt.Printf("Main floor: %i", GetMainFloor())
                    elev_driver.Elev_set_motor_direction((GetLastDirection()))
                    //elev_driver.Elev_set_motor_direction(1)
                    if(mainFloor == GetLastFloor()){
						state_slave = IDLE					
					}else{
						state_slave = MOVING
					}
				default:
					//ignore the allocation of queue if the order is not at idle
					prevfloor,prevbuttontype = -1,-1
                    
        }
}

func event_floorInQueue(floor int, button int){ // this is activated if one of the orders are on the floor.
	switch(state_slave){
	case MOVING:
		fmt.Println("Floor in queue called, moving state.")
		/*if(queue.Orders[floor][2] == 1){
			fmt.Println("Sup boy")
			elev_driver.Elev_set_motor_direction(0)
			queue.Orders[floor][2] = 0
			fmt.Println(queue.Orders)
			elev_driver.Elev_set_button_lamp(elev_driver.BUTTON_COMMAND,floor,0)
		}
		if( (queue.Last_direction == 1) && ( (queue.Orders[floor][elev_driver.BUTTON_CALL_UP] == 1) || queue.Orders[n_floors - 1][elev_driver.BUTTON_CALL_DOWN]  == 1)){
				elev_driver.Elev_set_motor_direction(0)
				queue.Orders[floor][0] = 0
				elev_driver.Elev_set_button_lamp(elev_driver.BUTTON_CALL_UP,floor,0)
		} else if( (queue.Last_direction == -1) && ( (queue.Orders[floor][elev_driver.BUTTON_CALL_DOWN] == 1) || queue.Orders[0][elev_driver.BUTTON_CALL_UP] == 1)){
				elev_driver.Elev_set_motor_direction(0)
				queue.Orders[floor][1] = 0
				elev_driver.Elev_set_button_lamp(elev_driver.BUTTON_CALL_DOWN,floor,0)
		}*/
		elev_driver.Elev_set_motor_direction(0)
		elev_driver.Elev_set_floor_indicator(floor+1)

		
		//timer := time.NewTimer(time.Second * 3 )
		
		//<- timer.C use this to indicate timer done
	case IDLE:
		fmt.Println("Floor in queue called, IDLE state") //this means an order was given while at the same floor.

	}
		if(button != BUTTON_COMMAND && !msg_handler.GetSingleElevatorState()){
			msg_handler.Send_externalCommandComplete(floor, button, msg_handler.GetID())
		}
		prevfloor = -1
		prevbuttontype = -1
		elev_driver.Elev_set_button_lamp(button, floor, 0)
		orders[floor][button] = 0
		elev_driver.Elev_set_door_open_lamp(1)
		state_slave = DOOR_OPEN
		fmt.Println(orders)
		//timer := time.NewTimer(time.Second * 3 )
		seconds := 3
		timer_seconds := time.Duration(seconds)*time.Second
		time.AfterFunc(timer_seconds,event_doorTimeout)

		

}

func event_doorTimeout(){
	fmt.Println("Timeoutcalled")
	switch(state_slave){
	case DOOR_OPEN:
		elev_driver.Elev_set_door_open_lamp(0)
		if(GetMainFloor() == GetLastFloor()){
			state_slave = IDLE
			fmt.Println("Main floor reached")
		}else{
			elev_driver.Elev_set_motor_direction(GetLastDirection())
			state_slave = MOVING
			fmt.Println("Still in moving state")
		}
	}

	elev_driver.Elev_set_door_open_lamp(0)

}
func event_newQueueRequest(floor int, button int){
// add order regardless of current state
	switch(state_slave){
	case NOT_INITIALIZATED:
			return
	}
	fmt.Printf("queueRequest in floor: %d \n", floor+1)
	orders[floor][button] = 1
	fmt.Println(orders)
	elev_driver.Elev_set_button_lamp(button, floor,1)
}

func event_outsideButtonPressed(floor int, button int){
	fmt.Println("Outside button pressed")
	switch(state_slave){
	case NOT_INITIALIZATED:
			return
	}
	//Accepting a button outside pressed does not matter in which state we are. All orders will be accessed.
	score := calculateOrderScore(floor, button)
	score_array := []float64{score}
	
	
	msg_handler.Send_requestedOrderEvaluation(score_array, floor,msg_handler.ButtonType(button), msg_handler.GetID() )			
}


// ADD LIGHTS TO THE BUTTONS
func event_evaluateRequest(floor int,  button int, elev_id int, elev_score []float64){
	fmt.Printf("Asked to calculate from elevator ID : %i \n",elev_id)
	elev_driver.Elev_set_button_lamp(button,floor,1)

	switch(state_slave){
	case NOT_INITIALIZATED:
			return
	}
	if(msg_handler.GetID()== elev_id){
		fmt.Println("received the calculating order back, lets see")
		highestElevatorScore := 0.0
		winningElevator := 0
		prevExternalButton = -1
		prevExternalFloor = -1 
		for i := 0 ; i<msg_handler.GetNelevators() ;i++{
			fmt.Printf("Elev score[%i] = %i \n",i, elev_score[i])
			if(elev_score[i] >= highestElevatorScore){
				highestElevatorScore = elev_score[i]
				winningElevator = i + 1
			}
		}
		
		if(elev_id == winningElevator ){
			//this elevator is the winner
			fmt.Println("this elevator is the winner")
			fmt.Print(" ")
			fmt.Println(winningElevator)
			

			orders[floor][button] = 1
			fmt.Println(orders)

		}else{
			msg_handler.Send_newOrder(floor, msg_handler.ButtonType(button), winningElevator)
			fmt.Println("Winning elevator has been calculated")
			fmt.Printf("winning elevator : %i \n", winningElevator)
		}
		
	}else{
		fmt.Println("forwarding")
		score := calculateOrderScore(floor, button)
		elev_score = append(elev_score, score)
		msg_handler.Send_requestedOrderEvaluation(elev_score, floor,msg_handler.ButtonType(button),elev_id)
	}
}



func setNFloors()(int){
	return msg_handler.NumberUserInput("Enter Amount of floors")
}
func setOrder()(int, int){
	return msg_handler.NumberUserInput("floor Order"),msg_handler.NumberUserInput("Button Order")
}



//func TestElevator(){
	/*for{
		if(InitDone){
			break
		}
	}*/
	/*fmt.Println("whats going on")
	//msg_handler.SemaphoreMessage <- 1
	nFloors := setNFloors()
	Event_init( nFloors )
	miliseconds := 40
	timer_miliseconds := time.Duration(miliseconds)*time.Millisecond
	const broadcastIP string = "129.241.187.255"
	const broadcastPort string = "20022"
//	msg_handler.LocalIP = netw.GetLocalIP()
//	msg_handler.LocalPort = netw.GetPort()
	broadcastAddr := broadcastIP+":"+broadcastPort
	brdcast_conn := netw.GetConnectionForDialing(broadcastAddr)
	msg_handler.TestSetNeighBourElevConnection(brdcast_conn)
	var floorb, buttontypeb int 
	for{
		time.Sleep(timer_miliseconds)
		if(orders[floorb][buttontypeb] != 1 ){
			floorb,buttontypeb =setOrder()	
			//queue.Orders[floorb][buttontypeb] = 1
			//broadcast_msg :=msg_handler.Message{MsgID : msg_handler.OrderRequest,TargetID : 1, Floor : floorb, Buttontype : buttontypeb }
			msg_handler.Send_newOrder(floorb,msg_handler.ButtonType(buttontypeb), 1)
			fmt.Println("Message sent")
			bufio.NewReader(os.Stdin).ReadBytes('\n') 
		}



	}

}*/