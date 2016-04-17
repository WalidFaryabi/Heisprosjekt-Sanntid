package main


import(
	"fmt"
	"time"
)

var timer *time.Timer
func main(){
	//define a timer.
	timer_seconds := time.Duration(3)*time.Second
	timer =time.AfterFunc(timer_seconds,callback)
	go stopTimer()
	for{}
	timer.Reset(timer_seconds)
}


func callback(){
	fmt.Println("This was called hoembro")
	timer.Reset(time.Duration(3)*time.Second)

}

func stopTimer(){
	for{
	time.Sleep(4 * time.Second)
	timer.Reset(time.Duration(3) * time.Second)
	}
}