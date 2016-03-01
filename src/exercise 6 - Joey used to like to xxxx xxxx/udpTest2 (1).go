package main
import ("fmt"
	"os"
	"os/exec"
	//"./src/network" 
	"net"
	"encoding/json" 
)

type Message struct {
	Msg string
    	Identifier int
}

func exc_cmd(path string) {
	fmt.Println("path is ", path)
	err := exec.Command("gnome-terminal", "-x", "sh", "-c", path).Run()
	if err != nil {
		fmt.Println("error occured")
		fmt.Printf("%s", err)	
	}
}


func main(){
	//CREATE UDPADDRESS
	ipAddress := ""
	port := "20021"
	addr, err := net.ResolveUDPAddr("udp", ipAddress+":" + port)
	if err != nil {
		fmt.Println("error occured")
		fmt.Println("%s", err)
	}

	//ETSTABLISH CONNECTION FOR LISTENING
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("error occured")
		fmt.Println("%s", err)
	}
		
	buffer := make([]byte,1024)	
	//_= localSocket.SetDeadline(time.Now().Add(time.Millisecond*100))
	
	var m Message
	counter := 0
	last_int := 0
	//READ MESSAGES
	for {	
		n, _ := conn.Read(buffer)
		err := json.Unmarshal(buffer[:n], &m) 
		if err != nil {
			fmt.Println("error in unmarshal")
			fmt.Println(err)
		}
		if (m.Msg !=  "") {
			counter = 0
			//fmt.Printf("%s, %d \n", m.Msg, m.Identifier)
			last_int = m.Identifier
		} else {
			if(counter >= 10000){
				//fmt.Println("MAIN PROGRAM CRAHSED")
				//last_int++
				//fmt.Printf("BACKUP: %d \n", last_int)
				fmt.Printf("LAST INTEGER: %d \n", last_int)
				break	
			}
			counter++
		}
	}
	conn.Close()
	
	//LAUNCH MAINPROGRAM AGAIN WITH LAST RECORDED INTEGER
	
	//RUN MAIN PROGRAM	
	exc_cmd("./udpTest")
	
	//ESTABLISH CONNECTION FOR SENDING
	conn1, err := net.Dial("udp", ":"+port)
	if err != nil {
		fmt.Println("error occured")
		fmt.Println("%s", err)
	}
	
	i:=0
	for {
		restart_message := Message{"RESTART", last_int}
		buffer, err := json.Marshal(restart_message)
		if err != nil {
			fmt.Println("error occured")
			fmt.Println("%s", err)
		}
		_,_ = conn1.Write(buffer)
		if(i >= 10000) {
			break
		}
		i++
	}
	fmt.Println("THIS GETS CALLED")
	conn1.Close()
	os.Exit(1)
	
}


