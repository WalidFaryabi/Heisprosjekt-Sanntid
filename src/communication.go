
package main
import ("fmt"
		"time"
		"net"
		"bufio"
		"os"
		"strconv"
		"strings"
		"./netw"
		"encoding/json"
)
//var broadcastPort int = getPort()
var listenPort int = getPort()
var neighborElevatorAddress string = ""
var addressOfDetectedElevator string = ""
//var elevatorID int
var numberOfElevators int
const broadcastIP string = "129.241.187.255"
const broadcastPort string = "20022"

type Message struct {
	Msg string
	Addr string
	NumberOfElevators int
}

func getPort() (int) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter listen port: ")
	portString, _ := reader.ReadString('\n')
	portString = strings.Replace(portString, "\n", "", -1)
	port,err := strconv.Atoi(portString)

	if err != nil {
		fmt.Println("ERROR IN CONVERSION")
		fmt.Printf("%s \n", err)
		return 0
	}
	return port
}

func isNeighborElevatorAddressEmtpy()(bool) {
	if (neighborElevatorAddress != "") {
		return false
	}
	return true
}

func broadcast(conn *net.UDPConn) {
	t0 := time.Now()
	run_bc := true
	msg := "Hello!"
	addr := netw.GetLocalIP()+":"+strconv.Itoa(20000+listenPort)
	
	broadcast_msg := Message{msg, addr, 0}

	for {
		if(run_bc){
			buffer,err := json.Marshal(broadcast_msg)
	
			if err != nil {
				fmt.Println("ERROR IN MARSHAL")
				fmt.Println("%s", err)
			}

			_,_ = conn.Write(buffer)

			t1 := time.Now()
			if(t1.Sub(t0) > 2*(1000*time.Millisecond)){
				run_bc = false
			}
		} else {
			break
		}
		
	}
	conn.Close()
}

func dial(address string) {

	dialConn := netw.GetConnectionForDialing(address)
	msg := "Hello!"
	addr := netw.GetLocalIP()+":"+strconv.Itoa(20000+listenPort)
	
	go func() {
		for {
			num := numberOfElevators
			dial_msg := Message{msg, addr, num}
			buffer,err := json.Marshal(dial_msg)
	
			if err != nil {
				fmt.Println("ERROR IN MARSHAL")
				fmt.Println("%s", err)
			}
			_,err = dialConn.Write(buffer)

			if err != nil {
				fmt.Println("ERROR IN DIALING")
				fmt.Println(err)		
			}
		}
	}()
}

func listen(conn *net.UDPConn, buffer []byte) {
	addr := ""
	t0 := time.Now()
	var message Message
	for {	
			if (neighborElevatorAddress == ""){
				n, err := conn.Read(buffer)
				t1 := time.Now()
				if(t1.Sub(t0) > 2*(1000*time.Millisecond)){
					numberOfElevators = 1
					fmt.Println("NUMBER OF ELEVATORS: 1")
				}

				if err != nil {
					fmt.Println("ERROR IN READING MESSAGE")
					fmt.Println(err)
				}

				if n != 0 {
					_ = json.Unmarshal(buffer[:n], &message)
					fmt.Println(message.Msg, message.Addr[len(message.Addr)-2:], message.NumberOfElevators)
					addr = message.Addr
					neighborElevatorAddress = addr
					numberOfElevators+=1
				}

			} else {
				n, err := conn.Read(buffer)
				if err != nil {
					fmt.Println("ERROR IN READING MESSAGE")
					fmt.Println(err)
				}
				if n != 0 {
					_ = json.Unmarshal(buffer[:n], &message)
					
					if(message.NumberOfElevators > numberOfElevators) {
						numberOfElevators = message.NumberOfElevators
					}
					fmt.Println(message.Msg, message.Addr[len(message.Addr)-2:], message.NumberOfElevators)
					addr = message.Addr
					if (addr != "" && neighborElevatorAddress != addr && addressOfDetectedElevator != addr) {
						fmt.Println("WE HAVE DETECTED A NEW ELEVATOR WITH ADDRESS: ", addr)
						addressOfDetectedElevator = addr
						numberOfElevators+=1
					}
				}		
			}
		}
}

func main() {
	
	broadcastAddr := broadcastIP+":"+broadcastPort
	brdcast_conn := netw.GetConnectionForDialing(broadcastAddr)
	
	broadcast(brdcast_conn)
	
	listenIP := ""
	listenAddr := listenIP + ":" + strconv.Itoa(20000+listenPort) 
	listen_conn := netw.GetConnectionForListening(listenAddr)
	
	buffer := make([]byte, 1024)
	
	go listen(listen_conn, buffer)
	
	for {
		if(!isNeighborElevatorAddressEmtpy()) {
			break		
		}
	}
	
	dial(neighborElevatorAddress)
	
    var input string
    fmt.Scanln(&input)

}
