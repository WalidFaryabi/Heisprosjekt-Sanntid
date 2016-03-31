package main
import (
	"fmt"

	//"./src/network" 
	//"net"
	//"encoding/json" 
	"elev_driver"
)



func main(){
	fmt.Println("Lol")
	elev_driver.Elev_init()
	elev_driver.Elev_set_door_open_lamp(1)
	elev_driver.Elev_set_motor_direction(-1)
	
	
}

