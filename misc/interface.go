package main

import "fmt"

type Interface interface {
	GetNextNumber() int
}

type Implementation struct {
	number int
}

// must have a pointer receiver to be able to change values inside the struct
func (i *Implementation) GetNextNumber() int {
	i.number += 1
	return i.number
}

func create() *Interface {
	var a Interface
	a = &Implementation{}
	b := &a
	return b
}

func main() {
	var i Interface
	//i = Implementation{} // type does not implement interface
	i = &Implementation{} // works as expected, may modify contents of the struct
	fmt.Println(i.GetNextNumber())
	fmt.Println(i.GetNextNumber())
	switch i.(type) {
	case *Implementation:
		fmt.Printf("implementation %d\n", i.GetNextNumber())
		i.(*Implementation).number = 100
	}
	fmt.Println(i.GetNextNumber())
}
