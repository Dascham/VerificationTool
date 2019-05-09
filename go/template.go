package main

import (
	"strconv"
	"strings"
)

type Template struct {
	LocalVariables  map[string]int
	InitialLocation *Location
	currentLocation *Location
}

//tostring, such that we can hash a template, which should be done only once,
//per reachability check, since the template doesn't change during model checking
func (t Template) ToString()string{
	var sb strings.Builder

	//location id, which is unique -> good
	sb.WriteString((strconv.Itoa(t.InitialLocation.LocationId)))

	//then all variables, which is basically the localvariables map
	for _, value := range t.LocalVariables {
		sb.WriteString(strconv.Itoa(value))
	}
	return sb.String()
}

//function that should enumerate all locations, done by following pointers
func (t Template) PrettyPrint(){

}
