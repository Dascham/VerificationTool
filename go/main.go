package main

import (
	"fmt"
)
func main(){
	//have global mutex, in order to change global state

	fmt.Printf("Hello, world!")
	localVariables := make(map[string]int)
	localVariables["b"] = 2
	var a Guard
	a.GuardValue = 5
	a.VariableToEvaluate = "b"
	a.ComparisonOperator = "<"
	fmt.Println(a.Evaluate(localVariables))

	//instantiate three processes, then model check

}


func Explore(map[string]Process) {

}


