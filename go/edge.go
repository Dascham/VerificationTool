package main

import "strings"

//an edge consists of src node, dst node, a possible guard, a possible synchronization,
// and a possible update on local or global variables.
type Edge struct {
	Src *Location
	Dst *Location
	Guard []Guard
	Ch string
	IsSend bool
	Update []Update
	name string
}

//! is send
//? is receive

func (e Edge) InitializeEdge() Edge{
	e.Guard = make([]Guard, 0, 0)
	e.Update = make([]Update, 0, 0)
	return e
}

func (e Edge) AcceptUpdates(args ...Update) Edge{
	for i := 0; i < len(args); i++{
		e.Update = append(e.Update, args[i])
	}
	return e
}
func (e Edge) AcceptGuards(args ...Guard) Edge{
	for i := 0; i < len(args); i++{
		e.Guard = append(e.Guard, args[i])
	}
	return e
}

func (e Edge) AssignSrcDst(src *Location, dst *Location) Edge{
	e.Src = src
	e.Dst = dst
	return e
}

func (e Edge) EdgeIsActive(localVariables map[string]int, s State) bool{
	var tempMap = CopyMap(localVariables)
	var result bool = true
	for i := 0; i < len(e.Guard); i++ {
		if (!e.Guard[i].Evaluate(localVariables) || !e.Guard[i].Evaluate(s.globalVariables)) {
			return false
		}
	}
	//eval channels, not done yet

	//then eval dst invariant, where we need to update first, and then check if invariant valid
	e.AtomicUpdate(tempMap, s.globalVariables)
	if (!e.Dst.Invariant.IsValid(tempMap) || !e.Dst.Invariant.IsValid(s.globalVariables)) { //add one more for chan
			return false
		}
	return result
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
func ValidValue(a int) bool{
	if (MinValue < a && a < MaxValue){
		return true
	}else {
		return false
	}
}

func (e Edge) AtomicUpdate(localVariables, globalVariables map[string]int) Edge{
	for i:=0;i<len(e.Update);i++{
		e.Update[i].Update(localVariables)
		e.Update[i].Update(globalVariables)
	}
	return e
}
func (e Edge) ToString()string{
	var sb strings.Builder
	for i := 0; i < len(e.Guard);i++ {
		sb.WriteString(e.Guard[i].ToString())
	}
	for j:=0; j<len(e.Update);j++{
		sb.WriteString(e.Update[j].ToString())
	}

	return sb.String()
}
