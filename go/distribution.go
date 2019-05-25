package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
)
//Ip addresses, master is always [0]
var ipaddresses []string = []string{"127.0.0.1", "172.28.211.53"}
var portNumbers1 []string = []string{":5000", ":5001", "5002"}
var lenOfIpaddreses uint32 = uint32(len(ipaddresses))

//StateInformation contains the information transmitted between nodes
type StateInformation struct{
	GlobalVariables map[string]int
	ListLocalVariables []map[string]int
	CurrentLocationIds []int
}

func (si StateInformation) GetEssentialInformation(s State) StateInformation{
	si.GlobalVariables = CopyMap(s.globalVariables)

	for _, template := range s.allTemplates{
		si.ListLocalVariables = append(si.ListLocalVariables, CopyMap(template.LocalVariables))
		si.CurrentLocationIds = append(si.CurrentLocationIds, template.currentLocation.LocationId)
	}
	return si
}
func initializeNodes(ipaddresses []string) {
	for i, address := range ipaddresses {
		if (i == 0){
			continue
		} else {
			conn, err := net.Dial("tcp", address+portNumbers1[0]) //0 is portnumber 5000
			if err != nil {
				fmt.Printf("Something wrong when dialing, initializeNode,: %s\n", err)
			}
			_, err1 := conn.Write([]byte(strconv.Itoa(i)))
			if err1 != nil {
				fmt.Printf("Something wrong when trying to conn.write: %s\n", err1)
			}
			err2 := conn.Close()
			if err2 != nil {
				fmt.Printf("Could not close connection: %s", err2)
			}
		}
	}
}
func GetInitialized(){
	ln, err := net.Listen("tcp", portNumbers1[0]) //portnumbers1[0], is port 5000
	println("listening")
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
	}

	conn, err1 := ln.Accept()
	if err1 != nil {
		fmt.Printf("Second layer of wrong: %s", err1)
	}
	var buff bytes.Buffer
	_, err2 := io.Copy(&buff, conn)
	if err2 != nil{
		fmt.Printf("Something with io.copy: %s", err2)
	}
	num, err3 := strconv.Atoi(buff.String())
	if err3 != nil{
		fmt.Printf("could not convert: %s", err3)
	}
	err4 := conn.Close()
	if err4 != nil{
		fmt.Printf("Could not close connection: %s", err4)
	}
	selfNodeNumber = num
}

func SendAState(s State, sendToNode uint32){
	//get stateinformation
	var si StateInformation = StateInformation{}
	si = si.GetEssentialInformation(s)

	json_si, err := json.Marshal(si)
	if err != nil{
		fmt.Println("ehhhhhhhhhhhhhhhhhhhhhhhhh")
	}

	//send
	conn, err1 := net.Dial("tcp", ipaddresses[sendToNode]+portNumbers1[1]) //5001

	if err1 != nil{
		fmt.Printf("Dial went wrong in SendState: %s", err1)
	}

	_, err2 := conn.Write(json_si)
	if err2 != nil{
		fmt.Printf("write jsonbytes went wrong: %s", err2)
	}
	err3 := conn.Close()
	if err3 != nil{
		fmt.Printf("Could not close connection: %s", err3)
	}
}
func ReceiveStates(channel chan State, s State) {
	ln, err := net.Listen("tcp", portNumbers1[1]) //5001
	println("listening")
	if err != nil {
		fmt.Printf("Something went wrong")
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Second layer of wrong")
		}
		go PutStateInChannel(conn, channel, s)
	}
}
func PutStateInChannel(conn net.Conn, channel chan State, s State) {
	//buffer := new(bytes.Buffer)
	//var msg []byte = make([]byte, 500)
	var buff bytes.Buffer
	_, err := io.Copy(&buff, conn)
	if err != nil{
		fmt.Printf("handleconnection error, something: %s \n", err)
		err1 := conn.Close()
		if err1 != nil{
			fmt.Printf("%s\n", err1)
		}
	}
	var si StateInformation
	err2 := json.Unmarshal(buff.Bytes(), &si)
	if err2 != nil{
		fmt.Printf( "%s\n",err2)
	}
	s = s.ConfigureState(si)

	err3 := conn.Close()
	if err3 != nil{
		fmt.Printf("Could not close connection: %s", err3)
	}
	channel <- s //send state on channel 'channel'
}

func SendExploredStates(){

}

func ReceiveExploredStates() []State{
	var list []State = make([]State, 0,0)

	return list
}


func Client() {
	var template Template = MainSetupCounterModel()
	var s State = State{}
	s.globalVariables = map[string]int{"x":5}
	s.allTemplates = make([]Template, 0,0)
	s.allTemplates = append(s.allTemplates, template)

	var si StateInformation = StateInformation{}
	si = si.GetEssentialInformation(s)

	jsonbytes, err := json.Marshal(si)
	if (err != nil) {
		fmt.Printf("Marshall error: %s\n", err)
	}
	conn, err1 := net.Dial("tcp", "172.28.211.53:5000")
	fmt.Println("Dialed")
	if err1 != nil {
		fmt.Printf("Something went wrong %s \n", err)
	}
	_, err2 := conn.Write(jsonbytes)
	if err2 != nil{
		fmt.Printf("Error: %s", err2)
	}
	err = conn.Close()
	if (err != nil) {
		fmt.Printf("printing error: %s", err)
	}
	fmt.Println("Client done")
}
func Client2(){
	var template Template = MainSetupCounterModel()
	var s State = State{}
	s.globalVariables = map[string]int{"x":5}
	s.allTemplates = make([]Template, 0,0)
	s.allTemplates = append(s.allTemplates, template)

	var si StateInformation = StateInformation{}
	si = si.GetEssentialInformation(s)

	jsonbytes, err := json.Marshal(si)
	if (err != nil) {
		fmt.Printf("Marshall error: %s\n", err)
	}
	conn, err1 := net.Dial("tcp", "127.0.0.1:5000")
	fmt.Println("Dialed")
	if err1 != nil {
		fmt.Printf("Something went wrong %s \n", err)
	}
	_, err2 := conn.Write(jsonbytes)
	if err2 != nil{
		fmt.Printf("Error: %s", err2)
	}
	err = conn.Close()
	if (err != nil) {
		fmt.Printf("printing error: %s", err)
	}
	fmt.Println("Client done")
}

func Server() {
	channel := make(chan StateInformation)

	ReceiveAndPrint := func() {
		s := <- channel
		fmt.Println(s.GlobalVariables)
		fmt.Println(s.ListLocalVariables)
		fmt.Println(s.CurrentLocationIds)
	}

	ln, err := net.Listen("tcp", ":5000")
	println("listening")
	if err != nil {
		fmt.Printf("Something went wrong")
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Second layer of wrong")
		}
		go handleConnection(conn, channel)
		println("Handling new connection!")
		go ReceiveAndPrint()
	}
}

func handleConnection(conn net.Conn, channel chan StateInformation) {
	//buffer := new(bytes.Buffer)
	//var msg []byte = make([]byte, 500)
	var buff bytes.Buffer
	_, err := io.Copy(&buff, conn)
	if err != nil{
		fmt.Printf("handleconncetion error, something \n")
		conn.Close()
	}
	fmt.Printf("Received the following message: %s", string(buff.Bytes()))

	fmt.Printf("Unmarshalling\n")
	var s StateInformation
	err1 := json.Unmarshal(buff.Bytes(), &s)
	if err1 != nil{
		fmt.Printf("Error: %s\n",err1)
	}

	fmt.Println("Closing connection: ")
	err2 := conn.Close()
	if err2 != nil{
		fmt.Printf("Could not close connection: %s", err2)
	} else {
		fmt.Println("Connection closed.")
	}
	channel <- s //send state on channel 'channel'
}


