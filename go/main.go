package main

import (
	"hash/fnv"
)

const MaxValue = 128
const MinValue = -127

func main(){
	var initialState State = State{}
	initialState.allTemplates = make([]Template, 0,0)
	initialState.allTemplates = append(initialState.allTemplates, MainSetupCounterModel())

	var list []State = Explore(initialState)

	for i := 0; i < len(list);i++{
		println(list[i].ToString())
	}
}


func remove(a []State, i int) []State {
	a[i] = a[len(a)-1] // Copy last element to index i.
	a[len(a)-1] = State{} // Erase last element (write zero value).
	a = a[:len(a)-1]   //truncate slice

	return a
}

func Explore(initialState State) []State{
	var waitingList []State = make([]State, 0,0) //size zero, always, cause append fixes size by itself
	waitingList = append(waitingList, initialState)

	var passedList []State = make([]State, 0,0)

	for len(waitingList) > 0 { //exploration loop
		println("First print")
		var currentState = waitingList[0]
		//remove the element
		waitingList = remove(waitingList, 0)
		for i := 0;i < len(currentState.allTemplates);i++{
			println("second for loop?")
			println(len(currentState.allTemplates[i].currentLocation.Edges))
			for j := 0 ; j < len(currentState.allTemplates[i].currentLocation.Edges); j++{
				println("Third for loop?")
				println(currentState.allTemplates[i].currentLocation.Edges[j].EdgeIsActive(currentState.allTemplates[i].LocalVariables))
				if (currentState.allTemplates[i].currentLocation.Edges[j].EdgeIsActive(currentState.allTemplates[i].LocalVariables)){
					//have to instantiate new state here
					newState := DeepCopyState(currentState)
					println("We here now")
					//do update on new state
					newState.allTemplates[i].currentLocation.Edges[j].AtomicUpdate(newState.allTemplates[i].LocalVariables)
					//then advance location, by looking to dst of edge that we took
					newState.allTemplates[i].currentLocation = newState.allTemplates[i].currentLocation.Edges[j].Dst
					//add newstate to waitinglist, for distributed, call distribute function, which hashes and does stuff
					//add only if map is valid, this should be made to better fix
					if (ValidMap(newState.allTemplates[i].LocalVariables)) {
						println("we here?")
						waitingList = append(waitingList, newState)
					}
				}
			}
		}
		//now we have tried everything that is possible for this state, therefore add to passed list
		passedList = append(passedList, currentState)
	}
	return passedList
}

//helper function, for finding an synchronization partner for an edge
func FindSync(e Edge) (bool, Edge){

	return false, Edge{}
}

func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func MainSetupCounterModel() Template{
	var localVariables map[string]int = map[string]int{"x":0}
	var template Template = Template{}
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
	edge = edge.AssignSrcDst(&location0, &location0)

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
	//newState.allTemplates = s.allTemplates

	//old way
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

func UpdateLocation(t Template, l *Location)Template{
	t.currentLocation = l

	return t
}