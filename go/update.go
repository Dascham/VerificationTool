package main

type Update struct {
	variableToUpdate string
	operator string
	updateValue int
}
func (u Update) Update(localVariables map[string]int){
	switch u.operator {
	case "+", "+=":
		localVariables[u.variableToUpdate] = localVariables[u.variableToUpdate] + u.updateValue
	case "-", "-=":
		localVariables[u.variableToUpdate] = localVariables[u.variableToUpdate] - u.updateValue
	case "*", "*=":
		localVariables[u.variableToUpdate] = localVariables[u.variableToUpdate] * u.updateValue
	case "/", "/=":
		localVariables[u.variableToUpdate] = localVariables[u.variableToUpdate] / u.updateValue
	case "=":
		localVariables[u.variableToUpdate] = u.updateValue
	case "":
		break
	}
}
