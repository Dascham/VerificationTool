package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func SetupMap() map[string]int {
	var localVariables map[string]int = make(map[string]int)
	localVariables["a"] = 2
	localVariables["b"] = 7
	//localVariables["x"] = 0
	return localVariables
}
func SetupInvalidMap() map[string]int {
	var localVariables map[string]int = make(map[string]int)
	localVariables["c"] = -129
	localVariables["d"] = 135
	return localVariables
}

//a template with a single location, no edges etc.
func SetupTemplate() Template{
	var template Template = Template{SetupMap(), &Location{}, &Location{}, "tempname"}
	var location Location = NewLocation("L0", Invariant{"x", "<", 20})
	template.InitialLocation = &location
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
	var update Update = Update{"x", "++", 0}
	//edge
	var edge Edge = Edge{}
	edge = edge.InitializeEdge()
	edge = edge.AcceptUpdates(update)
	edge = edge.AssignSrcDst(&location0, &location0)
	location0 = location0.AcceptOutGoingEdges(edge)

	return template
}



//this test should be the first test to call the function SetupTemplate(), otherwise it will fail

func TestTemplate_ToString(t *testing.T) {
	var template Template = SetupTemplate()
	var expected string = "Location: 0 a:2 b:7 "
	var s string = template.ToString()


	if (s != expected) {
		t.Errorf("Got: %s --- expected: %s", s, expected)
	}

}




func SetupFullModel() Template{
	//have global mutex, in order to change global state,
	// although not necessary for passed and waiting list implementation

	//define the things, guards
	var guard0 Guard = Guard{"x", ">", 5}
	var guard1 Guard = Guard{"y", ">", 22}
	var guard2 Guard = Guard{"x", "==", 0}
	var guard3 Guard = Guard{"y", "==", 0}

	var update0 Update = Update{"x", "++", 0}
	var update1 Update = Update{"y", "=", 3}
	var update2 Update = Update{"y", "*=", 2}
	var update3 Update = Update{"x", "=", 4}
	var update4 Update = Update{"x", "--", 0}
	var update5 Update = Update{"y", "=", 17}
	var update6 Update = Update{"y", "--", 0}
	var update7 Update = Update{"x", "=", 0}
	var update8 Update = Update{"y", "=", 0}

	//default invariant
	var emptyInvariant Invariant = Invariant{}
	var invariant0 Invariant = Invariant{"x", "<", 10}

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


func TestGuard_Evaluate(t *testing.T) {
	var localVariables map[string]int = SetupMap()
	var g Guard = Guard{"a", "<", 3}
	if g.Evaluate(localVariables) != true{
		t.Errorf("Evaluate function should return true, but returned false")
	}
}

func TestGuard_Evaluate2(t *testing.T) {
	var localVariables map[string]int = SetupMap()
	var g Guard = Guard{"b", "<", 3}
	if g.Evaluate(localVariables) == true{
		t.Errorf("Evaluate function should return false, but returned true")
	}
}
func TestUpdate_Update(t *testing.T) {
	var localVariables map[string]int = SetupMap()
	var u1 Update = Update{"x", "+=", 2}
	var u2 Update = Update{"x", "+", 2}
	u1.Update(localVariables)
	u2.Update(localVariables)
	if localVariables[u1.variableToUpdate] != localVariables[u2.variableToUpdate]{
		t.Errorf("'+' and '+=' do not return the same value, but should.")
	}
}



func TestUpdate_Update2(t *testing.T) {
	var s State = State{}
	s.allTemplates = make([]Template, 0, 0)

	var template1 Template = SetupCounterModel()

	s.allTemplates = append(s.allTemplates, template1)
	s.globalVariables = SetupMap()

	s1 := DeepCopyState(s)
	//println(s.allTemplates[0].LocalVariables["x"])
	//println(s1.allTemplates[0].LocalVariables["x"])

	//s1.allTemplates[0].InitialLocation.Edges[0] = s1.allTemplates[0].InitialLocation.Edges[0].AtomicUpdate(s1.allTemplates[0].LocalVariables)
	s1.allTemplates[0].InitialLocation.Edges[0].AtomicUpdate(s1.allTemplates[0].LocalVariables, s1.globalVariables)
	//println(s.allTemplates[0].LocalVariables["x"])
	//println(s1.allTemplates[0].LocalVariables["x"])

	if s.allTemplates[0].LocalVariables["x"] == s1.allTemplates[0].LocalVariables["x"]{
		t.Errorf("AtomicUpdate does not work: %d == %d",
			s.allTemplates[0].LocalVariables["x"], s1.allTemplates[0].LocalVariables["x"])
	}
}

func TestInvariant_IsValid(t *testing.T) {
	//two maps
	var map1 map[string]int = map[string]int{"f":5}
	//var map2 map[]
	var i Invariant = Invariant{"f", "<", 10}

	if (i.IsValid(map1)){

	}else{
		t.Errorf("Expected to be valid, but was not")
	}

}

func TestLocation_AcceptOutGoingEdges(t *testing.T) {
	var l0 Location = NewLocation("L0", Invariant{})
	var l1 Location = NewLocation("L1", Invariant{})
	l0.Edges = append(l0.Edges, Edge{})
	l0.Edges = append(l0.Edges, Edge{})
	l0.Edges = append(l0.Edges, Edge{})

	l1 = l1.AcceptOutGoingEdges(Edge{}, Edge{}, Edge{})

	if len(l0.Edges) != len(l1.Edges) {
		t.Errorf("function accept outgoing edges does not work")
	}

	//println("Len 0: ",len(l0.Edges))
	//println("Len 1: ", len(l1.Edges))
}

func TestCopyMap(t *testing.T) {
	var map1 map[string]int = make(map[string]int)

	map1["x"] = 2
	map2 := CopyMap(map1)
	map2["x"] = 5
	if map1["x"] == map2["x"]{
		t.Errorf("copymap does not work")
	}
}

func TestDeepCopyState(t *testing.T) {
	var temp Template = SetupCounterModel()
	var s State = State{}
	s.allTemplates = make([]Template, 0,0)
	s.allTemplates = append(s.allTemplates, temp)
	s.globalVariables = SetupMap()

	copy_s := DeepCopyState(s)
	copy_s.globalVariables["a"] = 17
	if copy_s.globalVariables["a"] == s.globalVariables["a"]{
		t.Errorf("deepcopystate does not work, prolly a reference thing")
	}
}
func TestDeepCopyState2(t *testing.T) {
	var temp Template = SetupCounterModel()
	var s State = State{}
	s.allTemplates = make([]Template, 0,0)
	s.allTemplates = append(s.allTemplates, temp)
	s.globalVariables = SetupMap()
	//copy state
	copy_s := DeepCopyState(s)
	//update
	copy_s.allTemplates[0].currentLocation.Edges[0].AtomicUpdate(copy_s.allTemplates[0].LocalVariables, copy_s.globalVariables)
	//advance location
	copy_s.allTemplates[0].currentLocation = copy_s.allTemplates[0].currentLocation.Edges[0].Dst

	//copy state
	copy_s1 := DeepCopyState(copy_s)
	//update
	copy_s1.allTemplates[0].currentLocation.Edges[0].AtomicUpdate(copy_s1.allTemplates[0].LocalVariables, copy_s1.globalVariables)
	//advance location
	copy_s1.allTemplates[0].currentLocation = copy_s1.allTemplates[0].currentLocation.Edges[0].Dst

	//shoud be some other condition, but fine for now
	if (s.ToString() == copy_s.ToString() || s.ToString() == copy_s1.ToString() || copy_s.ToString() == copy_s1.ToString()){
		t.Errorf("Expected increaments between the states, but the variables are the same, so update probably failed")
	}

}

func TestHash(t *testing.T) {
	var a string = "aksdksd<jfnjikdfhvjikdvhfjdvh"
	var b string = "sdkfj<hvuifdhvuizhvfbvhfbvvbvbvbvbvbvbvbvbvbvbvbvbvbvbvbvb"

	c := Hash(a)
	d := Hash(b)
	if c == d{
		t.Errorf("Hash function returns same value of different strings, but should not: '%d' '%d'", c,d)
	}
}

func TestValidMap1(t *testing.T) {
	if ValidMap(SetupMap()){

	} else {
		t.Errorf("SetupMap returned a map that is invalid")
	}
}
func TestValidMap2(t *testing.T) {
	if ValidMap(SetupInvalidMap()){
		t.Errorf("Setup invalid map evaluated to true, which it should not")
	} else {

	}
}

func TestEdge_EdgeIsActive(t *testing.T) {
	map1 := SetupMap()
	template1 := SetupCounterModel()
	template1.LocalVariables = map1
	//println(template1.currentLocation.Edges[0].ToString())
	if (template1.currentLocation.Edges[0].EdgeIsActive(template1.LocalVariables, State{})){

	} else{
		t.Errorf("Function evaluated to false, but expected true")
	}
}

func TestHashedStates(t *testing.T){
	s := State{}
	s.globalVariables = SetupMap()
	s.allTemplates = make([]Template, 0,0)
	s.allTemplates = append(s.allTemplates, SetupCounterModel())
	s.allTemplates = append(s.allTemplates, SetupCounterModel())
	hashedStates := make(map[string]string)

	//put in state
	hashedStates[s.ToString()] = s.ToString()
	//this test mail fail because of concurrency, because values might be switched

	if _, ok := hashedStates[s.ToString()]; ok{
		//println("We've seen before")
	}else {
		t.Errorf("We expect to see the state in the hashtable, since we just put it there. This " +
			"test may have failed, because concurrency.")
	}
}
func TestUpdate_Update3(t *testing.T) {
	template1 := SetupCounterModel()
	template1.LocalVariables = SetupMap()

	temp := template1.ToString()
	//empty map for global variables, which is right
	template1.currentLocation.Edges[0].AtomicUpdate(template1.LocalVariables, map[string]int{})

	if (temp != template1.ToString()){
		t.Errorf("Expected the ToString methods to return the same results, even though we update, because we update 'x', which " +
			"isn't in the map, and therefore no variables should be changed")
	}
	//if fail, prolly because 'x' has been added to localvariables

}

func TestExplore(t *testing.T) {
	var initialState State = SetupSimpleSyncModel()
	var list []State = Explore(initialState)

	if len(list) != 2{
		t.Errorf("Expected to find, in total, 2 states for setupsimplesyncmodel")

	}
	if(list[1].globalVariables["y"] != 18){
		t.Errorf("Variable should be '18', but is %d", list[1].globalVariables["y"])
	}
}

func SetupPotentiallyInfiniteModel() State{
	var update0 Update = Update{"x", "=", 1}
	var update1 Update = Update{"x", "=", 0}
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

func TestExplore2(t *testing.T) {
	s := SetupPotentiallyInfiniteModel()
	var list []State = make([]State, 0,0)

	list = Explore(s)
	if len(list) != 2 {
		t.Errorf("len of list should be 2, but got %d", len(list))
	}
	PrintStates(list)
}

func TestStateInformation_GetEssentialInformation(t *testing.T) {
	var s State = State{}
	s.globalVariables = SetupMap()
	s.allTemplates = make([]Template, 0,0)
	s.allTemplates = append(s.allTemplates, SetupCounterModel())

	var si StateInformation = StateInformation{}
	si.globalVariables = make(map[string]int)
	si = si.GetEssentialInformation(s)

	if (s.globalVariables["a"] != si.globalVariables["a"] || s.globalVariables["b"] != si.globalVariables["b"]){
		t.Errorf("function don't work")
	}
}

func TestState_ConfigureState(t *testing.T) {
	var s State = State{}
	var si StateInformation = StateInformation{}
	si = si.GetEssentialInformation(SetupPotentiallyInfiniteModel())

	s = s.ConfigureState(si)

}

func TestClient(t *testing.T) {
	var s State = State{}
	s.globalVariables = SetupMap()
	s.allTemplates = append(s.allTemplates, SetupCounterModel())
	fmt.Println(s.ToString())
	jsonbytes, err := json.Marshal(s)

	if err != nil{
		fmt.Println("something wrong")
	} else{
		fmt.Println(jsonbytes)
	}
	var s1 State = State{}
	json.Unmarshal(jsonbytes, &s1)
	println(s1.ToString())
	println(s1.globalVariables["b"])
}