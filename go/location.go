package main

const defaultEdgeAllocation = 1
var uniqueid int = 0

type Location struct{
	LocationName string
	LocationId int
	Edges []Edge //should be made into slice
	Invariant Invariant
}

func (l Location) AcceptOutGoingEdges(args ...Edge){
	
	for _, value := range args{
		l.Edges = append(l.Edges, value)
	}
}

func NewLocation(locationName string, i Invariant) Location{
	var a Location
	a.LocationName = locationName
	a.Invariant = i
	//id stuff
	a.LocationId = uniqueid
	uniqueid++

	//initialize slice edges
	a.Edges = make([]Edge, defaultEdgeAllocation)

	return a
}
