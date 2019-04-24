package main

type Update struct {
	updateValue int
	variableToUpdate string
	operator string
}
func (u Update) Update(localVariables map[string]int){
	switch u.operator {
	case "+":
		localVariables[u.variableToUpdate] = localVariables[u.variableToUpdate] + u.updateValue
	case "-":
		localVariables[u.variableToUpdate] = localVariables[u.variableToUpdate] - u.updateValue
	case "*":
		localVariables[u.variableToUpdate] = localVariables[u.variableToUpdate] * u.updateValue
	case "/":
		localVariables[u.variableToUpdate] = localVariables[u.variableToUpdate] / u.updateValue
	case "":
		localVariables[u.variableToUpdate] = u.updateValue
	}
}
