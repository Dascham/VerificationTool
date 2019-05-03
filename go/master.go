package main

import (
	"strconv"
)

var portNumbers []string = MakeRange(1, 10)
const localHost string = "127.0.0.1"
const hello = 1

func SendState(){

}

func ReceiveState()State{
	var state State = State{}

	return state
}
//flood localhost for now, for instances of workers
func ContactWorkers(){

}


func MakeRange(min, max int) []string {
	a := make([]string, max-min+1)
	for i := range a {
		a[i] = strconv.Itoa(min + i)
	}
	return a
}