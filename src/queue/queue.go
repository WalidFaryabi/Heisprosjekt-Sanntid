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

type queueCommands[] Button_type // queueorder commands
var Orders[] queueCommands
var Last_floor int
var Last_direction int
var n_floors int
var MainFloor int
func Queue_init(n_total_floors int){
       // queueSystems := make([]queue
        var quetest queueCommands
        quetest = []Button_type{0,0,0}
        for j := 0 ; j<n_total_floors; j++{
         Orders = append(Orders, quetest)
       }
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
	top_current_difference := n_floors - Last_floor //Amount of elevators between current and top floor
	first_current_difference := Last_floor - n_floors //Amount of elevators between current and first floor
	if(Orders[Last_floor][BUTTON_COMMAND] == 1){
	        MainFloor = Last_floor
	        return
	}
        if( Last_floor != n_floors-1 || Last_floor != 0){
	        for j := Last_floor ; j < n_floors ; j++ {
		        if(Orders[j][BUTTON_CALL_UP] == 1){
		                orders_up++
		                if(Orders[j][BUTTON_COMMAND] == 1){
		                        order_highestFloorOrder = j        
		                }
		        }
	        }
	        for j := Last_floor; j>=0 ; j++{
	                if(Orders[j][BUTTON_CALL_DOWN] == 1){
	                        orders_down++
	                        if(Orders[j][BUTTON_COMMAND] == 1){
	                                order_lowestFloorOrder = j
	                        }
	                }
	        }
	        if(orders_up > orders_down){
	                elevatorGap := order_highestFloorOrder - Last_floor //Amount of elevators between current and highest floor order
	                 if(top_current_difference < first_current_difference){ // If there are more up orders in upper direction while still having less floors away from top floor, then prioritize this
	                        MainFloor = order_highestFloorOrder
	                 
	                 } else if(elevatorGap < (Last_floor - order_lowestFloorOrder )){ 
	                        MainFloor = order_highestFloorOrder
	                 } else{
	                       MainFloor = order_lowestFloorOrder
	                                
	                 }   
	        }else{
	        
	                MainFloor = order_lowestFloorOrder
	        }
        }
        if(MainFloor > Last_floor){
                Last_direction = 1
        }else{
                Last_direction = -1
        }


}


