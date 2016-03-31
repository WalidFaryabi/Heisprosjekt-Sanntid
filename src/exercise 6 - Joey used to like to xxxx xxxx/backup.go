
package main

import(
	"fmt"
	"os/exec"
	"net"
	"encoding/json"
	"time" 
)

const PORT = "20021"

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

func getListenConnection() (*net.UDPConn){
	addr,_ := net.ResolveUDPAddr("udp", "localhost:" + PORT)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("error occured in establishing connection")
		fmt.Println("%s", err)
	}
	return conn
}

func getDialConnection() (*net.UDPConn) {
	addr, err := net.ResolveUDPAddr("udp", "localhost:" + PORT)
	if err != nil {
		fmt.Println("error occured in parsing ip addr")
		fmt.Println("%s", err)
	}
	conn, err := net.DialUDP("udp", nil,addr)	
	if err != nil {
		fmt.Println("error occured in establishing connection for sending")
		fmt.Println("%s", err)
	}
	return conn
}


func main(){
	
	/* BACKUP */	
	
	current_state := 0	
	conn := getListenConnection()
	
	var m Message	
	buffer := make([]byte,1024) // TO BE USED FOR RECEIVING

	for {
		_= conn.SetDeadline(time.Now().Add(time.Second*2))
		sizeOfBuffer, err := conn.Read(buffer)
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
	
	exc_cmd("./backup")
	conn1 := getDialConnection()

	for {
		// SENDING MESSAGE
		message := Message{"Im alive!", current_state}
		buffer1, err := json.Marshal(message)
		if err != nil {
			fmt.Println("error occured in marshal")
			fmt.Println("%s", err)
		}
		_,_ = conn1.Write(buffer1)
		
		current_state+=1;
		fmt.Println(current_state)
	}
}
