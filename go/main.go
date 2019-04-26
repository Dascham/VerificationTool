package main

import (
	"fmt"
)
func main(){
	//have global mutex, in order to change global state,
	// although not necessary for passed and waiting list implementation

	fmt.Printf("Hello, world!")
	localVariables := make(map[string]int)
	localVariables["b"] = 2
	var a Guard
	a.GuardValue = 5
	a.VariableToEvaluate = "b"
	a.ComparisonOperator = "<"
	fmt.Println(a.Evaluate(localVariables))

	//define the things, guards
	/*
	var guard1 Guard = Guard{"x", ">", 5}
	var guard2 Guard = Guard{"y", ">", 22}
	var guard3 Guard = Guard{"x", "==", 0}
	var guard4 Guard = Guard{"y", "==", 0}

	var update1 Update = Update{"y", "*=", 0}
	var update2 Update = Update{"x", "=", 0}
	var update3 Update = Update{"y", "=", 17}
	var update4 Update = Update{"x", "=", 0}
	var update5 Update = Update{"y", "=", 1}
	var Update6 Update = Update{"x", "++", 0}
	*/
	//instantiate 1 process, then model check

}

func Explore(map[string]Process) {

}


