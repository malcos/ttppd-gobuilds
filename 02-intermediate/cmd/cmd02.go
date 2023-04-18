package main

import (
	"fmt"

	"gitlab.com/malcos/demoapp/b"
)

func main() {
	fmt.Println("INTERMEDIATE command")
	fmt.Println("--------------------")
	fmt.Println(b.MessageFromB())
}
