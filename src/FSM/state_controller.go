package FSM

import(
	"fmt"
	"../elev_driver"
	"../queue"
	"time"
	"../netw"
	"../msg_handler"
	"../queue"
)


type state_slaveElevator int
const (
	INITIALIZATION state_slaveElevator = iota
	IDLE  
	MOVING 
	DOOR_OPEN 
	DOOR_CLOSED
	NOT_INITIALIZATED
)


func Thread_elevatorStateMachine(C_elevatorCommand chan int,C_order chan ch_elevOrder){		//Make a channel int in main
	Event_init()
	msgType := -1
	notSingleElevator := false
	for{
		//running elevator
		if(msg_handler.GetNelevators() > 1){
			notSingleElevator = true
		}
		else{
			notSingleElevator = false
		}
		msg_handler.GetNelevators
		for floor := 0; i<4;i++{
			for buttontype :=0;k < 3; k++{
				if(floor == 0 && buttontype == 1){
					continue
				}else if(floor == 3 && buttontype == 0){
					continue	
				}
				if((buttontype != 2 && elev_driver.Elev_get_button_signal(elev_driver.Elev_button_type_t(buttontype),floor) == 1) && notSingleElevator ){			//outside button pressed
					Event_OutsideButtonPressed(floor, buttontype)
				}		
				else if(elev_driver.Elev_get_button_signal(elev_driver.Elev_button_type_t(buttontype),floor) == 1){							//inside button pressed
					Event_newQueueRequest(floor,queue.Button_type(buttontype))					

				}
				if(queue.Orders[floor][buttontype] == 1){
					Event_queueNotEmpty()
				}
			}
		}
		currentfloor := currentFloor()
		//fmt.Println(currentfloor)
		elev_driver.Elev_set_floor_indicator(2)
		if(currentfloor >= 0 && currentfloor<4){
			elev_driver.Elev_set_floor_indicator(currentfloor)
			for i:= 0; i <3;i++{
				if(queue.Orders[currentfloor][i] == 1){
					Event_floorInQueue(currentfloor)
					fmt.Println(currentfloor)
					//fmt.Println(queue.Orders)
					//fmt.Printf("Queue at: %i", currentfloor)
				}
			}
		}
		if(!notSingleElevator){
			continue
			//skip the last part because it isonly for multiple elevators.
		}
		select{
			case msgrecv := <-C_elevator:
				msgType = msgrecv
			default:
				msgTy = -1 
		}
		if(msgType != -1){
			switch(msgType){
			case 1:
				order := <-C_order
				Event_evaluateRequest(order.Floor, order.Button)
				//Recv Order request evaluation
			case 2:
				order := <-C_order
				Event_newQueueRequest(order.Floor, order.Button)
				//Recv Order request
			case 3:
				//Recv elevator status from all elevators
			}


		}
	}

}

func Thread_elevatorCommRecv(C_elevatorCommand chan int,C_order chan ch_elevOrder){//assume it has made a connection for listening. This requires a connection first.
	//PAPPA is joey a faggot?
	conn := msg_handler.GetNeighbourConnection()
	buffer := make([]byte, 1024)
	var message msg_handler.Message
	for{
		n,err := conn.Read(buffer)
		if (err != nil){
			fmt.Println("ERROR IN READING ELEVATOR MESSAGE")
		fmt.Println(err)
		}
		if(n != 0){
			_ = json.Unmarshal(buffer[:n],&message)
		
			switch(message.MsgID){	
			case msg_handler.OrderRequestEvaluation:
				C_elevatorCommand <- 1
				C_order.Floor <- message.Floor
				C_order.Buttontype <-message.Buttontype
			case msg_handler.OrderRequest:
				if(message.Elev_targetID == msg_handler.GetID())
					C_elevatorCommand <- 2
					C_order.Floor <- message.Floor
					C_order.Buttontype <- message.Buttontype
				}else{
					msg_handler.Send_newOrder(message.Floor, message.Buttontype, message.Elev_targetID)	
				}
			case msg_handler.Elevator_initializationStatus:
					//something.
				if(message.elev_id == msg_handler.GetID()){
					//This elevator was the one initializing the command for status.
					if(message.SuccessfullInit == true){
						//All elevators are in working order fuck yeah! send nothing.
					}else{
						//Not all elevators are in working order. Must send a message to every elevator with the failed target id. If it doesnt work
						//after a second
						//message.Elev_failedID 

					}
				}
			}
		}
}




var state_slave state_slaveElevator

func Event_init(){
	init_success := elev_driver.Elev_init() // this should be tested during phase 1 initiliazation possibly. Also send a value indicating it was not properly initiated?
	if(init_success == 0){
		fmt.Println("Initialization failed")
		state_slave = NOT_INITIALIZATED
		Send_elevInitCompleted(false)
		return
	}
	queue.Queue_init(4) // we take 4 floors currently
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
	Send_elevInitCompleted(true)
	fmt.Println("Init done")
}

func Event_queueNotEmpty(){

        switch state_slave{
               case IDLE:
               	    fmt.Println("Queue not empty")
               	    //queue.MainFloor = 3
                    queue.SetNextMainFloor()
               	    //queue.MainFloor = 2
                    fmt.Printf("Main floor: %i", queue.MainFloor)
                    elev_driver.Elev_set_motor_direction(elev_driver.Elev_motor_direction_t(queue.Last_direction))
                    //elev_driver.Elev_set_motor_direction(1)
                    
                    state_slave = MOVING
        }
}

func Event_floorInQueue(floor int){ // this is activated if one of the orders are on the floor. Goes through every kind of order.
	queue.Last_floor = floor
	switch(state_slave){
	case MOVING:
		fmt.Println("Floor in queue called, moving state.")
		if(queue.Orders[floor][2] == 1){
			fmt.Println("Sup boy")
			elev_driver.Elev_set_motor_direction(0)
			queue.Orders[floor][2] = 0
			fmt.Println(queue.Orders)
			elev_driver.Elev_set_button_lamp(elev_driver.BUTTON_COMMAND,floor,0)
		}
		if( (queue.Last_direction == 1) && (queue.Orders[floor][elev_driver.BUTTON_CALL_UP] == 1) ){
				elev_driver.Elev_set_motor_direction(0)
				queue.Orders[floor][0] = 0
				elev_driver.Elev_set_button_lamp(elev_driver.BUTTON_CALL_UP,floor,0)
		} else if( (queue.Last_direction == -1) && (queue.Orders[floor][elev_driver.BUTTON_CALL_DOWN] == 1) ){
				elev_driver.Elev_set_motor_direction(0)
				queue.Orders[floor][1] = 0
				elev_driver.Elev_set_button_lamp(elev_driver.BUTTON_CALL_DOWN,floor,0)
		}
		elev_driver.Elev_set_floor_indicator(floor)

		elev_driver.Elev_set_door_open_lamp(1)
		state_slave = DOOR_OPEN
		//timer := time.NewTimer(time.Second * 3 )
		dur := 3
		something := time.Duration(dur)*time.Second
		time.AfterFunc(something,Event_doorTimeout)
		//<- timer.C use this to indicate timer done
	  	  break
	}
}

func Event_doorTimeout(){
	fmt.Println("This was called")
	switch(state_slave){
	case DOOR_OPEN:
		elev_driver.Elev_set_door_open_lamp(0)
		if(queue.MainFloor == queue.Last_floor){
			state_slave = IDLE
			fmt.Println("MAin floor reached")
		}else{
			elev_driver.Elev_set_motor_direction(elev_driver.Elev_motor_direction_t(queue.Last_direction))
			state_slave = MOVING
		}
	}

}
func Event_newQueueRequest(floor int, button queue.Button_type){
// add order regardless of current state
	switch(state_slave){
	case NOT_INITIALIZATED:
			return
	}
	fmt.Printf("queueRequest in floor: %d \n", floor+1)
	queue.Orders[floor][button] = 1
	fmt.Println(queue.Orders)
	elev_driver.Elev_set_button_lamp(elev_driver.Elev_button_type_t(button), floor,1)
}

func Event_outsideButtonPressed(floor int, button queue.Button_type){
	switch(state_slave){
	case NOT_INITIALIZATED:
			return
	}
	//Accepting a button outside pressed does not matter in which state we are. All orders will be accessed.
	score := queue.CalculateOrderScore(floor, button)
	score_array := []int{score}
	msg_handler.send_requestedOrder(score_array, floor,button)			
}


// ADD LIGHTS TO THE BUTTONS
func Event_evaluateRequest(requestedOrder msg_orderRequestEvaluation){
	switch(state_slave){
	case NOT_INITIALIZATED:
			return
	}
	if(elev_id == requestedOrder.elev_id){
		highestElevatorScore := 0
		winningElevator := 0
		for i := 0 ; i<msg_handler.GetNelevators() ;i++{
			if(requestedOrder.elev_score[i] >= highestElevatorScore){
				highestElevatorScore = requestedOrder.elev_score[i]
				winningElevator = i + 1
			}
		}
		if(elev_id == winningElevator ){
			//this elevator is the winner
			queue.Orders[requestedOrder.Floor]requestedOrder.Buttontype] = 1
		}else{
			Send_newOrder(requestedOrder.Floor, requestedOrder.Button_type, winningElevator)
		}
		break
	}

	else{
		score := queue.CalculateOrderScore(requestedOrder.Floor, requestedOrder.Buttontype)
		requestedOrder.Elev_score = append(requestedOrder.Elev_score, score)
		msg_handler.send_requestedOrder(requestedOrder.Elev_score, requestedOrder.Floor,requestedOrder.Button)
	}

}




