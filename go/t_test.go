package main

import (
	"testing"
)

/*
func TestcalcThing(t *testing.T){
	value := calcThing(5.5, 7)
	if value != 12.5{
		t.Fail()
		//t.Errorf("calcThing was incorrect, got: %f, but expected %f", value, 12.5)
	}
}
*/
func SetupMap() map[string]int {
	var localVariables map[string]int = make(map[string]int)
	localVariables["a"] = 2
	localVariables["b"] = 7
	return localVariables
}

func SetupTemplate() Template{
	var template Template = Template{SetupMap(), Location{}}
	template.InitialLocation = NewLocation("L0", Invariant{"x", "<", 20})
	return template
}

func TestGuard_Evaluate(t *testing.T) {
	var localVariables map[string]int = SetupMap()
	var g Guard = Guard{"a", "<", 3}
	if g.Evaluate(localVariables) != true{
		t.Errorf("Evaluate function should return true, but returned false")
	}
}
func TestGuard_Evaluate2(t *testing.T) {
	var localVariables map[string]int = SetupMap()
	var g Guard = Guard{"b", "<", 3}
	if g.Evaluate(localVariables) == true{
		t.Errorf("Evaluate function should return false, but returned true")
	}
}
func TestUpdate_Update(t *testing.T) {
	var localVariables map[string]int = SetupMap()
	var u1 Update = Update{"x", "+=", 2}
	var u2 Update = Update{"x", "+", 2}
	u1.Update(localVariables)
	u2.Update(localVariables)
	if localVariables[u1.variableToUpdate] != localVariables[u2.variableToUpdate]{
		t.Errorf("'+' and '+=' do not return the same value, but should.")
	}
}

func TestTemplate_ToString(t *testing.T) {
	var template Template = SetupTemplate()
	var expected string = "027"
	var s string = template.ToString()

	if (s != expected) {
		t.Errorf("Got: %s --- expected: %s", s, expected)
	}
}

func TestHash(t *testing.T) {
	var a string = "aksdksd<jfnjikdfhvjikdvhfjdvh"
	var b string = "sdkfj<hvuifdhvuizhvfbvhfbvvbvbvbvbvbvbvbvbvbvbvbvbvbvbvbvb"

	c := Hash(a)
	d := Hash(b)
	if c == d{
		t.Errorf("Hash function returns same value of different strings, but should not: '%d' '%d'", c,d)
	}
}