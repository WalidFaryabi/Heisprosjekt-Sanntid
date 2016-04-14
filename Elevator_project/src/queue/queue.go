package queue
import(
        ."fmt"
        "math"
)


var Button_type int
const (
	BUTTON_CALL_UP int= iota //declares type with 0 and increments for each new variable
	BUTTON_CALL_DOWN
	BUTTON_COMMAND
)

type CheckQueueOption int  //Different options for how you want to use the checkQueueList function. Total orders? Orders in upper direction? etc.
const(
	BUTTON_UP_ORDERS CheckQueueOption=iota // 0
	BUTTON_DOWN_ORDERS 	// 1
	BUTTON_COMMAND_ORDERS //2
	ALL_ORDERS  //3
//	BUTTON_DOWN_ORDERS_UP 
	/*BUTTON_DOWN_ORDERS_DOWN = 0x05
	BUTTON_COMMAND_UP = 0x06
	BUTTON_COMMAND_DOWN = 0x07
	NO_DIRECTION = 0x10			//if you want to check amount of orders without specifying direction, then OR this with the desired command BUTTON_COMMAND_UP | NO_DIRECTION FOR INSTANCE
	// last identifier _UP/_DOWN the direction */
)
var Orders[][] int
var last_floor int
var last_direction int
var n_floors int

var order_highestOrder int
var order_lowestOrder int

var MainFloor int
func Queue_init(n_total_floors int)(){
        for j := 0 ; j<n_total_floors; j++{
         Orders = append(Orders, []int {0,0,0})
    	}
        n_floors = n_total_floors
        last_floor = 0
        last_direction = 1        
        Println(Orders) // confirmation that the order queue equals the correct amount of floors.
        
}

//Algorithm for choosing the next queue
func SetNextMainFloor(){
	//count amount of orders in the same direction it previously stopped
	//orders_up :=0
	//orders_down := 0
	//order_highestFloorOrder := 0 //highest floor order
	//order_lowestFloorOrder :=  0//lowest floor order
//	top_current_difference := n_floors - Last_floor //Amount of elevators between current and top floor
//	first_current_difference := Last_floor - 0 //Amount of elevators between current and first floor
	
	if((Orders[n_floors-1][BUTTON_COMMAND] == 1 || Orders[n_floors-1][BUTTON_CALL_DOWN] == 1) && (last_floor == 0)){
	        Println("1")
	        MainFloor = n_floors - 1
	        evaluateDirection()
	        return
	}else if((Orders[0][BUTTON_COMMAND] == 1 || Orders[0][BUTTON_CALL_UP] == 1) && (last_floor == n_floors-1)){
		MainFloor = 0	
		Println("2")
		evaluateDirection()
		return
	}
	if(CheckQueueList(ALL_ORDERS, 0,1) == 1){
		Println("3")
		for floor := 0; floor <n_floors; floor++{
			for button := 0; button < 3 ; button++{
				if(Orders[floor][button] == 1){
					Println("4")
					MainFloor = floor
					if(last_floor > floor){
						last_direction = -1
					}else if (last_floor < floor){
						last_direction = 1
					}else{
						last_direction = 0		
					}
					evaluateDirection()
					return
				}
			}
		}	
	}
	Println("5")

	//First rake up amount of floor orders in each direction
	orders_total_up := CheckQueueList(BUTTON_UP_ORDERS, last_floor,1)
	orders_total_up+= CheckQueueList(BUTTON_COMMAND_ORDERS,last_floor,1)
	orders_total_down := CheckQueueList(BUTTON_DOWN_ORDERS, last_floor,-1)
	orders_total_down+= CheckQueueList(BUTTON_COMMAND_ORDERS, last_floor,-1)
	
	var highestFloor int = 0
	var lowestFloor int = 0
	
	for floor := last_floor; floor >= 0; floor--{
				if(Orders[floor][BUTTON_CALL_DOWN] == 1 || Orders[floor][BUTTON_COMMAND] == 1){
					lowestFloor = floor
				}
	}
	
	for floor := last_floor; floor < n_floors; floor++{
				if(Orders[floor][BUTTON_CALL_UP] == 1 || Orders[floor][BUTTON_COMMAND] == 1){
					highestFloor = floor
				}
	}
	
	
	
	
	if(orders_total_up >= orders_total_down){
		if(last_direction == 1 || last_direction == 0){
			MainFloor = highestFloor
			
			
		
		}else if(last_direction == -1){
			if(last_floor == 0){
				MainFloor = highestFloor
			}else if (lowestFloor !=0){
				MainFloor = lowestFloor
			}else{
				MainFloor = highestFloor
			}
		}	
	}else{
	
		if(last_direction == -1 || last_direction == 0 ){
			
			MainFloor = lowestFloor
			
		}else if(last_direction == 1){	
			if(last_floor == n_floors-1){
				MainFloor = lowestFloor

			}else if(highestFloor != 0){
				MainFloor = highestFloor
			}else{
				MainFloor = lowestFloor
			}
			
		}
	}
	
	
	/* for j := last_floor ; j < n_floors ; j++ {
	 	orders_up+= 1			//Orders[j][BUTTON_CALL_UP] + Orders[j][BUTTON_COMMAND] // count up relevant orders.
	        if(Orders[j][BUTTON_CALL_UP] == 1 || Orders[j][BUTTON_COMMAND] == 1){
	        	order_highestFloorOrder = j         	                         
	        }
     }
      for j := last_floor; j>=0 ; j--{
      	orders_down+= Orders[j][BUTTON_CALL_DOWN] + Orders[j][BUTTON_COMMAND]
                if(Orders[j][BUTTON_CALL_DOWN] == 1 || Orders[j][BUTTON_COMMAND] == 1 ){
   					order_lowestFloorOrder = j    
                }
        }

    if( last_floor != n_floors-1 && last_floor != 0){
    	Println("6")
        if(orders_up > orders_down){
        		MainFloor = order_highestFloorOrder
                /*elevatorGap := order_highestFloorOrder - last_floor //Amount of elevators between current and highest floor order
                 if(top_current_difference < first_current_difference){ // If there are more up orders in upper direction while still having less floors away from top floor, then prioritize this
                        MainFloor = order_highestFloorOrder
                 
                 } else if(elevatorGap < (last_floor - order_lowestFloorOrder )){ 
                        MainFloor = order_highestFloorOrder
                 } else{
                       MainFloor = order_lowestFloorOrder
                                
                 }   */
       /* }else{
        
                MainFloor = order_lowestFloorOrder
        }
    }else{ // on first or last floor
    	Println("7")
    	if(last_floor == 0){
    		MainFloor = order_highestFloorOrder

    	}else{
    		MainFloor = order_lowestFloorOrder
    	}


    }*/

    evaluateDirection()


}


func evaluateDirection(){
    if(MainFloor > last_floor){
 
            last_direction = 1
    }else if(MainFloor < last_floor){
            last_direction = -1
    } else{
    	last_direction = 0
    }
}

func CheckQueueList(desiredOption CheckQueueOption,startFloor, direction int)(int){ 	//this function will return amount of orders of a given button command or all commands, given on a startfloor, 
	
	total_orders := 0														// in the relevant direction. Direction is either -1 or 1.
	lastFloorToCount := ( n_floors + (n_floors * direction) ) / 2 		//This will set the floor to either 0 or top floor. Reducing if statements.
	conditionalOffset := (1 - direction)/-2 				//if direction is -1, then this will be -1. if it is 1 then it will be 0. will use to offset correct conditions for for loop.	
	//Printf("%i , %i, %i",startFloor,lastFloorToCount,conditionalOffset) Highest/Lowest floor order will also be returned
	
	var highestLowestOrder int = 0
	
	switch(desiredOption){
		case ALL_ORDERS:
			for floor:= startFloor; (floor != lastFloorToCount + conditionalOffset) ;floor+=direction{
				for button:= BUTTON_CALL_UP; button<= BUTTON_COMMAND; button++{
					//Printf("%i \n",button)
					if(Orders[floor][button] ==1){
						total_orders++
						highestLowestOrder = floor
					}
				}
		
			}
			
			//return total_orders//,highestLowestOrder
		case BUTTON_UP_ORDERS,BUTTON_DOWN_ORDERS, BUTTON_COMMAND_ORDERS:
			for floor := startFloor ; floor != lastFloorToCount + conditionalOffset ; floor+=direction{
				if(Orders[floor][desiredOption] == 1){
					total_orders++
					highestLowestOrder = floor
				}
			}
			//return total_orders//, highestLowestOrder
		
			
	}
	if(direction == 1){
		order_highestOrder = highestLowestOrder
	}else if (direction == -1){
		order_lowestOrder = highestLowestOrder
	}
	if(total_orders > 0){
		//Println("returns total orders")
		//Println(total_orders)
		return total_orders
	}else{
		//Println("Re2")
		return -1
	}
}

func CalculateOrderScore(floor int, button int)(float64){ // algorithm for calculating whether an elevator should take an order or not.
	var priorityScore float64 = 0

	nCurrentOrders := CheckQueueList(ALL_ORDERS,0,1)
	if(nCurrentOrders == -1){
		priorityScore +=100
		if(last_floor == floor){
			priorityScore +=200 
			return priorityScore
		}else{
			priorityScore+=  150/(math.Abs(float64(last_floor-floor) ) )
			return priorityScore 

		}
	}
	
	if(Orders[floor][button] == 1){
		priorityScore+=50
	}
	if(last_floor == floor){
		priorityScore+=20
	}else if(last_floor < floor){
		if(MainFloor > floor){
			priorityScore+=45
			if(button == BUTTON_CALL_UP){
				priorityScore+=155
			}
		}else if(MainFloor < floor && last_direction == 1){
			priorityScore+=10
			if(button == BUTTON_CALL_UP){
				priorityScore+=60
			}else{
				priorityScore+=30
			}
		}else if(MainFloor < floor && last_direction == -1){
			priorityScore+=(100/(float64(floor-MainFloor) )) // maximum point given by that is 50
		}

	}else{ //Last_floor > floor
		if(MainFloor < floor){
			priorityScore+=45
			if(button == BUTTON_CALL_DOWN){
				priorityScore+=155
			}
		}else if(MainFloor > floor && last_direction == -1){
			priorityScore+=10
			if(button == BUTTON_CALL_DOWN){
				priorityScore+=60
			}else{
				priorityScore+=30
			}
		}else if(MainFloor > floor && last_direction == 1){
			priorityScore+=(100/(float64(MainFloor - floor) ))
			
		}
	}
	Println("priority score achieved?")
	Println(priorityScore)
	return priorityScore

}

func GetNFloors()(int){
	return n_floors
}

func GetLastDirection()(int){
	return last_direction
}

func GetLastFloor()(int){
	return last_floor
}

func SetLastFloor(currentFloor int){
	last_floor = currentFloor

}
