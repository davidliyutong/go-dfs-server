package main

import "fmt"

func main() {
	var a = make([]int, 1)
	fmt.Println(a[:0])
}
