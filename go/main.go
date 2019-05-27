package main

import (
	"fmt"
	"time"
)

//IF A STRUCT IS TO BE MARSHALLED, ONE SHOULD "EXPORT" THE FIELDS/ATTRIBUTES OF THAT STRUCT
//EXPORTING THE FIELDS IS DONE BY HAVING THE NAME OF THE FIELDS, CAPITALIZED. CAPITALIZATION OF THE NAME OF A
//FIELD, EXPORTS IT.

var selfNodeNumber int = 0

func main(){
	start := time.Now()
	s := EmptyState()
	s.globalVariables = make(map[string]int)
	s.allTemplates = append(s.allTemplates, SetupCounterModel(), SetupCounterModel(), SetupCounterModel())
	list := Explore(s)
	print("Number of states explored: ")
	println(len(list))
	println(time.Since(start).String())
}
func Master(){
	selfNodeNumber = 0
	//assign numbers to nodes
	initializeNodes(ipaddresses)
	println("Initialization done")

	list := ExploreDistributed(ParallelSetup(SetupCounterModel()))
	PrintStates(list)
	print("Explored states: ")
	println(len(list))

}
func Node(){
	GetInitialized()
	list := ExploreDistributed(ParallelSetup(SetupCounterModel()))
	PrintStates(list)
	print("Explored states: ")
	println(len(list))
}

func ExploreDistributed(initialState State) []State{
	//var mutex sync.Mutex = sync.Mutex{}
	var channel chan State = make(chan State, 10.000) //buffer size 10.000, used for transmitting states
	//var chanDonezo chan bool = make(chan bool)

	var waitingList []State = make([]State, 0,0) //size zero, always, cause append fixes size by itself
	var passedList []State = make([]State, 0,0)
	var hashedStates map[string]string = make(map[string]string)

	go ReceiveStates(channel, DeepCopyState(initialState)) //this concurrently receives states from the network, and puts them in a buffered channel
	fmt.Printf("The initialstate: %s\n", initialState.ToString())
	if (selfNodeNumber != 0){ //this blocks non-master nodes from exploring, until they receive a state
		initialState = <- channel

	}
	//master starts the exploration
	hashedStates[initialState.ToString()] = initialState.ToString()
	waitingList = append(waitingList, initialState)

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
									num := Hash(newState.ToString())
									if (num%lenOfIpaddreses == uint32(selfNodeNumber)){
										waitingList = append(waitingList, newState)
									}else{
										SendAState(newState, num)
									}
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
						if ValidMap(newState.allTemplates[i].LocalVariables) && ValidMap(newState.globalVariables){
							var num uint32 = Hash(newState.ToString())
							var sendToNode uint32 = num%lenOfIpaddreses
							if (sendToNode == uint32(selfNodeNumber)){
								waitingList = append(waitingList, newState)
							}else{
								SendAState(newState, sendToNode)
							}
						}
					}
				}
			}
		}
		//now we have tried everything that is possible for this state, therefore add to passed list
		passedList = append(passedList, currentState)

		//put all states received from other nodes in waitinglist
		tempList := FromChannelToList(channel)
		for _, state := range tempList{
			_, ok := hashedStates[state.ToString()] //ok is true if state is seen
			if !ok { //if not okay, then add stuff, otherwise skip
				hashedStates[state.ToString()] = state.ToString()
				waitingList = append(waitingList, state)
			}
		}

		if(len(waitingList) == 0){
			time.Sleep(900*time.Millisecond) //if empty we wait a bit, to see if other machines send some states that need exploration
		}
		tempList1 := FromChannelToList(channel)
		waitingList = append(waitingList, tempList1...)
		//if waitingList still empty, we prolly done with exploring
	}
	/*
	if (selfNodeNumber == 0){
		passedList = append(passedList, MasterReceiveExploredStates(initialState)...)
	}
	if (selfNodeNumber != 0){
		//chanDonezo <- true //this should be redundant now?
		NodeSendExploredStates(passedList)
	}
	 */

	return passedList
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