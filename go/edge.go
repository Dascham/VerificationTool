package main

//an edge consists of src node, dst node, a possible guard, a possible synchronization,
// and a possible update on local or global variables.
type Edge struct {
	Src *Location
	Dst *Location
	Guard Guard
	Ch string
	IsSend bool
	Update Update
}

//! is send
//? is receive

func (e Edge) EdgeIsActive(localVariables map[string]int) bool{
	var result bool = false
	if (e.Guard.Evaluate(localVariables) &&
		e.Dst.Invariant.IsValid(localVariables)){ //add one more for chan
		result = true
	}
	return result
}