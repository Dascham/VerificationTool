package main

import (
	"testing"
)

func SetupMap() map[string]int {
	var localVariables map[string]int = make(map[string]int)
	localVariables["a"] = 2
	localVariables["b"] = 7
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
	var localVariables map[string]int = map[string]int{"x":0}
	var template Template = Template{}
	template.LocalVariables = localVariables
	var location0 Location = NewLocation("L0", Invariant{})
	template.InitialLocation = &location0

	//update
	var update Update = Update{"x", "++", 0}
	//edge
	var edge Edge = Edge{}
	edge = edge.InitializeEdge()
	edge = edge.AcceptUpdates(update)
	edge = edge.AssignSrcDst(location0, location0)

	location0 = location0.AcceptOutGoingEdges(edge)

	return template
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
	edge0 = edge0.AssignSrcDst(location0, location0)
	var edge1 = Edge{}
	edge1 = edge1.InitializeEdge()
	edge1 = edge1.AcceptGuards(guard0)
	edge1 = edge1.AcceptUpdates(update1)
	edge1 = edge1.AssignSrcDst(location0, location1)
	var edge2 = Edge{}
	edge2 = edge2.InitializeEdge()
	edge2 = edge2.AcceptUpdates(update2)
	edge2 = edge2.AssignSrcDst(location1, location1)
	var edge3 = Edge{}
	edge3 = edge3.InitializeEdge()
	edge3 = edge3.AcceptGuards(guard1)
	edge3 = edge3.AcceptUpdates(update3)
	edge3 = edge3.AssignSrcDst(location1, location2)
	var edge4 = Edge{}
	edge4 = edge4.InitializeEdge()
	edge4 = edge4.AcceptUpdates(update4)
	edge4 = edge4.AssignSrcDst(location2, location2)
	var edge5 = Edge{}
	edge5 = edge5.InitializeEdge()
	edge5 = edge5.AcceptGuards(guard2)
	edge5 = edge5.AcceptUpdates(update5)
	edge5 = edge5.AssignSrcDst(location2, location3)
	var edge6 = Edge{}
	edge6 = edge6.InitializeEdge()
	edge6 = edge6.AcceptUpdates(update6)
	edge6 = edge6.AssignSrcDst(location3, location3)
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

//this test should be the first test to call the function SetupTemplate(), otherwise it will fail
func TestTemplate_ToString(t *testing.T) {
	var template Template = SetupTemplate()
	var expected string = "027"
	var s string = template.ToString()

	if (s != expected) {
		t.Errorf("Got: %s --- expected: %s", s, expected)
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
	s1.allTemplates[0].InitialLocation.Edges[0].Update[0].Update(s1.allTemplates[0].LocalVariables)
	//println(s.allTemplates[0].LocalVariables["x"])
	//println(s1.allTemplates[0].LocalVariables["x"])

	if s.allTemplates[0].LocalVariables["x"] == s1.allTemplates[0].LocalVariables["x"]{
		t.Errorf("AtomicUpdate does not work: %d == %d",
			s.allTemplates[0].LocalVariables["x"], s1.allTemplates[0].LocalVariables["x"])
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

func TestHash(t *testing.T) {
	var a string = "aksdksd<jfnjikdfhvjikdvhfjdvh"
	var b string = "sdkfj<hvuifdhvuizhvfbvhfbvvbvbvbvbvbvbvbvbvbvbvbvbvbvbvbvb"

	c := Hash(a)
	d := Hash(b)
	if c == d{
		t.Errorf("Hash function returns same value of different strings, but should not: '%d' '%d'", c,d)
	}
}