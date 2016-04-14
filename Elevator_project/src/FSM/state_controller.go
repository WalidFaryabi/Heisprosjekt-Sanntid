package FSM

import(
	"fmt"
	"../elev_driver"
	"../queue"
	"time"
	//"../netw"
	"../msg_handler"
	//"encoding/json"
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

func currentFloor()(int){
	for i:= 0; i<4;i++{
		if(elev_driver.Elev_get_floor_sensor_signal() == i){
			return i
		}
	}
	return -1
}

var state_slave state_slaveElevator
func Thread_elevatorStateMachine(C_elevatorCommand chan int,C_order chan msg_handler.Ch_elevOrder){		//Make a channel int in main
	fmt.Println("We get not here")
	Event_init(4)
	fmt.Println("We get here")
	miliseconds := 40
	timer_miliseconds := time.Duration(miliseconds)*time.Millisecond
	msgType := -1
	notSingleElevator := true
	n_floors := queue.GetNFloors() 
	prevfloor,prevbuttontype := 0,0
	for{
		time.Sleep(timer_miliseconds)
		//running elevator
		if(msg_handler.GetNelevators() > 1){
			notSingleElevator = true
		}else{
			//notSingleElevator = false
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
				if((buttontype != elev_driver.BUTTON_COMMAND && elev_driver.Elev_get_button_signal(buttontype,floor) == 1) && false){// && notSingleElevator ){			//outside button pressed					
					fmt.Println("This is where u fucked up walid")
					fmt.Printf(" %i		%i 	\n", floor,buttontype)
					Event_outsideButtonPressed(floor,buttontype)
				}else if(elev_driver.Elev_get_button_signal(buttontype,floor) == 1 && queue.Orders[floor][buttontype] == 0){							//inside button pressed
					Event_newQueueRequest(floor,buttontype)					

				}
				if(queue.Orders[floor][buttontype] == 1  && prevfloor != floor && prevbuttontype != buttontype){
					prevfloor,prevbuttontype = floor,buttontype //to avoid excessive button mashing order...
					Event_queueNotEmpty()
				}
			}	
		}
		currentfloor := currentFloor()
		

		//elev_driver.Elev_set_floor_indicator(2)
		
		if(currentfloor >= 0 && currentfloor<4){
		
			elev_driver.Elev_set_floor_indicator(currentfloor)
			last_direction := queue.GetLastDirection()
			queue.SetLastFloor(currentfloor)
			
			for button:= elev_driver.BUTTON_CALL_UP; button <=elev_driver.BUTTON_COMMAND;button++{
				switch(button){
					case elev_driver.BUTTON_CALL_UP: 
						if( (last_direction == 1) && (queue.Orders[currentfloor][button] == 1)){
							Event_floorInQueue(currentfloor,button)
						}else if(last_direction == -1 && queue.Orders[currentfloor][button] == 1 && (queue.CheckQueueList(queue.ALL_ORDERS,0,1) == 1) || 															(queue.Orders[0][button] == 1 && currentfloor == 0) ){
							Event_floorInQueue(currentfloor,button)
						}
					case elev_driver.BUTTON_CALL_DOWN:
						if(last_direction == -1 && queue.Orders[currentfloor][button] == 1){
							Event_floorInQueue(currentfloor,1)
						}else if( (last_direction == 1) && (queue.Orders[currentfloor][button] == 1) && (queue.CheckQueueList(queue.ALL_ORDERS,0,1) == 1) || 								( (queue.Orders[n_floors - 1][button] == 1) && 				currentfloor == n_floors-1 )) {
							
							Event_floorInQueue(currentfloor,button)
						}
					case elev_driver.BUTTON_COMMAND:
						if(queue.Orders[currentfloor][button] == 1){
							Event_floorInQueue(currentfloor,button)
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
			case 1:
				order := <-C_order
				Event_evaluateRequest(order.Floor, int(order.Button),order.Elev_id, order.Elev_score)
				//Recv Order request evaluation
			case 2:
			
				order := <-C_order
				if(queue.Orders[order.Floor][order.Button] == 1){
						
				}else{
				Event_newQueueRequest(order.Floor,int(order.Button))
				}
				//Recv Order request
			case 3:
				//Recv elevator status from all elevators
			}


		}
	}

}
													




func Event_init(nTotalFloors int){
	init_success := elev_driver.Elev_init() // this should be tested during phase 1 initiliazation possibly. Also send a value indicating it was not properly initiated?
	if(init_success == 0){
		fmt.Println("Initialization failed")
		state_slave = NOT_INITIALIZATED
		//msg_handler.Send_elevInitCompleted(false)
		return
	}
	queue.Queue_init(nTotalFloors) // we take 4 floors currently
	
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
                    elev_driver.Elev_set_motor_direction((queue.GetLastDirection()))
                    //elev_driver.Elev_set_motor_direction(1)
                    if(queue.MainFloor == queue.GetLastFloor()){
						state_slave = IDLE					
					}else{
						state_slave = MOVING
					}
					
                    
        }
}

func Event_floorInQueue(floor int, button int){ // this is activated if one of the orders are on the floor.
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
		elev_driver.Elev_set_floor_indicator(floor)

		
		//timer := time.NewTimer(time.Second * 3 )
		
		//<- timer.C use this to indicate timer done
	case IDLE:
		fmt.Println("Floor in queue called, IDLE state") //this means an order was given while at the same floor.

	}
		elev_driver.Elev_set_button_lamp(button, floor, 0)
		queue.Orders[floor][button] = 0
		elev_driver.Elev_set_door_open_lamp(1)
		state_slave = DOOR_OPEN
		fmt.Println(queue.Orders)
		//timer := time.NewTimer(time.Second * 3 )
		seconds := 3
		timer_seconds := time.Duration(seconds)*time.Second
		time.AfterFunc(timer_seconds,Event_doorTimeout)

}

func Event_doorTimeout(){
	fmt.Println("Timeoutcalled")
	switch(state_slave){
	case DOOR_OPEN:
		elev_driver.Elev_set_door_open_lamp(0)
		if(queue.MainFloor == queue.GetLastFloor()){
			state_slave = IDLE
			fmt.Println("Main floor reached")
		}else{
			elev_driver.Elev_set_motor_direction(queue.GetLastDirection())
			state_slave = MOVING
		}
	}

	elev_driver.Elev_set_door_open_lamp(0)

}
func Event_newQueueRequest(floor int, button int){
// add order regardless of current state
	switch(state_slave){
	case NOT_INITIALIZATED:
			return
	}
	fmt.Printf("queueRequest in floor: %d \n", floor+1)
	queue.Orders[floor][button] = 1
	fmt.Println(queue.Orders)
	elev_driver.Elev_set_button_lamp(button, floor,1)
}

func Event_outsideButtonPressed(floor int, button int){
	switch(state_slave){
	case NOT_INITIALIZATED:
			return
	}
	//Accepting a button outside pressed does not matter in which state we are. All orders will be accessed.
	score := queue.CalculateOrderScore(floor, button)
	score_array := []float64{score}
	
	fmt.Println(score_array)
	//msg_handler.Send_requestedOrderEvaluation(score_array, floor,msg_handler.ButtonType(button) )			
}


// ADD LIGHTS TO THE BUTTONS
func Event_evaluateRequest(floor int,  button int, elev_id int, elev_score []float64){
	switch(state_slave){
	case NOT_INITIALIZATED:
			return
	}
	if(msg_handler.GetID()== elev_id){
		highestElevatorScore := 0.0
		winningElevator := 0
		for i := 0 ; i<msg_handler.GetNelevators() ;i++{
			if(elev_score[i] >= highestElevatorScore){
				highestElevatorScore = elev_score[i]
				winningElevator = i + 1
			}
		}
		if(elev_id == winningElevator ){
			//this elevator is the winner
			queue.Orders[floor][button] = 1
		}else{
			msg_handler.Send_newOrder(floor, msg_handler.ButtonType(button), winningElevator)
		}
		
	}else{
		score := queue.CalculateOrderScore(floor, button)
		elev_score = append(elev_score, score)
		msg_handler.Send_requestedOrderEvaluation(elev_score, floor,msg_handler.ButtonType(button))
	}
}




