package main

type Location struct{
	LocationName string
	LocationId int
	Edge []Edge //should be made into slice
	Invariant Invariant
}