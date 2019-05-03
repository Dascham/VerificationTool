package main

import (
	"hash/fnv"
)

func main(){
	//have global mutex, in order to change global state,
	// although not necessary for passed and waiting list implementation

}
func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}