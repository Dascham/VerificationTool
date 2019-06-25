package main

import (
	"fmt"
	"time"
)
//v_0 global
//v_i local
//4x(4.3), 1x(4.2)
func Experiment1(){
	start := time.Now()
	fmt.Println("Experiment 1")
	s := EmptyState()
	s.allTemplates = make([]Template, 0,0)
	s.allTemplates = append(s.allTemplates, SetupTestTemplate2(), SetupTestTemplate3(),SetupTestTemplate3(), SetupTestTemplate3())
	s.globalVariables = map[string]int{"b":0}
	list,_ := Explore(s)
	fmt.Print("Number of states explored: ")
	fmt.Println(len(list))
	fmt.Println("Running time: "+time.Since(start).String())
}
//4x, 4.3 (updated 4.3, without sync and v_0)
func Experiment2(){
	start := time.Now()
	fmt.Println("Experiment 2")
	s := EmptyState()
	s.allTemplates = make([]Template, 0,0)
	s.allTemplates = append(s.allTemplates, SetupTestTemplate4_nosync(),SetupTestTemplate4_nosync(),SetupTestTemplate4_nosync(), SetupTestTemplate4_nosync())
	s.globalVariables = map[string]int{"b":0}
	list,_ := Explore(s)
	fmt.Print("Number of states explored: ")
	fmt.Println(len(list))
	fmt.Println("Running time: "+time.Since(start).String())
}
//4x, 4.1
func Experiment3(){
	start := time.Now()
	fmt.Println("Experiment 3")
	s := EmptyState()
	s.allTemplates = make([]Template, 0,0)
	s.allTemplates = append(s.allTemplates, SetupTestTemplate1(), SetupTestTemplate1(), SetupTestTemplate1(), SetupTestTemplate1())
	s.globalVariables = map[string]int{"b":0}
	list,_ := Explore(s)
	fmt.Print("Number of states explored: ")
	fmt.Println(len(list))
	fmt.Println("Running time: "+time.Since(start).String())
}
