package main

import (
	"fmt"
	"hash/fnv"
)

const MaxValue = 128
const MinValue = -127

//
/*----------------------------------------------- Useful functions ---------------------------------------------------*/
//

func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))

	return h.Sum32()
}
func ValidValue(a int) bool{
	if (MinValue < a && a < MaxValue){
		return true
	}else {
		return false
	}
}
func ValidMap(a map[string]int) bool{
	for _,value := range a{
		if (ValidValue(value)){

		} else {
			return false
		}
	}
	return true
}
func PrintStates(list []State) {
	for _, s := range list{
		fmt.Printf(s.ToString()+"\n")
	}
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
func remove(a []State, i int) []State {
	a[i] = a[len(a)-1] // Copy last element to index i.
	a[len(a)-1] = State{} // Erase last element (write zero value).
	a = a[:len(a)-1]   //truncate slice

	return a
}
func removeLocation(a []*Location, i int) []*Location {
	a[i] = a[len(a)-1] // Copy last element to index i.
	a[len(a)-1] = &Location{} // Erase last element (write zero value).
	a = a[:len(a)-1]   //truncate slice

	return a
}

func FromChannelToList(channel chan State)[]State{
	var tempList []State = make([]State,0,0)
	for { //for loop that consumes everything from the buffed channel 'channel' and appends to waiting list, which is returned when channel is empty
		select{
		case state := <- channel:
			tempList = append(tempList, state)
		default:
			return tempList
		}
	}
}
//
//
//
/*------------------------------- Functions that setup models to be model checked ------------------------------------*/
//
//
//
func SetupPotentiallyInfiniteModel() State{
	var update0 Update = Update{"x", "=", 1,""}
	var update1 Update = Update{"x", "=", 0,""}
	var location0 = NewLocation("L0", Invariant{})
	var location1 = NewLocation("L1", Invariant{})

	var edge0 Edge = Edge{}
	edge0 = edge0.InitializeEdge()
	edge0 = edge0.AcceptUpdates(update0)
	edge0 = edge0.AssignSrcDst(&location0, &location1)

	var edge1 Edge = Edge{}
	edge1 = edge1.InitializeEdge()
	edge1 = edge1.AcceptUpdates(update1)
	edge1 = edge1.AssignSrcDst(&location1, &location0)

	location0 = location0.AcceptOutGoingEdges(edge0)
	location1 = location1.AcceptOutGoingEdges(edge1)

	var template0 Template = Template{}
	template0.currentLocation = &location0
	template0.InitialLocation = &location0
	template0.LocalVariables = map[string]int{"x":0}

	var s State = State{}
	s.allTemplates = make([]Template, 0,0)
	s.allTemplates = append(s.allTemplates, template0)
	s.globalVariables = make(map[string]int)

	return s
}

func SetupFullModel() Template{
	//have global mutex, in order to change global state,
	// although not necessary for passed and waiting list implementation

	//define the things, guards
	var guard0 Guard = Guard{"x", ">", 5, ""}
	var guard1 Guard = Guard{"y", ">", 22, ""}
	var guard2 Guard = Guard{"x", "==", 0, ""}
	var guard3 Guard = Guard{"y", "==", 0, ""}

	var update0 Update = Update{"x", "++", 0, ""}
	var update1 Update = Update{"y", "=", 3,""}
	var update2 Update = Update{"y", "*=", 2,""}
	var update3 Update = Update{"x", "=", 4,""}
	var update4 Update = Update{"x", "--", 0,""}
	var update5 Update = Update{"y", "=", 17,""}
	var update6 Update = Update{"y", "--", 0,""}
	var update7 Update = Update{"x", "=", 0,""}
	var update8 Update = Update{"y", "=", 0,""}

	//default invariant
	var emptyInvariant Invariant = Invariant{}
	var invariant0 Invariant = Invariant{"x", "<", 10,""}

	var location0 Location = NewLocation("L0", invariant0)
	var location1 Location = NewLocation("L1", emptyInvariant)
	var location2 Location = NewLocation("L2", emptyInvariant)
	var location3 Location = NewLocation("L3", emptyInvariant)

	//edges
	var edge0 Edge = Edge{}
	edge0 = edge0.InitializeEdge()
	edge0 = edge0.AcceptUpdates(update0)
	edge0 = edge0.AssignSrcDst(&location0, &location0)
	var edge1 = Edge{}
	edge1 = edge1.InitializeEdge()
	edge1 = edge1.AcceptGuards(guard0)
	edge1 = edge1.AcceptUpdates(update1)
	edge1 = edge1.AssignSrcDst(&location0, &location1)
	var edge2 = Edge{}
	edge2 = edge2.InitializeEdge()
	edge2 = edge2.AcceptUpdates(update2)
	edge2 = edge2.AssignSrcDst(&location1, &location1)
	var edge3 = Edge{}
	edge3 = edge3.InitializeEdge()
	edge3 = edge3.AcceptGuards(guard1)
	edge3 = edge3.AcceptUpdates(update3)
	edge3 = edge3.AssignSrcDst(&location1, &location2)
	var edge4 = Edge{}
	edge4 = edge4.InitializeEdge()
	edge4 = edge4.AcceptUpdates(update4)
	edge4 = edge4.AssignSrcDst(&location2, &location2)
	var edge5 = Edge{}
	edge5 = edge5.InitializeEdge()
	edge5 = edge5.AcceptGuards(guard2)
	edge5 = edge5.AcceptUpdates(update5)
	edge5 = edge5.AssignSrcDst(&location2, &location3)
	var edge6 = Edge{}
	edge6 = edge6.InitializeEdge()
	edge6 = edge6.AcceptUpdates(update6)
	edge6 = edge6.AssignSrcDst(&location3, &location3)
	var edge7 = Edge{}
	edge7 = edge7.InitializeEdge()
	edge7 = edge7.AcceptGuards(guard3)
	edge7 = edge7.AcceptUpdates(update7, update8)
	//locations
	location0 = location0.AcceptOutGoingEdges(edge0, edge1)
	location1 = location1.AcceptOutGoingEdges(edge2, edge3)
	location2 = location2.AcceptOutGoingEdges(edge4, edge5)
	location3 =location3.AcceptOutGoingEdges(edge6, edge7)

	var template Template = Template{}
	template.InitialLocation = &location0

	return template
}
func SetupCounterModel() Template{
	var localVariables map[string]int = map[string]int{"x":0, "y":5,"z":10}
	var template Template = Template{}
	template.LocalVariables = localVariables
	var location0 Location = NewLocation("L0", Invariant{})
	template.InitialLocation = &location0
	template.currentLocation = &location0

	//update
	var update Update = Update{"x", "++", 0, ""}
	//edge
	var edge Edge = Edge{}
	edge = edge.InitializeEdge()
	edge = edge.AcceptUpdates(update)
	edge = edge.AssignSrcDst(&location0, &location0)
	location0 = location0.AcceptOutGoingEdges(edge)

	return template
}
func EmptyState()State{

	var s State = State{}
	s.globalVariables = make(map[string]int)
	s.allTemplates = make([]Template,0,0)
	return s
}

//a template with a single location, no edges etc.
func SetupTemplate() Template{
	var template Template = Template{SetupMap(), &Location{}, &Location{}, "tempname"}
	var location Location = NewLocation("L0", Invariant{"x", "<", 20, ""})
	template.InitialLocation = &location
	return template
}
func SetupSimpleSyncModel() State{
	var initialState State = State{}
	initialState.globalVariables = SetupMap3()

	//guards
	var guard0 Guard = Guard{"x", "==", 0,""}
	var guard1 Guard = Guard{"z", "<", 10,""}
	var update0 Update = Update{"y", "++", 0,""}
	var update1 Update = Update{"y", "*", 2,""}
	var invariant0 Invariant = Invariant{"z", "<", 10,""}
	var location0 Location = NewLocation("L0", Invariant{})
	var Location1 Location = NewLocation("L1", Invariant{})
	var location2 Location = NewLocation("L2", Invariant{})
	var location3 Location = NewLocation("L3", invariant0)


	var edge0 Edge = Edge{}
	edge0 = edge0.InitializeEdge()
	edge0 = edge0.AcceptGuards(guard0)
	edge0 = edge0.AcceptUpdates(update0)
	edge0 = edge0.AssignSrcDst(&location0, &Location1)
	edge0.Ch = "a"
	edge0.IsSend = true

	var edge1 Edge = Edge{}
	edge1 = edge1.InitializeEdge()
	edge1 = edge1.AcceptGuards(guard1)
	edge1 = edge1.AcceptUpdates(update1)
	edge1 = edge1.AssignSrcDst(&location2, &location3)
	edge1.Ch = "a"
	edge1.IsSend = false

	location0 = location0.AcceptOutGoingEdges(edge0)
	location2 = location2.AcceptOutGoingEdges(edge1)

	//template0
	var template0 Template = Template{}
	template0.LocalVariables = map[string]int{"x":0}
	template0.InitialLocation = &location0
	template0.currentLocation = &location0

	var template1 Template = Template{}
	template1.LocalVariables = make(map[string]int) //let be empty
	template1.InitialLocation = &location2
	template1.currentLocation = &location2

	initialState.allTemplates = append(initialState.allTemplates, template0, template1)

	return initialState
}
func MainSetupCounterModel() Template{
	var localVariables map[string]int = map[string]int{"x":0}
	var template Template = Template{}
	template.LocalVariables = localVariables
	var location0 Location = NewLocation("L0", Invariant{})
	template.InitialLocation = &location0
	template.currentLocation = &location0

	//update
	var update Update = Update{"x", "++", 0,""}
	//edge
	var edge Edge = Edge{}
	edge = edge.InitializeEdge()
	edge = edge.AcceptUpdates(update)
	edge = edge.AssignSrcDst(&location0, &location0)

	location0 = location0.AcceptOutGoingEdges(edge)
	return template
}
//v_i == a,
func SetupTestModel1()Template{
	var guard0 Guard = Guard{"a", ">", 2, ""}
	var guard1 Guard = Guard{"a", "<", 16, ""}
	var guard2 Guard = Guard{"a", ">", 5, ""}

	var update0 Update = Update{"a", "+=", 3, ""}
	var update1 Update = Update{"a", "+=", 2, ""}
	var update2 Update = Update{"a", "--", 1, ""}
	var update3 Update = Update{"a", "++", 1, ""}
	var update4 Update = Update{"a", "-=", 3, ""}

	var location0 Location = NewLocation("L0", Invariant{})
	location0.BlockId = 0
	var location1 Location = NewLocation("L1", Invariant{})
	location1.BlockId = 1
	var location2 Location = NewLocation("L2", Invariant{})
	location2.BlockId = 2
	var location3 Location = NewLocation("L3", Invariant{})
	location3.BlockId = 3
	var location4 Location = NewLocation("L4", Invariant{})
	location4.BlockId = 4
	var location5 Location = NewLocation("L5", Invariant{})
	location5.BlockId = 4
	var location6 Location = NewLocation("L6", Invariant{})
	location6.BlockId = 5

	//0,1
	var edge01 Edge = Edge{}
	edge01 = edge01.InitializeEdge()
	edge01 = edge01.AcceptGuards(guard0)
	edge01 = edge01.AssignSrcDst(&location0, &location1)

	//0,3
	var edge03 Edge = Edge{}
	edge03 = edge03.InitializeEdge()
	edge03 = edge03.AcceptGuards(guard1)
	edge03 = edge03.AssignSrcDst(&location0, &location3)

	//1,2
	var edge12 Edge = Edge{}
	edge12 = edge12.InitializeEdge()
	edge12 = edge12.AcceptUpdates(update2)
	edge12 = edge12.AssignSrcDst(&location1, &location2)

	//1,3
	var edge13 Edge = Edge{}
	edge13 = edge13.InitializeEdge()
	edge13 = edge13.AcceptGuards(guard1)
	edge13 = edge13.AcceptUpdates(update0)
	edge13 = edge13.AssignSrcDst(&location1, &location3)

	//2,4
	var edge24 Edge = Edge{}
	edge24 = edge24.InitializeEdge()
	edge24 = edge24.AssignSrcDst(&location2, &location4)

	//3,4
	var edge34 Edge = Edge{}
	edge34 = edge34.InitializeEdge()
	edge34 = edge34.AssignSrcDst(&location3, &location4)

	//3,5
	var edge35 Edge = Edge{}
	edge35 = edge35.InitializeEdge()
	edge35 = edge35.AcceptUpdates(update1)
	edge35 = edge35.AssignSrcDst(&location3, &location5)
	//4,5
	var edge45 Edge = Edge{}
	edge45 = edge45.InitializeEdge()
	edge45 = edge45.AssignSrcDst(&location4, &location5)

	//5,6
	var edge56 Edge = Edge{}
	edge56 = edge56.InitializeEdge()
	edge56 = edge56.AcceptGuards(guard1)
	edge56 = edge56.AssignSrcDst(&location5, &location6)

	//5,0
	var edge50 Edge = Edge{}
	edge50 = edge50.InitializeEdge()
	edge50 = edge50.AcceptGuards(guard2)
	edge50 = edge50.AcceptUpdates(update4)
	edge50 = edge50.AssignSrcDst(&location5, &location0)

	//6,0
	var edge60 Edge = Edge{}
	edge60 = edge60.InitializeEdge()
	edge60 = edge60.AcceptUpdates(update3)
	edge60 = edge60.AssignSrcDst(&location6, &location0)

	location0 = location0.AcceptOutGoingEdges(edge01, edge03)
	location1 = location1.AcceptOutGoingEdges(edge12, edge13)
	location2 = location2.AcceptOutGoingEdges(edge24)
	location3 = location3.AcceptOutGoingEdges(edge34, edge35)
	location4 = location4.AcceptOutGoingEdges(edge45)
	location5 = location5.AcceptOutGoingEdges(edge56, edge50)
	location6 = location6.AcceptOutGoingEdges(edge60)

	var t Template = Template{}
	t.InitialLocation = &location0
	t.currentLocation = &location0
	t.LocalVariables = map[string]int{"a":0}
	return t
}

//
//
//
/*------------------------------------------ Setup of maps used in the models ----------------------------------------*/
//
//
//
func SetupMap() map[string]int {
	var localVariables map[string]int = make(map[string]int)
	localVariables["a"] = 2
	localVariables["b"] = 7
	return localVariables
}
func SetupMap1() map[string]int {
	var localVariables map[string]int = make(map[string]int)
	localVariables["c"] = 6
	localVariables["d"] = 14
	return localVariables
}
func SetupMap2() map[string]int {
	var localVariables map[string]int = make(map[string]int)
	localVariables["y"] = 9
	localVariables["z"] = 21
	return localVariables
}
func SetupMap3()map[string]int{
	var localVariables map[string]int = map[string]int{"y":8,"z":5}
	return localVariables
}
func SetupInvalidMap() map[string]int {
	var localVariables map[string]int = make(map[string]int)
	localVariables["c"] = -129
	localVariables["d"] = 135
	return localVariables
}


