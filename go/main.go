package main

//IF A STRUCT IS TO BE MARSHALLED, ONE SHOULD "EXPORT" THE FIELDS/ATTRIBUTES OF THAT STRUCT
//EXPORTING THE FIELDS IS DONE BY HAVING THE NAME OF THE FIELDS, CAPITALIZED. CAPITALIZATION OF THE NAME OF A
//FIELD, EXPORTS IT.

func main(){
	list := Explore(SetupSimpleSyncModel())
	PrintStates(list)
}

func Explore(initialState State) []State{
	var hashedStates map[string]string = make(map[string]string)
	hashedStates[initialState.ToString()] = initialState.ToString()
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