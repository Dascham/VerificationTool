package main

//an edge consists of src node, dst node, a possible guard, a possible synchronization,
// and a possible update on local or global variables.
type Edge struct {
	Src *Location
	Dst *Location
	Guard []Guard
	Ch string
	IsSend bool
	Update []Update
}

//! is send
//? is receive

func (e Edge) InitializeEdge(){
	e.Guard = make([]Guard, 2)
	e.Update = make([]Update, 2)
}

func (e Edge) AcceptUpdates(args ...Update){
	for i := 0; i < len(args); i++{
		e.Update = append(e.Update, args[i])
	}
}
func (e Edge) AcceptGuards(args ...Guard){
	for i := 0; i < len(args); i++{
		e.Guard = append(e.Guard, args[i])
	}
}

func (e Edge) AssignSrcDst(src Location, dst Location){
	e.Src = &src
	e.Dst = &dst
}

func (e Edge) EdgeIsActive(localVariables map[string]int) bool{
	var result bool = false
	for i := 0; i < len(e.Guard); i++ {
		if (e.Guard[i].Evaluate(localVariables) &&
			e.Dst.Invariant.IsValid(localVariables)) { //add one more for chan
			result = true
			}
		}
	return result
}