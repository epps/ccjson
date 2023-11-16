package main

import (
	"fmt"
)

func main() {
	fmt.Println("Parsing JSON ...")
	input := "{}"

	p := NewParser(input)

	output, err := p.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Parsed output: %v\n", output)
}
