package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go-envdir path/to/env-dir command <arg1> <arg2> ...")
		os.Exit(1)
	}
	dir := os.Args[1]
	env, err := ReadDir(dir)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error reading env dir:", err)
		os.Exit(1)
	}

	code := RunCmd(os.Args[2:], env)
	os.Exit(code)
}
