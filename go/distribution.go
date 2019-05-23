package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

func Client() {
	//prep some data to send
	var s State = State{}
	jsonbytes, err := json.Marshal(s)
	println("Printing json string")
	fmt.Println(string(jsonbytes))
	//var address string = "127.0.0.1:"+os.Args[1] //should be first argument

	if (err != nil) {
		fmt.Printf("Marshall error: %s\n", err)
	}

	conn, err1 := net.Dial("tcp", "127.0.0.1:5000")
	fmt.Printf("Dialed\n")
	if err1 != nil {
		// handle error
		fmt.Printf("Something went wrong %s \n", err)
	}
	fmt.Printf("Going to send: ")
	println(jsonbytes)
	_, err2 := conn.Write(jsonbytes)
	if err2 != nil{
		fmt.Printf("Error: %s", err2)
	}
	err = conn.Close()
	if (err != nil) {
		fmt.Printf("printing error: %s", err)
	}
}


func Server() {
	channel := make(chan State)

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
	}
	s := <- channel
	println(s.ToString())

}

func handleConnection(conn net.Conn, channel chan State) {
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
	var s State
	err1 := json.Unmarshal(buff.Bytes(), &s)
	if err1 != nil{
		fmt.Printf("Error: %s\n",err1)
	}
	fmt.Println(s.ToString())

	fmt.Println("Closing connection: ")
	err2 := conn.Close()
	if err2 != nil{
		fmt.Printf("Could not close connection: %s", err2)
	} else {
		fmt.Println("Connection closed.")
	}
	channel <- s //send state on channel 'channel'
}

type StateInformation struct{
	globalVariables map[string]int
	localVariables []map[string]int
	currentLocationIds []int
}

func (si StateInformation) GetEssentialInformation(s State) StateInformation{
	si.globalVariables = s.globalVariables

	for i, template := range s.allTemplates{
		si.localVariables[i] = template.LocalVariables
		si.currentLocationIds[i] = template.currentLocation.LocationId
	}
	return si
}

