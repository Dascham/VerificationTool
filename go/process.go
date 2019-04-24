package main

import "sync"

//need process abstraction, consists of local variables
type Process struct {
	LocalVariables  map[string]int
	InitialLocation Location
	Mutex sync.Mutex //global mutex
}
//essentially, do reachability
//this can be run in go routine
func (p Process) Explore() {
	//do things
	//check at either beginning or end whether next (global) state has been seen before

}
