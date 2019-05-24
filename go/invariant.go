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
	result := false
	if (!ValidValue(localVariables[i.VariableToEvaluate])){
		return false
	}

	if _, ok := localVariables[i.VariableToEvaluate]; ok || i.VariableToEvaluate=="" {
		switch i.ComparisonOperator {
		case "<":
			if (localVariables[i.VariableToEvaluate] < i.InvariantValue) {
				result = true
			}
		case ">":
			if (localVariables[i.VariableToEvaluate] > i.InvariantValue) {
				result = true
			}
		case "==":
			if (localVariables[i.VariableToEvaluate] == i.InvariantValue) {
				result = true
			}
		case "":
			result = true
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