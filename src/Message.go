package Message

import (
	"fmt"
	"net"
	"encoding/json"
)

type messageType int
const (
	I_AM_ALIVE messageType = iota
	STATE_INFO
	COMMAND
	CONFIG
)

type Message struct {
	msgType messageType;
	text string
	n_elev int
	n_floors int
	*ports int
	*orders int
}

func SendMessage(msg Message, conn *net.UDPConn) {
	// WE SHOULD CONSIDER USING GO CHANNELS
	switch msg.MessageType {
		case I_AM_ALIVE:
			// SEND IM ALIVE MESSAGES TO SLAVES (OR MASTER)
			_,_ = conn.Write([]byte("Im Alive"))
		case STATE_INFO:
			// SEND ORDERS, CURRENTFLOOR AND CURRENTDIRECTION TO MASTER
		case COMMAND:
			// SEND COMMAND TO SLAVE
		case CONFIG:
			// SEND CONFIG INFO TO ALL SLAVES
		default:
			fmt.Println("Undefined case")
	}
}

func ReadMessage(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	

}


