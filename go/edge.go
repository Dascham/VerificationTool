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

func (e Edge) AssignSrcDst(src Location, dst Location) Edge{
	e.Src = &src
	e.Dst = &dst
	return e
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

func (e Edge) AtomicUpdate(localVariables map[string]int) Edge{
	for i:=0;i<len(e.Update);i++{
		e.Update[i].Update(localVariables)
	}
	return e
}
