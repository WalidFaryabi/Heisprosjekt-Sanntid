package FSM_Elevator

import(
	".fmt"
	"../Elev_driver"
	"../queue"
)


type state_slaveElevator int
const (
	INITIALIZATION state_slaveElevator = iota
	IDLE state_slaveElevator 
	MOVING state_slaveElevator
	DOOR_OPEN state_slaveElevator
	DOOR_CLOSED state_slaveElevator
)

var state_slave state_slaveElevator

func Event_init(init_floor int){
	elev_init // this should be tested during phase 1 initiliazation possibly. Also send a value indicating it was not properly initiated?
	if(init_floor == 0){
		Elev_driver.Elev_set_motor_direction(-1)
		while(Elev_driver.Elev_get_floor_sensor_signal() != 0);
		Elev_driver.Elev_set_motor_direction(0) // c
	}
	state_slave = IDLE

}

func Event_queueNotEmpty(){
	queue.SetNextMainFloor()
	Elev_driver.Elev_set_motor_direction(queue.Last_direction)
	state_slave = MOVING
}

func elevator_inFloor(floor int){
	queue.ClearRequest()
	Elev_driver.Elev_set_motor_direction(0)
	Elev_driver.Elev_set_button_lamp()

}
// Non-event functions creating the structure of the elevator system
