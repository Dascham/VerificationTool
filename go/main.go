package main

import (
	"hash/fnv"
)

const MaxValue = 128
const MinValue = -127

func main(){


}
func remove(a []State, i int) []State {
	a[i] = a[len(a)-1] // Copy last element to index i.
	a[len(a)-1] = State{}   // Erase last element (write zero value).
	a = a[:len(a)-1]   //truncate slice

	return a
}

/*
func Explore(initialState State) []State{
	var waitingList []State = make([]State, 10000) //dunno how big should be, but many states
	waitingList = append(waitingList, initialState)
	var passedList []State = make([]State, 10000)

	for len(waitingList) > 0 { //exploration loop
		var currentState = waitingList[0]
		for i := 0; i < len(currentState.allTemplates);i++ {
			currentState.allTemplates[i].InitialLocation.
		}
	}


	return passedList
}

 */

func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func MainSetupCounterModel() Template{
	var localVariables map[string]int = make(map[string]int)
	localVariables["x"] = 0

	var template Template = Template{}
	template.tempName = "name"
	//template.LocalVariables = make(map[string]int)
	template.LocalVariables = localVariables

	var location0 Location = NewLocation("L0", Invariant{})
	template.InitialLocation = &location0
	template.currentLocation = &location0

	//update
	var update Update = Update{"x", "++", 0}
	//edge
	var edge Edge = Edge{}
	edge = edge.InitializeEdge()
	edge = edge.AcceptUpdates(update)
	edge = edge.AssignSrcDst(location0, location0)
	edge.name = "edge name"

	location0 = location0.AcceptOutGoingEdges(edge)
	return template
}

func MainSetupMap() map[string]int {
	var localVariables map[string]int = make(map[string]int)
	localVariables["a"] = 2
	localVariables["b"] = 7
	return localVariables
}

func CopyMap(originalMap map[string]int) map[string]int{
	var newMap map[string]int = make(map[string]int)
	for key, value := range originalMap {
		newMap[key] = value
	}
	return newMap
}

func DeepCopyState(s State) State{
	var newState State = State{}
	newState.allTemplates = make([]Template, 0,0)

	//copy templates
	for i := 0; i<len(s.allTemplates);i++{
		newState.allTemplates = append(newState.allTemplates, s.allTemplates[i])
	}

	newState.globalVariables = CopyMap(s.globalVariables)

	for i := 0; i<len(s.allTemplates);i++ {
		newState.allTemplates[i].LocalVariables = CopyMap(s.allTemplates[i].LocalVariables)
	}

	return newState
}