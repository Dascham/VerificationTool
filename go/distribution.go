package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

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
