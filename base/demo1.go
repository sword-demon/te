package main

import "fmt"

func readData(data []byte, x byte) {
	data[0] = x
}

func main() {
	var a [10]byte

	readData(a[0:], 1)
	fmt.Println(a)

	readData(a[1:], 2)
	fmt.Println(a)
}
