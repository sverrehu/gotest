package main

import "fmt"

type Interface interface {
	GetNextNumber() int
}

type Implementation struct {
	number int
}

func (i *Implementation) GetNextNumber() int {
	i.number += 1
	return i.number
}

func main() {
	var i Interface
	//i = Implementation{} // type does not implement interface
	i = &Implementation{} // works as expected, may modify contents of the struct
	fmt.Println(i.GetNextNumber())
	fmt.Println(i.GetNextNumber())
}
