package main

import (
	"strconv"
	"strings"
)

type Guard struct {
	VariableToEvaluate string
	ComparisonOperator string
	GuardValue int
	GuardVar string
}
func (g Guard) Evaluate(localVariables map[string]int) bool {
	var result bool = false
	if(g.VariableToEvaluate == ""){
		return true
	}

	x, ok := localVariables[g.VariableToEvaluate] //lookup variable value, and assign to x
	y, ok1 := localVariables[g.GuardVar]
	//lookup right hand side variable, if in localvariables, we should use that, otherwise, we just use guardvalue
	if (!ok1){//if no GuardVar, then use guardvalue
		y = g.GuardValue
	}

	if (ok) {
		if (!ValidValue(x)) { //check on min max
			return false
		} else {
			switch g.ComparisonOperator {
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
			default:
				result = false
			}
		}
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

