package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func reverseString(s string) string {
	return reverse.String(s)
}

func main() {
	ret := reverseString("Hello, OTUS!")
	fmt.Println(ret)
}
