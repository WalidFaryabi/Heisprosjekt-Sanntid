package main
import ("fmt"
	//"os"
	"os/exec"
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

func restart()(int) {
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
	//conn1,_ := net.Dial("udp", ":"+port)
		
	buffer := make([]byte,1024)
	
	var m Message
	j:=0
	
	for j = 0;j<1000;j++{
		n, _ := conn.Read(buffer)
		err := json.Unmarshal(buffer[:n], &m) 
		if err != nil {
			fmt.Println("error in unmarshal")
			fmt.Println(err)
		}
		if(m.Msg == "RESTART") {
			fmt.Println("WE GOT TO RESTART")
			conn.Close()
			return m.Identifier
		}
	}
	conn.Close()
	return 0
}


func main(){
	
	j := restart()
	fmt.Printf("AFTER RESTART:%d \n",j)
	
	//CREATE UDPADDRESS
	ipAddress := ""
	port := "20021"
	_, err := net.ResolveUDPAddr("udp4", ipAddress + ":" + port)
	if err != nil {
		fmt.Println("error occured")
		fmt.Println("%s", err)
	}
	
	//RUN TEST PROGRAM	
	exc_cmd("./udpTest2")	
	
	//ESTABLISH CONNECTION FOR SENDING
	conn, err := net.Dial("udp", ":"+port)
	if err != nil {
		fmt.Println("error occured")
		fmt.Println("%s", err)
	}

	//SEND MESSAGES
	for {
		fmt.Println(j)
		m := Message{"Plass 21", j}
		
		buffer, err := json.Marshal(m)
		if err != nil {
			fmt.Println("error occured")
			fmt.Println("%s", err)
		}
		_,_ = conn.Write(buffer)

		j++
	}
}

