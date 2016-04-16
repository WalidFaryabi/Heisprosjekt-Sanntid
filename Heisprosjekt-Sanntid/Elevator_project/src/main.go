package main
import (
	"fmt"
	"./FSM"
	//"./src/network" 
	//"net"
	//"encoding/json" 
	"./elev_driver"
	"./queue"
)


func checkQueueList()(int){
	for i:= 0; i<4;i++{
		for k:= 0; k<3; k++{
			if(queue.Orders[i][k] ==1){
				return 1
			}
		}
	}
	return -1
}

func currentFloor()(int){
	for i:= 0; i<4;i++{
		if(elev_driver.Elev_get_floor_sensor_signal() == i){
			return i
		}
	}
	return -1
}

//timer afterfunc 
func main(){
//first we have to initialize the elevator of course
	FSM.Event_init()
	fmt.Println(queue.Orders)
	for{
		//check button inputs and queue orders
		for i := 0; i<4;i++{
			for k :=0;k < 3; k++{
				if(i == 0 && k == 1){
					continue
				}else if(i == 3 && k == 0){
					continue	
				}
				if(elev_driver.Elev_get_button_signal(elev_driver.Elev_button_type_t(k),i) == 1){
					FSM.Event_newQueueRequest(i,queue.Button_type(k))

					
				}
				if(queue.Orders[i][k] == 1){
					FSM.Event_queueNotEmpty()
					
				}
			}
		}
		//fmt.Println("stfu")
		//fmt.Println(queue.Orders)

		//queue.Orders[2][1] = 1
				//fmt.Println(queue.Orders)
//				fmt.Println(queue.Ordersorders)


	

		currentfloor := currentFloor()
		//fmt.Println(currentfloor)
		elev_driver.Elev_set_floor_indicator(2)
		if(currentfloor >= 0 && currentfloor<4){
			elev_driver.Elev_set_floor_indicator(currentfloor)
			for i:= 0; i <3;i++{
				if(queue.Orders[currentfloor][i] == 1){
					FSM.Event_floorInQueue(currentfloor)
					fmt.Println(currentfloor)
					//fmt.Println(queue.Orders)
					//fmt.Printf("Queue at: %i", currentfloor)
				}
			}
		} 
	
	}

}
