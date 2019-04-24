package main

type Guard struct {
	GuardValue int
	ComparisonOperator string
	VariableToEvaluate string
}
//a function that only applies to guards
func (g Guard) Evaluate(localVariables map[string]int) bool{
	var result bool = false
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
	default:
		result = false
	}
	return result
}
