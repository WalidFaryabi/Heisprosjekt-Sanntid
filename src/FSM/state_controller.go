package FSM

import(
	"fmt"
	"../elev_driver"
	"../queue"
	"time"
	"network"
	"msg_handler"
)


type state_slaveElevator int
const (
	INITIALIZATION state_slaveElevator = iota
	IDLE  
	MOVING 
	DOOR_OPEN 
	DOOR_CLOSED
)






var state_slave state_slaveElevator

func Event_init(){
	elev_driver.Elev_init() // this should be tested during phase 1 initiliazation possibly. Also send a value indicating it was not properly initiated?
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
	fmt.Printf("queueRequest in floor: %d \n", floor+1)
	queue.Orders[floor][button] = 1
	fmt.Println(queue.Orders)
	elev_driver.Elev_set_button_lamp(elev_driver.Elev_button_type_t(button), floor,1)
}

func Event_OutsideButtonPressed(floor int, button queue.Button_type){
	//Accepting a button outside pressed does not matter in which state we are. All orders will be accessed.
	score := queue.CalculateOrderScore(floor, button)
	score_array := []int{score}
	msg_handler.send_requestedOrder(score_array, floor,button)			
}
// ADD LIGHTS TO THE BUTTONS
func Event_evaluateRequest(requestedOrder msg_orderRequest){
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
func Send_newOrder(floor int, button ButtonType, chosenElevator int){