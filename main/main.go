package main

import "fmt"

func test() interface{} {
	return "hello world!"
}

func main() {
	output := test()
	fmt.Println(output.([]byte))
}
