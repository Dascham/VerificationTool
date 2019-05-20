package main

import (
	"strconv"
	"strings"
)

type State struct {
	allTemplates []Template
	globalVariables map[string]int

}
//this to string produces the following format:
//(TemplateID ValueofLocalvariables)* ValueofGlobalVariables
func (s State) ToString() string{
	var sb strings.Builder
	for i := 0; i < len(s.allTemplates); i++{
		sb.WriteString(s.allTemplates[i].ToString())
	}
	for _, value := range s.globalVariables{
		sb.WriteString(strconv.Itoa(value))
	}
	return sb.String()
}