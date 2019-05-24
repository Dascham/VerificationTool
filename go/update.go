package main

import (
	"strconv"
	"strings"
)

type Update struct {
	variableToUpdate string
	operator string
	updateValue int
	updateVar string
}
func (u Update) Update(newMap map[string]int) {
	if _, ok := newMap[u.variableToUpdate]; ok{
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
}

func (u Update) ToString()string{
	var sb strings.Builder
	sb.WriteString(u.variableToUpdate)
	sb.WriteString(u.operator)
	sb.WriteString(strconv.Itoa(u.updateValue))
	sb.WriteString("\n")
	return sb.String()
}
