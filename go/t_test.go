package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
)

//this test should be the first test to call the function SetupTemplate(), otherwise it will fail
func TestTemplate_ToString(t *testing.T) {
	var template Template = SetupTemplate()
	var expected string = "Location: 0 a:2 b:7 "
	var s string = template.ToString()
	if (s != expected) {
		t.Errorf("Got: %s --- expected: %s", s, expected)
	}
}
func TestGuard_Evaluate(t *testing.T) {
	var localVariables map[string]int = SetupMap()
	var g Guard = Guard{"a", "<", 3,""}
	if g.Evaluate(localVariables) != true{
		t.Errorf("Evaluate function should return true, but returned false")
	}
}
func TestGuard_Evaluate2(t *testing.T) {
	var localVariables map[string]int = SetupMap()
	var g Guard = Guard{"b", "<", 3,""}
	if g.Evaluate(localVariables) == true{
		t.Errorf("Evaluate function should return false, but returned true")
	}
}
func TestUpdate_Update(t *testing.T) {
	var localVariables map[string]int = SetupMap()
	var u1 Update = Update{"x", "+=", 2,""}
	var u2 Update = Update{"x", "+", 2,""}
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

	s1.allTemplates[0].InitialLocation.Edges[0].AtomicUpdate(s1.allTemplates[0].LocalVariables, s1.globalVariables)

	if s.allTemplates[0].LocalVariables["x"] == s1.allTemplates[0].LocalVariables["x"]{
		t.Errorf("AtomicUpdate does not work: %d == %d",
			s.allTemplates[0].LocalVariables["x"], s1.allTemplates[0].LocalVariables["x"])
	}
}

func TestInvariant_IsValid(t *testing.T) {
	//two maps
	var map1 map[string]int = map[string]int{"f":5}
	//var map2 map[]
	var i Invariant = Invariant{"f", "<", 10,""}

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
func TestExplore2(t *testing.T) {
	s := SetupPotentiallyInfiniteModel()
	var list []State = make([]State, 0,0)

	list = Explore(s)
	if len(list) != 2 {
		t.Errorf("len of list should be 2, but got %d", len(list))
	}
	PrintStates(list)
}
func TestState_ConfigureState(t *testing.T) {
	var template Template = SetupCounterModel()
	var s State = State{}
	s.globalVariables = SetupMap()
	s.allTemplates = make([]Template, 0)
	s.allTemplates = append(s.allTemplates, template)

	var si StateInformation = StateInformation{}
	si = si.GetEssentialInformation(s)

	temp := s.ToString()
	s.globalVariables["a"] = 15
	s = s.ConfigureState(si)
	if temp != s.ToString(){
		t.Errorf("stateinformation did not properly reconfigure our state s")
	}
}
func TestMarshal(t *testing.T) {
	var slice []map[string]int = make([]map[string]int, 0,0)
	slice = append(slice, SetupMap())
	slice = append(slice, SetupMap())

	fmt.Println(slice)

	b, err := json.Marshal(slice)

	fmt.Println(err)

	fmt.Print("printing json: ")
	fmt.Println(b)

	var map1 []map[string]int = make([]map[string]int,0,0)
	err1 := json.Unmarshal(b, &map1)

	fmt.Println(err1)

	fmt.Println(map1)
}
func TestMarshal2(t *testing.T){
	//marshal stateinformation struct
	var si StateInformation
	si.GlobalVariables = make(map[string]int)
	si.GlobalVariables = SetupMap()

	si.ListLocalVariables = make([]map[string]int, 0,0)
	si.ListLocalVariables = append(si.ListLocalVariables, SetupMap1(), SetupMap2())
	var slice []int = []int{2, 5, 7, 9}
	si.CurrentLocationIds = slice

	//print si information,
	jsonbytes, _:= json.Marshal(si)

	var si1 StateInformation
	si1.GlobalVariables = make(map[string]int)
	_ = json.Unmarshal(jsonbytes, &si1)

}
func TestNetworkCommunication(t *testing.T){
	//run main in terminal window
	//Server()
	//fmt.Println("Everything Done")
}

func TestConcurrentListAdd(t *testing.T){
	var channel chan string = make(chan string, 10)
	for i:=0;i<10;i++{
		channel <- "vip "+strconv.Itoa(i)
	}
	var slice []string = make([]string,0,0)
	for {
		select{
		case string:= <- channel:
			slice = append(slice, string)
		default:
			return
		}
	}
}

//in initializeNodes and getinitialized we convert a lot, test this
func TestConversions(t *testing.T){

}

//test buffer writing and &si
func TestBuffer(t *testing.T){
	s := SetupSimpleSyncModel()
	s.globalVariables = SetupMap3()
	var si StateInformation = StateInformation{}
	si = si.GetEssentialInformation(s)
	jsonbytes,_ := json.Marshal(si)

	var buff bytes.Buffer

	//single write is okay, but double is bad
	buff.Write(jsonbytes)
	buff.Write(jsonbytes)

	si1 := StateInformation{}
	//json.Unmarshal(buff.ReadBytes([]byte(io.EOF)), &si1)

	fmt.Println(si1.GlobalVariables)
}

func TestExplore3(t *testing.T) {
	s := EmptyState()
	s.globalVariables = make(map[string]int)
	s.allTemplates = append(s.allTemplates, SetupCounterModel(), SetupCounterModel(), SetupCounterModel())
	list := Explore(s)
	print("Number of states explored: ")
	println(len(list))

}

