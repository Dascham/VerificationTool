package main

import "hash/fnv"

const MaxValue = 128
const MinValue = -127

func main(){
	var s State = State{}
	var s1 State = State{}
	s.allTemplates = make([]Template, 2)

	var t Template = MainSetupCounterModel()

	s.allTemplates = append(s.allTemplates, t)
	s.globalVariables = MainSetupMap()

	s1 = s
	println(s.allTemplates[0].LocalVariables["x"])
	println(s1.allTemplates[0].LocalVariables["x"])

	var newMap map[string]int = make(map[string]int)
	for key, value := range s1.allTemplates[0].LocalVariables {
		newMap[key] = value
	}
	s1.allTemplates[0].InitialLocation.Edges[0].AtomicUpdate(newMap)
	//s1.allTemplates[0].LocalVariables = newMap

	println(s.allTemplates[0].LocalVariables["x"])
	println(s1.allTemplates[0].LocalVariables["x"])
	/*
	var mylist []State = make([]State, 2)
	mylist = append(mylist, s)

	//create new state
	s.allTemplates[0].currentLocation = s.allTemplates[0].InitialLocation
	s.allTemplates[0].currentLocation.Edges[0].Update
	*/
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

	var template Template
	template.LocalVariables = make(map[string]int)
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
