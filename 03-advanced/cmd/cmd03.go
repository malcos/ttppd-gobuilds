package main

import (
	"fmt"

	"gitlab.com/malcos/demoapp/a"
	"gitlab.com/malcos/demoapp/b"
)

func main() {
	fmt.Println("ADVANCED SINGLE command")
	fmt.Println("--------------------")
	fmt.Println(a.MessageFromA())
	fmt.Println(b.MessageFromB())
}
