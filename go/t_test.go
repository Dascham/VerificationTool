package main

import "testing"

/*
func TestcalcThing(t *testing.T){
	value := calcThing(5.5, 7)
	if value != 12.5{
		t.Fail()
		//t.Errorf("calcThing was incorrect, got: %f, but expected %f", value, 12.5)
	}
}
*/
func Setup() map[string]int {
	var localVariables map[string]int = make(map[string]int)
	localVariables["a"] = 2
	localVariables["b"] = 7
	return localVariables
}

func TestGuard_Evaluate(t *testing.T) {
	var localVariables map[string]int = Setup()
	var g Guard = Guard{3, "<", "a"}
	if g.Evaluate(localVariables) != true{
		t.Errorf("Evaluate function should return true, but returned false")
	}
}
func TestGuard_Evaluate2(t *testing.T) {
	var localVariables map[string]int = Setup()
	var g Guard = Guard{3, "<", "b"}
	if g.Evaluate(localVariables) == true{
		t.Errorf("Evaluate function should return false, but returned true")
	}
}
