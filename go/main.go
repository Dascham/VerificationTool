package main

import (
	"fmt"
	"sync"
)

//need process abstraction, consists of local variables
type Process struct {
	localVariables  map[string]int
	initialLocation Location
	mutex sync.Mutex //global mutex
}
//essentially, do reachability
//this can be run in go routine
func (p Process) Explore() {
	//do things
	//check at either beginning or end whether next (global) state has been seen before
}

type Location struct{
	locationName string
	locationId int
	edge []Edge //should be made into slice
	invariant Invariant
}
//an edge consists of src node, dst node, a possible guard, a possible synchronization,
// and a possible update on local or global variables.
type Edge struct {
	src *Location
	dst *Location
	guard Guard
	ch chan int
	isSend bool
	update Update
}

func (e Edge) EdgeIsActive(localVariables map[string]int){
	if (e.guard.Evaluate(localVariables) &&
		e.dst.invariant.IsValid(localVariables)) { //add one more for chan

	}
}
type Guard struct {
	guardValue int
	comparisonOperator string
	variableToEvaluate string
}
//a function that only applies to guards
func (g Guard) Evaluate(localVariables map[string]int) bool{
	var result bool = false
	switch g.comparisonOperator {
	case "<":
		if(localVariables[g.variableToEvaluate] < g.guardValue){
			result = true
		}
	case ">":
		if(localVariables[g.variableToEvaluate] > g.guardValue){
			result = true
		}
	default:
		result = false
	}
	return result
}

type Invariant struct {
	invariantValue int
	comparisonOperator string
	variableToEvaluate string
}

func (i Invariant) IsValid(localVariables map[string]int) bool {
	result := false
	switch i.comparisonOperator {
	case "<":
		if (i.invariantValue < localVariables[i.variableToEvaluate]) {
			result = true
		}
	case ">":
		if (i.invariantValue > localVariables[i.variableToEvaluate]){
			result = true
		}
	}
	return result
}
type Update struct {
	updateValue int
	variableToUpdate string
}
func (u Update) Update(localVariables map[string]int){
	localVariables[u.variableToUpdate] = u.updateValue
}













func main(){
	//have global mutex, in order to change global state


	fmt.Printf("Hello, world!")
	localVariables := make(map[string]int)
	localVariables["b"] = 2
	var a Guard
	a.guardValue = 5
	a.variableToEvaluate = "b"
	a.comparisonOperator = "<"
	fmt.Println(a.Evaluate(localVariables))

}


