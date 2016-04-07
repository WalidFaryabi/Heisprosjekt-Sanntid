package main
import (
	"fmt"
	"FSM"
	//"./src/network" 
	//"net"
	//"encoding/json" 
	"elev_driver"
	"queue"
)


func checkQueueList()(int){
	for i:= 0; i<4;i++{
		for k:= 0; k<3; k++{
			if(queue.Orders[i][k] ==1){
				return 1
			}
		}
	}

}

func currentFloor()(int){
	for i:= 0; i<4;i++{
		if(elev_driver.Elev_get_floor_sensor_signal() == i){
			return i
		}
	}
	return -1
}


func main(){
//first we have to initialize the elevator of course
	FSM.Event_int()
	for{

		for i := 0; i<4;i++{
			for k :=0;k < 3; k++{
				if(elev_driver.Elev_get_button_signal(k,i) == 1){
					FSM.Event_newQueueRequest(i,k)
				}
			}
		}

		if(checkQueueList() == 1){
			FSM.Event_queueNotEmpty
		}
		currentfloor := currentFloor()
		if(currentFloor >= 0 || currentFloor<4){
			for i:= 0; i <3;i++{
				if(queue.Orders[currentFloor][i]){
					FSM.Event_floorInQueue
				}
			}
		}
	}
}


