package FSM_Elevator

import(
	//"fmt"
	"../elev_driver"
	"../queue"
	"time"
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

func Event_init(init_floor int){
	elev_driver.Elev_init() // this should be tested during phase 1 initiliazation possibly. Also send a value indicating it was not properly initiated?
	if(init_floor == 0){
		elev_driver.Elev_set_motor_direction(-1)
		for{
		        if(elev_driver.Elev_get_floor_sensor_signal() != 0){
		        break
		        }
		}
		elev_driver.Elev_set_motor_direction(0) // c
	}
	state_slave = IDLE

}

func Event_queueNotEmpty(){

        switch state_slave{
               case IDLE:
                    queue.SetNextMainFloor()
                    elev_driver.Elev_set_motor_direction(elev_driver.Elev_motor_direction_t(queue.Last_direction))
        }

	queue.SetNextMainFloor()
	elev_driver.Elev_set_motor_direction(elev_driver.Elev_motor_direction_t(queue.Last_direction))
	state_slave = MOVING
}

func elevator_floorInQueue(floor int){

	elev_driver.Elev_set_motor_direction(0)
	elev_driver.Elev_set_floor_indicator(floor)
	if(floor != queue.MainFloor){
	        elev_driver.Elev_set_button_lamp(elev_driver.Elev_button_type_t(queue.Last_direction),floor,0)
	} // change this..
	elev_driver.Elev_set_button_lamp(elev_driver.BUTTON_COMMAND,floor,0)
	elev_driver.Elev_set_door_open_lamp(1)
	state_slave = DOOR_OPEN
	timer := time.NewTimer(time.Second * 3 )
        <- timer.C
        //this is bad
       elev_driver.Elev_set_door_open_lamp(0)

}


// Non-event functions creating the structure of the elevator system
