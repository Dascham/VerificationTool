package main

const defaultEdgeAllocation = 0
var uniqueid int = 0

type Location struct{
	LocationName string
	LocationId int
	Edges []Edge //should be made into slice
	Invariant Invariant
	BlockId int
}

func (l Location) AcceptOutGoingEdges(args ...Edge) Location{
	for i := 0; i < len(args); i++{
		l.Edges = append(l.Edges, args[i])
	}
	return l
}

func NewLocation(locationName string, i Invariant) Location{
	var a Location
	a.LocationName = locationName
	a.Invariant = i
	//id stuff
	a.LocationId = uniqueid
	uniqueid++

	//initialize slice edges
	a.Edges = make([]Edge, 0, 0)

	return a
}
