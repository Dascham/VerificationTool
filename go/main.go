package main

import (
	"hash/fnv"
)

func main(){
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
	edge0.InitializeEdge()
	edge0.AcceptUpdates(update0)
	edge0.AssignSrcDst(location0, location0)
	var edge1 = Edge{}
	edge1.InitializeEdge()
	edge1.AcceptGuards(guard0)
	edge1.AcceptUpdates(update1)
	edge1.AssignSrcDst(location0, location1)
	var edge2 = Edge{}
	edge2.InitializeEdge()
	edge2.AcceptUpdates(update2)
	edge2.AssignSrcDst(location1, location1)
	var edge3 = Edge{}
	edge3.InitializeEdge()
	edge3.AcceptGuards(guard1)
	edge3.AcceptUpdates(update3)
	edge3.AssignSrcDst(location1, location2)
	var edge4 = Edge{}
	edge4.InitializeEdge()
	edge4.AcceptUpdates(update4)
	edge4.AssignSrcDst(location2, location2)
	var edge5 = Edge{}
	edge5.InitializeEdge()
	edge5.AcceptGuards(guard2)
	edge5.AcceptUpdates(update5)
	edge5.AssignSrcDst(location2, location3)
	var edge6 = Edge{}
	edge6.InitializeEdge()
	edge6.AcceptUpdates(update6)
	edge6.AssignSrcDst(location3, location3)
	var edge7 = Edge{}
	edge7.InitializeEdge()
	edge7.AcceptGuards(guard3)
	edge7.AcceptUpdates(update7, update8)
	//locations
	location0.AcceptOutGoingEdges(edge0, edge1)
	location1.AcceptOutGoingEdges(edge2, edge3)
	location2.AcceptOutGoingEdges(edge4, edge5)
	location3.AcceptOutGoingEdges(edge6, edge7)

	var template Template = Template{}
	template.InitialLocation = &location0

}
func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}