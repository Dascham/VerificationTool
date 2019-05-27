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
	x, ok := newMap[u.variableToUpdate]
	y, ok1 := newMap[u.updateVar]
	if(!ok1) { //if updateVar == "", i.e. not in map, then we just use updatevalue
		y = u.updateValue
	}

	if(ok) {
		switch u.operator {
		case "+", "+=":
			newMap[u.variableToUpdate] = x + y
		case "-", "-=":
			newMap[u.variableToUpdate] = x - y
		case "*", "*=":
			newMap[u.variableToUpdate] = x * y
		case "/", "/=":
			newMap[u.variableToUpdate] = x / y
		case "=":
			newMap[u.variableToUpdate] = y
		case "++":
			newMap[u.variableToUpdate] = x + 1
		case "--":
			newMap[u.variableToUpdate] = x - 1
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
