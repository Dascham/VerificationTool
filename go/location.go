package main

const defaultEdgeAllocation int = 3
var uniqueid int = 0

type Location struct{
	LocationName string
	LocationId int
	Edges []Edge //should be made into slice
	Invariant Invariant
}

func (l Location) AcceptOutGoingEdges(args ...Edge){
	for i := 0; i < len(args); i++{
		l.Edges = append(l.Edges, )
	}
}

func NewLocation(locationName string, i Invariant) Location{
	var a Location = Location{}
	a.LocationName = locationName
	a.Invariant = i
	//id stuff
	a.LocationId = uniqueid
	uniqueid++

	//initialize slice edges
	a.Edges = make([]Edge, defaultEdgeAllocation)

	return a
}
