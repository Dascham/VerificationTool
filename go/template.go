package main

import (
	"sort"
	"strconv"
	"strings"
)

type Template struct {
	LocalVariables  map[string]int
	InitialLocation *Location
	currentLocation *Location
	tempName string
}

//tostring, such that we can hash a template, which should be done only once,
//per reachability check, since the template doesn't change during model checking
func (t Template) ToString()string{
	var sb strings.Builder

	//location id, which is unique -> good
	sb.WriteString("Location: "+(strconv.Itoa(t.currentLocation.LocationId)+" "))

	//then all variables, which is basically the localvariables map
	var keys []string
	for key, _ := range t.LocalVariables{
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for i:=0;i<len(keys);i++{
		sb.WriteString(keys[i]+":"+strconv.Itoa(t.LocalVariables[keys[i]])+" ")
	}

	return sb.String()
}

//function that should enumerate all locations, done by following pointers
func (t Template) PrettyPrint(){

}
