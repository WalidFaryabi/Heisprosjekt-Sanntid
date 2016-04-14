package elev_driver

import (
	"fmt"
)

const N_FLOORS int = 4
const N_BUTTONS int = 3
const MOTOR_SPEED int = 2800

var Elev_motor_direction_t int
const (
	DIRN_DOWN int = -1
	DIRN_STOP int = 0
	DIRN_UP int = 1
)
var Elev_button_type_t int
const (
	BUTTON_CALL_UP int = iota //declares type with 0 and increments for each new variable
	BUTTON_CALL_DOWN
	BUTTON_COMMAND

)

var lamp_channel_matrix = [N_FLOORS][N_BUTTONS] int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_channel_matrix = [N_FLOORS][N_BUTTONS] int {
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

func elev_checkLegalFloors(button int, floor int)(int){
	if(floor <0){
		fmt.Println("YOU ARE ASKING BELOW FLOOR 0 HOW? FATAL ERROR")
		return 1	
	}
	if(floor >= N_FLOORS){
		fmt.Println("YOU ARE ASKING FOR A FLOOR ABOVE THE MAXIMUM? WTF? FATAL ERROR")
		return 1
	}	
	if(button < 0 || button > BUTTON_COMMAND || ( (button == BUTTON_CALL_UP) && (floor == 3) ) || ( (button == BUTTON_CALL_DOWN) && (floor == 0))) {
		fmt.Println("You are asking for a button that does not exist.")
		return 1
	}
	return 0
}

func Elev_init()(int){
	init_success := Io_init()

	if(init_success == 0){
		fmt.Println("unsuccessfull elev init")
		return 0
	}
	for  f := 0; f<N_FLOORS; f++{
		for b := BUTTON_CALL_UP ; b <= BUTTON_COMMAND; b++{
			if(	f == 0 && b == 1 || (f == N_FLOORS -1) && b == 0){
				continue			
			}
			Elev_set_button_lamp(b,f,0)		
		}
	}
	Elev_set_stop_lamp(0)
	Elev_set_door_open_lamp(0)
	Elev_set_floor_indicator(0)
	return 1
}

func Elev_set_motor_direction(dirn int){
	if(dirn == 0){
		Io_write_analog(MOTOR,0)
	}else if (dirn > 0){
		Io_clear_bit(MOTORDIR)
		Io_write_analog(MOTOR,MOTOR_SPEED)
	}else if (dirn < 0){
		Io_set_bit(MOTORDIR)
		Io_write_analog(MOTOR,MOTOR_SPEED)
	}
	
}

func Elev_set_button_lamp(button,floor,value int){
	if(elev_checkLegalFloors(button,floor) == 1){
		return 
	}
	if (value==1){
		Io_set_bit(lamp_channel_matrix[floor][button])
	}else{
		Io_clear_bit(lamp_channel_matrix[floor][button])
	}	
	
}

func Elev_set_floor_indicator(floor int){
	if(floor < 0 || floor > 3){
		fmt.Println("Fatal error, non existing floor")
		return
	}
	if((floor & 0x02) == 1){
		Io_set_bit(LIGHT_FLOOR_IND1)
		
	}else{
		Io_clear_bit(LIGHT_FLOOR_IND1)
	}
	if((floor & 0x01 ) == 1){
		Io_set_bit(LIGHT_FLOOR_IND2)
	}else{
		Io_clear_bit(LIGHT_FLOOR_IND2)
	}
	
}

func Elev_set_door_open_lamp(value int){
	if(value == 1){
		Io_set_bit(LIGHT_DOOR_OPEN)
	}else{
		Io_clear_bit(LIGHT_DOOR_OPEN)
	}

}

func Elev_set_stop_lamp(value int){
	if(value == 1){
		Io_set_bit(LIGHT_STOP)
	}else{
		Io_clear_bit(LIGHT_STOP)
	}
}

func Elev_get_button_signal(button , floor int) (int){
	if(elev_checkLegalFloors(button,floor) == 1){
		fmt.Println("%i, %i 	\n",floor, button)
		return 0
	}
	if(Io_read_bit(button_channel_matrix[floor][button]) ==1){
		return 1	
	}else{
		return 0
	}
}


func Elev_get_floor_sensor_signal()(int){
	if(Io_read_bit(SENSOR_FLOOR1) == 1){
		return 0
	}else if(Io_read_bit(SENSOR_FLOOR2) == 1){
		return 1
	}else if(Io_read_bit(SENSOR_FLOOR3) == 1){
		return 2
	}else if(Io_read_bit(SENSOR_FLOOR4) == 1){
		return 3
	}else{
		return -1
	}
}

func Elev_get_stop_signal()(int){
	return Io_read_bit(STOP)
}

func Elev_get_obstruction_signal()(int){
	return Io_read_bit(OBSTRUCTION)
}
