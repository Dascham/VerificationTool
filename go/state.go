package main

import (
	"sort"
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

	var keys []string
	for key, _ := range s.globalVariables{
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for i:=0;i<len(keys);i++{
		sb.WriteString(keys[i]+":"+strconv.Itoa(s.globalVariables[keys[i]])+" ")
	}
	return sb.String()
}