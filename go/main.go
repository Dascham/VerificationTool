package main

import (
	"hash/fnv"
)
const MaxValue = 128
const MinValue = -127

func main(){
	var template Template = SetupCounterModel()
	var state State = State{}

	for i:=0;i<2;i++{
		state.allTemplates = append(state.allTemplates, template)
	}
	state.globalVariables = SetupMap()
	print(len(state.allTemplates))

}
func remove(a []State, i int) []State {
	a[i] = a[len(a)-1] // Copy last element to index i.
	a[len(a)-1] = State{}   // Erase last element (write zero value).
	a = a[:len(a)-1]   //truncate slice

	return a
}

func Explore(initialState State) []State{
	var waitingList []State = make([]State, 10000) //dunno how big should be, but many states
	waitingList = append(waitingList, initialState)
	var passedList []State = make([]State, 10000)

	for len(waitingList) > 0 {
		var currentState = waitingList[0]
		for template := range len(currentState.allTemplates){

		}
	}


	return passedList
}


func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func SetupCounterModel() Template{
	var localVariables map[string]int = map[string]int{"x":0}
	var template Template
	template.LocalVariables = localVariables
	var location0 Location = NewLocation("L0", Invariant{})
	template.InitialLocation = &location0

	//update
	var update Update = Update{"x", "++", 0}
	//edge
	var edge Edge = Edge{}
	edge.InitializeEdge()
	edge.AcceptUpdates(update)
	edge.AssignSrcDst(location0, location0)

	location0.AcceptOutGoingEdges(edge)

	return template
}

func SetupMap() map[string]int {
	var localVariables map[string]int = make(map[string]int)
	localVariables["a"] = 2
	localVariables["b"] = 7
	return localVariables
}