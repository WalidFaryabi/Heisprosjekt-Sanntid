package main

import (
	"fmt"
	"time"
)

type testStruct struct{
	insignificantNr int
	msg string
}

func go1(c chan testStruct){
	for{
		select{
		case testMsg := <-c:
			fmt.Println(testMsg.Hei)
			fmt.Println(testMsg.msg)
		
		default:
			//fmt.Println("none received")
		}
	}
	
}

func go2(c chan testStruct){
	structtest := testStruct{insignificantNr : 0, msg : "Joey u failed"}
	i := 0
	for{
		time.Sleep(2 * time.Second)
		c <- structtest
		i++
		structtest.insignificantNr = i
	}

}
func main() {
	fmt.Println("Hello, playground")
	testChannel := make(chan testStruct)

	go go1(testChannel)
	go go2(testChannel)
	
	/* Run infinite loop */
	for{
	}
}
