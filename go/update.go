package main

type Update struct {
	variableToUpdate string
	operator string
	updateValue int
}
func (u Update) Update(newMap map[string]int) {
	//var newMap = CopyMap(localVariables)
	switch u.operator {
	case "+", "+=":
		newMap[u.variableToUpdate] = newMap[u.variableToUpdate] + u.updateValue
	case "-", "-=":
		newMap[u.variableToUpdate] = newMap[u.variableToUpdate] - u.updateValue
	case "*", "*=":
		newMap[u.variableToUpdate] = newMap[u.variableToUpdate] * u.updateValue
	case "/", "/=":
		newMap[u.variableToUpdate] = newMap[u.variableToUpdate] / u.updateValue
	case "=":
		newMap[u.variableToUpdate] = u.updateValue
	case "++":
		newMap[u.variableToUpdate] = newMap[u.variableToUpdate] + 1
	case "--":
		newMap[u.variableToUpdate] = newMap[u.variableToUpdate] - 1
	case "":
		break
	}
}
