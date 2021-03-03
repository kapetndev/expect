package main

import "fmt"

func main() {
	fmt.Println(greeting("Picard"))
}

func greeting(name string) string {
	return "Hello, " + name
}
