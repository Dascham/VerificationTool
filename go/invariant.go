package main

import (
	"fmt"
	"strconv"
)

type Invariant struct {
	VariableToEvaluate string
	ComparisonOperator string
	InvariantValue int
}
func (i Invariant) iPrintln() {
	fmt.Printf("'%s':'%s':'%s'", i.VariableToEvaluate, i.ComparisonOperator, strconv.Itoa(i.InvariantValue))
}

func (i Invariant) IsValid(localVariables map[string]int) bool {
	result := false
	switch i.ComparisonOperator {
	case "<":
		if (i.InvariantValue < localVariables[i.VariableToEvaluate]) {
			result = true
		}
	case ">":
		if (i.InvariantValue > localVariables[i.VariableToEvaluate]){
			result = true
		}
	case "==":
		if (i.InvariantValue == localVariables[i.VariableToEvaluate]){
			result = true
		}
	case "":
		result = true
	}
	return result
}
