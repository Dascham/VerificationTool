package main

type Invariant struct {
	InvariantValue int
	ComparisonOperator string
	VariableToEvaluate string
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
	}
	return result
}
