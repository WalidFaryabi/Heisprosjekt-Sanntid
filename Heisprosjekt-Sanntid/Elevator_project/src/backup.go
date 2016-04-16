
package main

import(
	"fmt"
	"os/exec"
	"net"
	"encoding/json"
	"time" 
)

type Message struct {
	Msg string
    	State int
}

func exc_cmd(path string) {
	fmt.Println("path is ", path)
	err := exec.Command("gnome-terminal", "-x", "sh", "-c", path).Run()
	if err != nil {
		fmt.Println("error occured")
		fmt.Printf("%s", err)	
	}
}

func getListenConnection(port int) UDPConn{
	p := 20000+port
	addr,_ := net.ResolveUDPAddr("udp", ":" + string(p))
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("error occured in establishing connection")
		fmt.Println("%s", err)
	}
	return conn
}

func getDialConnection(port int) UDPConn {
	p := 20000+port
	conn, err := net.Dial("udp", ":"+string(p))
	if err != nil {
		fmt.Println("error occured in establishing connection for sending")
		fmt.Println("%s", err)
	}
	return conn
}


func main{
	
	/* BACKUP */	
	
	current_state = 0	
	conn := getListenConnection(21)
	
	var m Message	
	buffer := make([]byte,1024) // TO BE USED FOR RECEIVING

	for {
		_= conn.SetDeadline(time.Now().Add(time.Millisecond*100))
		sizeOfBuffer, err := conn.Read(buffer1)
		if err != nil {
			fmt.Println("error in conn.Read")
			fmt.Println(err)
			conn.Close()
			break; // Primary is dead
		}
		if (sizeOfBuffer != 0) {
			_ = json.Unmarshal(buffer[:sizeOfBuffer], &m)
			fmt.Println(m.Msg)
			current_state = m.State
		}
	}
	
	/* MAKE BACKUP PRIMARY */
	
	conn := getDialConnection(21)
	
	exc_cmd("./backup") 

	for {
		// SENDING MESSAGE
		message := Message{"Im alive!", current_state}
		buffer1, err := json.Marshal(start_message)
		if err != nil {
			fmt.Println("error occured in marshal")
			fmt.Println("%s", err)
		}
		_,_ = sending_conn.Write(buffer1)
		
		current_state+=1;
		fmt.Println(current_state)
	}
	
	
}
