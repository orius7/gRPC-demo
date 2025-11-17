package main

import "fmt"

func main() {
	fmt.Println("Practice package")

	var xx [4]int32
	fmt.Println(xx[0])

	xx[1] = 42
	fmt.Println(&xx[1])
}
