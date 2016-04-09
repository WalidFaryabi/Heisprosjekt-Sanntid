package queue
import(
        ."fmt"
)


type Button_type int
const (
	BUTTON_CALL_UP Button_type = iota //declares type with 0 and increments for each new variable
	BUTTON_CALL_DOWN
	BUTTON_COMMAND

)

type queueCommands[] Button_type // queueorder commands u dont need it.
var Orders[][] int
var Last_floor int
var Last_direction int
var n_floors int
//var Ordersorders [][]int
var MainFloor int
func Queue_init(n_total_floors int){
       // queueSystems := make([]queue
        //var quetest queueCommands
        //quetest = []Button_type{0,0,0}
        for j := 0 ; j<n_total_floors; j++{
         Orders = append(Orders, []int {0,0,0})
       }
      
       
      // slice1 := []int{0,0,0}
       //for k:= 0; k<n_total_floors ; k++{
       	//	Ordersorders = append(Ordersorders,[]int{0,0,0})
      // }

        n_floors = n_total_floors
        Last_floor = 0
        Last_direction = 1      
        //Commands,len(queueSystem),n_floors)
        //queueSystem = queueSystems
        
        Println(Orders)

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
	if(Orders[Last_floor][BUTTON_COMMAND] == 1){
	        MainFloor = Last_floor
	        return
	}

	//First rake up amount of floor orders in eacch direction
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


