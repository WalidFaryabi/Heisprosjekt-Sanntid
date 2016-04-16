package main

import (
	"fmt"
	"time"
)

type moren struct{
	Hei int
	msg string
}

func go1(c chan moren){
	for{
		select{
		case bleble := <-c:
			fmt.Println(bleble.Hei)
			fmt.Println(bleble.msg)
		
		default:
			//fmt.Println("none received")
		}
	}
	
}

func go2(c chan moren){
	structtest := moren{Hei : 0, msg : "Joey u failed"}
	i := 0
	for{
		time.Sleep(2 * time.Second)
		c <- structtest
		i++
		structtest.Hei = i
	}

}
func main() {
	fmt.Println("Hello, playground")
	bla := make(chan moren)
	//structbla := moren{Hei : 1, msg : "jOEY LIKES TO SUCK HUUGE COCKULU"}
	
	go go1(bla)
	go go2(bla)
	
	for{
	}
}
