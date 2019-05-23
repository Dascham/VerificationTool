package main

import (
	"fmt"
	"hash/fnv"
)

const MaxValue = 128
const MinValue = -127

func main(){
	var initialState State = SetupSimpleSyncModel()
	var list []State = Explore(initialState)

	fmt.Printf("len of list: %d \n", len(list))

	for i := 0; i < len(list);i++{
		println(list[i].ToString())
	}
}


func remove(a []State, i int) []State {
	a[i] = a[len(a)-1] // Copy last element to index i.
	a[len(a)-1] = State{} // Erase last element (write zero value).
	a = a[:len(a)-1]   //truncate slice

	return a
}

func Explore(initialState State) []State{
	var hashedStates map[string]string = make(map[string]string)
	var waitingList []State = make([]State, 0,0) //size zero, always, cause append fixes size by itself
	waitingList = append(waitingList, initialState)
	var passedList []State = make([]State, 0,0)

	for len(waitingList) > 0 { //exploration loop
		var currentState = waitingList[0]
		waitingList = remove(waitingList, 0)
		for i := 0;i < len(currentState.allTemplates);i++{
			for j := 0 ; j < len(currentState.allTemplates[i].currentLocation.Edges); j++{
				if (currentState.allTemplates[i].currentLocation.Edges[j].EdgeIsActive(currentState.allTemplates[i].LocalVariables, currentState)){
					//check if edge has hand-shake sync
					if (currentState.allTemplates[i].currentLocation.Edges[j].Ch != ""){
						foundedges, ok, templateNumbers := FindSyncEdges(currentState.allTemplates[i].currentLocation.Edges[j],currentState,i)
						if ok { //means we found some edges
							for k, edge := range foundedges{
								//have to instantiate new state here
								newState := DeepCopyState(currentState)
								//do update on new state -> both edges
								newState.allTemplates[i].currentLocation.Edges[j].AtomicUpdate(newState.allTemplates[i].LocalVariables, newState.globalVariables)
								edge.AtomicUpdate(newState.allTemplates[i].LocalVariables, newState.globalVariables)
								//then we advance location of both templates
								newState.allTemplates[i].currentLocation = newState.allTemplates[i].currentLocation.Edges[j].Dst
								newState.allTemplates[templateNumbers[k]].currentLocation = edge.Dst //I thinks 'k' in subscript is correct

								//check if state has been encountered before, must be done after advance
								temp := newState.ToString()
								if _, ok := hashedStates[temp]; ok{
									break; //
								}else{
									hashedStates[temp] = temp
								}
								//add state to waitinglist, add only if map is valid
								if (ValidMap(newState.allTemplates[i].LocalVariables) && ValidMap(newState.globalVariables)){
									waitingList = append(waitingList, newState)
								}
							}
						}
					} else{ //we do basic transition
						//have to instantiate new state here
						newState := DeepCopyState(currentState)
						//do update on new state
						newState.allTemplates[i].currentLocation.Edges[j].AtomicUpdate(newState.allTemplates[i].LocalVariables, newState.globalVariables)

						//then advance location, by looking to dst of edge that we took
						newState.allTemplates[i].currentLocation = newState.allTemplates[i].currentLocation.Edges[j].Dst

						//check if state has been encountered before
						temp := newState.ToString()
						if _, ok := hashedStates[temp]; ok{
							break;
						} else{
							hashedStates[temp] = temp
						}
						//add newstate to waitinglist, for distributed, call distribute function, which hashes and does stuff
						//add only if map is valid, this should be made to better fix
						if (ValidMap(newState.allTemplates[i].LocalVariables) && ValidMap(currentState.globalVariables)) {
						//println("we here?")
						waitingList = append(waitingList, newState)
						}
					}
				}
			}
		}
		//now we have tried everything that is possible for this state, therefore add to passed list
		passedList = append(passedList, currentState)
	}
	return passedList
}

//helper function, for finding an synchronization partner for an edge
func FindSyncEdges(e Edge, currentState State, activeTemplate int) ([]Edge, bool, []int){
	var foundEdges []Edge = make([]Edge, 0,0)
	var result bool = true
	var templateNumbers []int = make([]int, 0,0)
	//essentially, find all edges, which match
	for i, template := range currentState.allTemplates{
		if (i == activeTemplate){ //skip the template, that we are "in"
			continue;
		}else {
			for _, edge := range template.currentLocation.Edges {
				if (e.IsSend){
					if (!edge.IsSend && edge.Ch == e.Ch && edge.EdgeIsActive(template.LocalVariables, currentState)) { //is not send
						foundEdges = append(foundEdges, edge)
						templateNumbers = append(templateNumbers, i)
					}
				} /*else {
					if (edge.IsSend && edge.Ch == e.Ch && edge.EdgeIsActive(template.LocalVariables, currentState)){
						foundEdges = append(foundEdges, edge)
						templateNumbers = append(templateNumbers, i)
					}
				}
				*/
			}
		}
	}

	if len(foundEdges) == 0{
		result = false
	}

	return foundEdges, result, templateNumbers
}


func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func MainSetupCounterModel() Template{
	var localVariables map[string]int = map[string]int{"x":0}
	var template Template = Template{}
	template.LocalVariables = localVariables
	var location0 Location = NewLocation("L0", Invariant{})
	template.InitialLocation = &location0
	template.currentLocation = &location0

	//update
	var update Update = Update{"x", "++", 0}
	//edge
	var edge Edge = Edge{}
	edge = edge.InitializeEdge()
	edge = edge.AcceptUpdates(update)
	edge = edge.AssignSrcDst(&location0, &location0)

	location0 = location0.AcceptOutGoingEdges(edge)
	return template
}

func MainSetupMap() map[string]int {
	var localVariables map[string]int = make(map[string]int)
	localVariables["a"] = 2
	localVariables["b"] = 7
	return localVariables
}

func CopyMap(originalMap map[string]int) map[string]int{
	var newMap map[string]int = make(map[string]int)
	for key, value := range originalMap {
		newMap[key] = value
	}
	return newMap
}

func DeepCopyState(s State) State{
	var newState State = State{}
	//newState.allTemplates = s.allTemplates

	//old way
	newState.allTemplates = make([]Template, 0,0)
	//copy templates

	for i := 0; i<len(s.allTemplates);i++{
		newState.allTemplates = append(newState.allTemplates, s.allTemplates[i])
	}


	newState.globalVariables = CopyMap(s.globalVariables)

	for i := 0; i<len(s.allTemplates);i++ {
		newState.allTemplates[i].LocalVariables = CopyMap(s.allTemplates[i].LocalVariables)
	}

	return newState
}

func SetupSimpleSyncModel() State{
	var initialState State = State{}
	initialState.globalVariables = map[string]int{"y":8,"z":5}

	//guards
	var guard0 Guard = Guard{"x", "==", 0}
	var guard1 Guard = Guard{"z", "<", 10}

	var update0 Update = Update{"y", "++", 0}
	var update1 Update = Update{"y", "*", 2}

	var invariant0 Invariant = Invariant{"z", "<", 10}

	var location0 Location = NewLocation("L0", Invariant{})
	var Location1 Location = NewLocation("L1", Invariant{})
	var location2 Location = NewLocation("L2", Invariant{})
	var location3 Location = NewLocation("L3", invariant0)


	var edge0 Edge = Edge{}
	edge0 = edge0.InitializeEdge()
	edge0 = edge0.AcceptGuards(guard0)
	edge0 = edge0.AcceptUpdates(update0)
	edge0 = edge0.AssignSrcDst(&location0, &Location1)
	edge0.Ch = "a"
	edge0.IsSend = true

	var edge1 Edge = Edge{}
	edge1 = edge1.InitializeEdge()
	edge1 = edge1.AcceptGuards(guard1)
	edge1 = edge1.AcceptUpdates(update1)
	edge1 = edge1.AssignSrcDst(&location2, &location3)
	edge1.Ch = "a"
	edge1.IsSend = false

	location0 = location0.AcceptOutGoingEdges(edge0)
	location2 = location2.AcceptOutGoingEdges(edge1)

	//template0
	var template0 Template = Template{}
	template0.LocalVariables = map[string]int{"x":0}
	template0.InitialLocation = &location0
	template0.currentLocation = &location0

	var template1 Template = Template{}
	template1.LocalVariables = make(map[string]int) //let be empty
	template1.InitialLocation = &location2
	template1.currentLocation = &location2

	initialState.allTemplates = append(initialState.allTemplates, template0, template1)

	return initialState
}