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

var Orders[][] int
var Last_floor int
var Last_direction int
var n_floors int

var MainFloor int
func Queue_init(n_total_floors int)(){
        for j := 0 ; j<n_total_floors; j++{
         Orders = append(Orders, []int {0,0,0})
    	}
        n_floors = n_total_floors
        Last_floor = 0
        Last_direction = 1        
        Println(Orders) // confirmation that the order queue equals the correct amount of floors.
        
}

//Algorithm for choosing the next queue
func SetNextMainFloor(){
	//count amount of orders in the same direction it previously stopped
	orders_up :=0
	orders_down := 0
	order_highestFloorOrder := 0 //highest floor order
	order_lowestFloorOrder :=  0//lowest floor order
//	top_current_difference := n_floors - Last_floor //Amount of elevators between current and top floor
//	first_current_difference := Last_floor - 0 //Amount of elevators between current and first floor
	if((Orders[n_floors-1][BUTTON_COMMAND] == 1 || Orders[n_floors-1][BUTTON_CALL_DOWN] == 1) && (Last_floor == 0)){
	        MainFloor = n_floors - 1
	        return
	}else if((Orders[0][BUTTON_COMMAND] == 1 || Orders[0][BUTTON_CALL_UP] == 1) && (Last_floor == n_floors-1)){
		MainFloor = 0	
		return
	}
	if(CheckQueueList() == 1){
		for floor := 0; floor <n_floors; floor++{
			for button := 0; button < 3 ; button++{
				if(Orders[floor][button] == 1){
					MainFloor = floor
					if(Last_floor > floor){
						Last_direction = -1
					}else if (Last_floor < floor){
						Last_direction = 1
					}else{
						Last_direction = 0		
					}
					return
				}
			}
		}	
	}

	//First rake up amount of floor orders in each direction
	 for j := Last_floor ; j < n_floors ; j++ {
	 	orders_up+= Orders[j][BUTTON_CALL_UP] + Orders[j][BUTTON_COMMAND] // count up relevant orders.
	        if(Orders[j][BUTTON_CALL_UP] == 1 || Orders[j][BUTTON_COMMAND] == 1){
	        	order_highestFloorOrder = j         	                         
	        }
     }
      for j := Last_floor; j>=0 ; j--{
      	orders_down+= Orders[j][BUTTON_CALL_DOWN] + Orders[j][BUTTON_COMMAND]
                if(Orders[j][BUTTON_CALL_DOWN] == 1 || Orders[j][BUTTON_COMMAND] == 1 ){
   					order_lowestFloorOrder = j    
                }
        }

    if( Last_floor != n_floors-1 && Last_floor != 0){
        if(orders_up > orders_down){
        		MainFloor = order_highestFloorOrder
                /*elevatorGap := order_highestFloorOrder - Last_floor //Amount of elevators between current and highest floor order
                 if(top_current_difference < first_current_difference){ // If there are more up orders in upper direction while still having less floors away from top floor, then prioritize this
                        MainFloor = order_highestFloorOrder
                 
                 } else if(elevatorGap < (Last_floor - order_lowestFloorOrder )){ 
                        MainFloor = order_highestFloorOrder
                 } else{
                       MainFloor = order_lowestFloorOrder
                                
                 }   */
        }else{
        
                MainFloor = order_lowestFloorOrder
        }
    }else{ // on first or last floor
    	if(Last_floor == 0){
    		MainFloor = order_highestFloorOrder

    	}else{
    		MainFloor = order_lowestFloorOrder
    	}


    }

    if(MainFloor > Last_floor){
            Last_direction = 1
    }else if(MainFloor < Last_floor){
            Last_direction = -1
    } else{
    	Last_direction = 0
    }


}

func CheckQueueList()(int){
	total_orders := 0
	for i:= 0; i<4;i++{
		for k:= 0; k<3; k++{
			if(Orders[i][k] ==1){
				total_orders++
			}
		}
	
	}
	if(total_orders > 0){
		return total_orders	
	}
	return -1
}

func CalculateOrderScore(floor int, button int)(float64){ // algorithm for calculating whether an elevator should take an order or not.
	var priorityScore float64 = 0

	nCurrentOrders := CheckQueueList()
	if(nCurrentOrders == -1){
		priorityScore +=100
		if(Last_floor == floor){
			priorityScore +=200 
			return priorityScore
		}else{
			priorityScore+=  150/(math.Abs(float64(Last_floor-floor) ) )
			return priorityScore 

		}
	}
	if(Orders[floor][button] == 1){
		priorityScore+=50
	}
	if(Last_floor == floor){
		priorityScore+=20
	}else if(Last_floor < floor){
		if(MainFloor > floor){
			priorityScore+=45
			if(button == BUTTON_CALL_UP){
				priorityScore+=155
			}
		}else if(MainFloor < floor && Last_direction == 1){
			priorityScore+=10
			if(button == BUTTON_CALL_UP){
				priorityScore+=60
			}else{
				priorityScore+=30
			}
		}else if(MainFloor < floor && Last_direction == -1){
			priorityScore+=(100/(float64(floor-MainFloor) )) // maximum point given by that is 50
		}

	}else{ //Last_floor > floor
		if(MainFloor < floor){
			priorityScore+=45
			if(button == BUTTON_CALL_DOWN){
				priorityScore+=155
			}
		}else if(MainFloor > floor && Last_direction == -1){
			priorityScore+=10
			if(button == BUTTON_CALL_DOWN){
				priorityScore+=60
			}else{
				priorityScore+=30
			}
		}else if(MainFloor > floor && Last_direction == 1){
			priorityScore+=(100/(float64(MainFloor - floor) ))
			
		}
	}
	return priorityScore

}

func GetNFloors()(int){
	return n_floors
}


