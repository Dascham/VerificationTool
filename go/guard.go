package main

import (
	"strconv"
	"strings"
)

type Guard struct {
	VariableToEvaluate string
	ComparisonOperator string
	GuardValue int
}
//a function that only applies to guards
func (g Guard) Evaluate(localVariables map[string]int) bool{
	var result bool = false

	if (MinValue < localVariables[g.VariableToEvaluate] && localVariables[g.VariableToEvaluate] < MaxValue){

	} else {
		return false
	}

	switch g.ComparisonOperator {
	case "<":
		if(localVariables[g.VariableToEvaluate] < g.GuardValue){
			result = true
		}
	case ">":
		if(localVariables[g.VariableToEvaluate] > g.GuardValue){
			result = true
		}
	case "==":
		if(localVariables[g.VariableToEvaluate] == g.GuardValue){
			result = true
		}
	case "":
		result = true
	default:
		result = false
	}
	return result
}
func (g Guard) ToString()string{
	var sb strings.Builder
	sb.WriteString(g.VariableToEvaluate)
	sb.WriteString(g.ComparisonOperator)
	sb.WriteString(strconv.Itoa(g.GuardValue))
	sb.WriteString("\n")
	return sb.String()
}
