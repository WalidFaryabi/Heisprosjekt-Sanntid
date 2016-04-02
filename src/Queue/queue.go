package main
import(
        ."fmt"
)


type button_type int
const (
	BUTTON_CALL_UP button_type = iota //declares type with 0 and increments for each new variable
	BUTTON_CALL_DOWN
	BUTTON_COMMAND

)

type queueCommands[] button_type // queueorder commands
var orders[] queueCommands
var last_floor int
var Last_direction int
var n_floors int
var MainFloor int
func queue_init(n_total_floors int){
       // queueSystems := make([]queue
        var quetest queueCommands
        quetest = []button_type{0,0,0}
        for j := 0 ; j<n_total_floors; j++{
         orders = append(orders, quetest)
       }
        n_floors = n_total_floors
        last_floor = 0
        last_direction = 1      
        //Commands,len(queueSystem),n_floors)
        //queueSystem = queueSystems
        
        Println(orders)

}

//Algorithm for choosing the next queue
func setNextMainFloor(){
	//count amount of orders in the same direction it previously stopped
	orders_up :=0
	orders_down := 0
	order_highestFloorOrder := 0 //highest floor order
	order_lowestFloorOrder :=  0//lowest floor order
	top_current_difference := n_floors - last_floor //Amount of elevators between current and top floor
	first_current_difference := last_floor - n_floors //Amount of elevators between current and first floor
        if( last_floor != n_floors-1 || last_floor != 0){
	        for j := last_floor ; j < n_floors ; j++ {
		        if(orders[j][BUTTON_CALL_UP] == 1){
		                orders_up++
		                if(orders[j][BUTTON_COMMAND] == 1){
		                        order_highestFloorOrder = j        
		                }
		        }
	        }
	        for j := last_floor; j>=0 ; j++{
	                if(orders[j][BUTTON_CALL_DOWN] == 1){
	                        orders_down++
	                        if(orders[j][BUTTON_COMMAND] == 1){
	                                order_lowestFloorOrder = j
	                        }
	                }
	        }
	        if(orders_up > orders_down){
	                elevatorGap := order_highestFloorOrder - last_floor //Amount of elevators between current and highest floor order
	                 if(top_current_difference < first_current_difference){ // If there are more up orders in upper direction while still having less floors away from top floor, then prioritize this
	                        MainFloor = order_highestFloorOrder
	                 
	                 } else if(elevatorGap < (last_floor - order_lowestFloorOrder )){ 
	                        MainFloor = order_highestFloorOrder
	                 } else{
	                       MainFloor = order_lowestFloorOrder
	                                
	                 }   
	        }else{
	        
	                MainFloor = order_lowestFloorOrder
	        }
        }
        if(MainFloor > lastfloor){
                last_direction = 1
        }else{
                last_direction = -1
        }


}

func main(){
        queue_init(5)
}
