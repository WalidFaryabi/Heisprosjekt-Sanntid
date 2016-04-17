package FSM
import(
        ."fmt"
        "math"
)

/***************************************************************************************************************************\\\\\\\\
This package contains the elevator system.

In this file, the queue is written. It is only accessed by state controller.

It can be configured for more than 4 floors.

*****************************************************************************************************************************//////

/***************************	***************private variables*****************************************/
var orders[][] int 				//Queue matrix slice. Is dynamically allocated.
var lastFloor int
var lastDirection int
var nFloors int

var orderHighestOrder int
var orderLowestOrder int

var mainFloor int 				// Contains the floor that the elevator moves to before being in idle state.

//potentially remove this
var buttonType int
const (
	BUTTON_CALL_UP int= iota //declares type with 0 and increments for each new variable
	BUTTON_CALL_DOWN
	BUTTON_COMMAND
)

type checkQueueOption int  //Different options for how you want to use the checkQueueList function. Total orders? Orders in upper direction? etc.
const(
	BUTTON_UP_ORDERS checkQueueOption=iota // 0
	BUTTON_DOWN_ORDERS 	// 1
	BUTTON_COMMAND_ORDERS //2
	ALL_ORDERS  //3
)

func queue_init(nTotalFloors int)(){
        for j := 0 ; j<nTotalFloors; j++{
         orders = append(orders, []int {0,0,0})
    	}
        nFloors = nTotalFloors
        lastFloor = 0
        lastDirection = 1        
        Println(orders) // confirmation that the order queue equals the correct amount of floors.
        
}

//Algorithm for setting the next main floor
func setNextMainFloor(){
	//count amount of orders in the same direction it previously stopped
	//orders_up :=0
	//orders_down := 0
	//order_highestFloorOrder := 0 //highest floor order
	//order_lowestFloorOrder :=  0//lowest floor order
//	top_current_difference := n_floors - Last_floor //Amount of elevators between current and top floor
//	first_current_difference := Last_floor - 0 //Amount of elevators between current and first floor
	
	if((orders[nFloors-1][BUTTON_COMMAND] == 1 || orders[nFloors-1][BUTTON_CALL_DOWN] == 1) && (lastFloor == 0)){
	        Println("1")
	        mainFloor = nFloors - 1
	        evaluateDirection()
	        return
	}else if((orders[0][BUTTON_COMMAND] == 1 || orders[0][BUTTON_CALL_UP] == 1) && (lastFloor == nFloors-1)){
		mainFloor = 0	
		Println("2")
		evaluateDirection()
		return
	}
	if(checkQueueList(ALL_ORDERS, 0,1) == 1){
		Println("3")
		for floor := 0; floor <nFloors; floor++{
			for button := 0; button < 3 ; button++{
				if(orders[floor][button] == 1){
					Println("4")
					mainFloor = floor
					if(lastFloor > floor){
						lastDirection = -1
					}else if (lastFloor < floor){
						lastDirection = 1
					}else{
						lastDirection = 0		
					}
					evaluateDirection()
					return
				}
			}
		}	
	}
	Println("5")

	//First rake up amount of floor orders in each direction
	ordersTotalUp := checkQueueList(BUTTON_UP_ORDERS, lastFloor,1)
	ordersTotalUp+= checkQueueList(BUTTON_COMMAND_ORDERS,lastFloor,1)
	ordersTotalDown := checkQueueList(BUTTON_DOWN_ORDERS, lastFloor,-1)
	ordersTotalDown+= checkQueueList(BUTTON_COMMAND_ORDERS, lastFloor,-1)
	
	var highestFloor int = 0
	var lowestFloor  int = 0
	
	for floor := lastFloor; floor >= 0; floor--{
				if(orders[floor][BUTTON_CALL_DOWN] == 1 || orders[floor][BUTTON_COMMAND] == 1){
					lowestFloor = floor
				}
	}
	
	for floor := lastFloor; floor < nFloors; floor++{
				if(orders[floor][BUTTON_CALL_UP] == 1 || orders[floor][BUTTON_COMMAND] == 1){
					highestFloor = floor
				}
	}
	
	
	
	
	if(ordersTotalUp >= ordersTotalDown){
		if(lastDirection == 1 || lastDirection == 0){
			mainFloor = highestFloor
			
			
		
		}else if(lastDirection == -1){
			if(lastFloor == 0){
				mainFloor = highestFloor
			}else if (lowestFloor !=0){
				mainFloor = lowestFloor
			}else{
				mainFloor = highestFloor
			}
		}	
	}else{
	
		if(lastDirection == -1 || lastDirection == 0 ){
			
			mainFloor = lowestFloor
			
		}else if(lastDirection == 1){	
			if(lastFloor == nFloors-1){
				mainFloor = lowestFloor

			}else if(highestFloor != 0){
				mainFloor = highestFloor
			}else{
				mainFloor = lowestFloor
			}
			
		}
	}
	
    evaluateDirection()


}

func evaluateDirection(){
    if(mainFloor > lastFloor){
 
            lastDirection = 1
    }else if(mainFloor < lastFloor){
            lastDirection = -1
    } else{
    	lastDirection = 0
    }
}

func checkQueueList(desiredOption checkQueueOption,startFloor, direction int)(int){ 	//this function will return amount of orders of a given button command or all commands, given on a startfloor, 
	
	total_orders := 0														// in the relevant direction. Direction is either -1 or 1.
	lastFloorToCount := ( nFloors + (nFloors * direction) ) / 2 		//This will set the floor to either 0 or top floor. Reducing if statements.
	conditionalOffset := (1 - direction)/-2 				//if direction is -1, then this will be -1. if it is 1 then it will be 0. will use to offset correct conditions for for loop.	
	//Printf("%i , %i, %i",startFloor,lastFloorToCount,conditionalOffset) Highest/Lowest floor order will also be returned
	
	var highestLowestOrder int = 0
	
	switch(desiredOption){
		case ALL_ORDERS:
			for floor:= startFloor; (floor != lastFloorToCount + conditionalOffset) ;floor+=direction{
				for button:= BUTTON_CALL_UP; button<= BUTTON_COMMAND; button++{
					//Printf("%i \n",button)
					if(orders[floor][button] ==1){
						Println("We somehow have a order.")
						total_orders++
						highestLowestOrder = floor
					}
				}
		
			}
			
			//return total_orders//,highestLowestOrder
		case BUTTON_UP_ORDERS,BUTTON_DOWN_ORDERS, BUTTON_COMMAND_ORDERS:
			for floor := startFloor ; floor != lastFloorToCount + conditionalOffset ; floor+=direction{
				if(orders[floor][desiredOption] == 1){
					total_orders++
					highestLowestOrder = floor
				}
			}
			//return total_orders//, highestLowestOrder
		
			
	}
	if(direction == 1){
		orderHighestOrder = highestLowestOrder
	}else if (direction == -1){
		orderLowestOrder = highestLowestOrder
	}
	if(total_orders > 0){
		return total_orders
	}else{
		Println("no orders were in the queue")
		return -1
	}
}

func calculateOrderScore(floor int, button int)(float64){ // algorithm for calculating whether an elevator should take an order or not.
	var priorityScore float64 = 0

	nCurrentOrders := checkQueueList(ALL_ORDERS,0,1)
	if(nCurrentOrders == -1){
		Println("we should get here definitely")
		priorityScore +=100
		if(lastFloor == floor){
			priorityScore +=200 
			return priorityScore
		}else{
			priorityScore+=  150/(math.Abs(float64(lastFloor-floor) ) )
			return priorityScore 

		}
	}
	
	if(orders[floor][button] == 1){
		priorityScore+=50
	}
	if(lastFloor == floor){
		priorityScore+=40
		if(BUTTON_CALLL_UP == button && lastDirection == 1 || BUTTON_CALL_DOWN ==button && lastDirection == -1 ){
			priorityScore+=100
		}
	}else if(lastFloor < floor){
		if(mainFloor > floor){
			priorityScore+=45
			if(button == BUTTON_CALL_UP){
				priorityScore+=155
			}
		}else if(mainFloor < floor && lastDirection == 1){
			priorityScore+=10
			if(button == BUTTON_CALL_UP){
				priorityScore+=60
			}else{
				priorityScore+=30
			}
		}else if(mainFloor < floor && lastDirection == -1){
			priorityScore+=(100/(float64(floor-mainFloor) )) // maximum point given by that is 50
		}

	}else{ //Last_floor > floor
		if(mainFloor < floor){
			priorityScore+=45
			if(button == BUTTON_CALL_DOWN){
				priorityScore+=155
			}
		}else if(mainFloor > floor && lastDirection == -1){
			priorityScore+=10
			if(button == BUTTON_CALL_DOWN){
				priorityScore+=60
			}else{
				priorityScore+=30
			}
		}else if(mainFloor > floor && lastDirection == 1){
			priorityScore+=(100/(float64(mainFloor - floor) ))
			
		}
	}
	Println("priority score achieved?")
	Println(priorityScore)
	return priorityScore

}

func GetNFloors()(int){
	return nFloors
}

func GetLastDirection()(int){
	return lastDirection
}

func GetLastFloor()(int){
	return lastFloor
}

func GetMainFloor()(int){
	return mainFloor
}

func SetLastFloor(currentFloor int){
	lastFloor = currentFloor

}
