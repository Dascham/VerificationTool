package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Invariant struct {
	VariableToEvaluate string
	ComparisonOperator string
	InvariantValue int
	InvariantVar string
}
func (i Invariant) iPrintln() {
	fmt.Printf("'%s':'%s':'%s'", i.VariableToEvaluate, i.ComparisonOperator, strconv.Itoa(i.InvariantValue))
}

func (i Invariant) IsValid(localVariables map[string]int) bool {
	var result bool = false
	if (i.VariableToEvaluate == ""){ //if there ain't no variable to evaluate, must be empty invariant, and that's fine
		return true
	}
	x, ok := localVariables[i.VariableToEvaluate]
	y, ok1 := localVariables[i.InvariantVar]
	if(!ok1){
		y = i.InvariantValue
	}

	if (ok) {
		if(!ValidValue(x)){
			return false
		}else {
			switch i.ComparisonOperator {
			case "<":
				if (x < y) {
					result = true
				}
			case ">":
				if (x > y) {
					result = true
				}
			case "==":
				if (x == y) {
					result = true
				}
			case "":
				result = true
			}
		}
	}
	return result
}

func (i Invariant) ToString()string{
	var sb strings.Builder
	sb.WriteString(i.VariableToEvaluate)
	sb.WriteString(i.ComparisonOperator)
	sb.WriteString(strconv.Itoa(i.InvariantValue))
	sb.WriteString("\n")
	return sb.String()
}