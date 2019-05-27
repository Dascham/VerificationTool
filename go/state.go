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
		s.allTemplates[i].currentLocation = helperConfigLocation(s.allTemplates[i], si.CurrentLocationIds[i])
	}

	return s
}
func helperConfigLocation(t Template, locationid int) *Location{
	var hashTable map[int]bool = make(map[int]bool)

	var correctLocation *Location
	var locations []*Location = make([]*Location, 0,0)
	locations = append(locations, t.InitialLocation)
	hashTable[t.InitialLocation.LocationId] = true
	for len(locations) > 0 {
		if (locations[0].LocationId == locationid){
			correctLocation = locations[0]
			break;
		}else{
			//add all locations to list, from outgoing edges
			for _,edge := range locations[0].Edges{
				_, ok := hashTable[edge.Dst.LocationId] //ok is true, if locationId is already in hashtable
				if !ok { //if not ok, then we have not seen the dst before and we should add it
					locations = append(locations, edge.Dst)
					hashTable[edge.Dst.LocationId] = true
				}
			}
			removeLocation(locations, 0)
		}
	}
	return correctLocation
}

