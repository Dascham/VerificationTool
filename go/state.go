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

func (s State) ConfigureState(si StateInformation) State{
	s.globalVariables = CopyMap(si.GlobalVariables)
	for i, tempMap := range si.ListLocalVariables{
		s.allTemplates[i].LocalVariables = CopyMap(tempMap)
		s.allTemplates[i].currentLocation = configLocation(s.allTemplates[i], si.CurrentLocationIds[i])
	}

	return s
}
func configLocation(t Template, locationid int) *Location{
	var correctLocation *Location
	var locations []*Location = make([]*Location, 0,0)
	locations = append(locations, t.InitialLocation)
	for len(locations) > 0 {
		if (locations[0].LocationId == locationid){
			correctLocation = locations[0]
			break;
		}else{
			//add all locations to list, from outgoing edges
			for _,edge := range locations[0].Edges{
				locations = append(locations, edge.Dst)
			}
			removeLocation(locations, 0)
		}
	}
	return correctLocation
}
