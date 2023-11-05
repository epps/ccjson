package main

import (
	"fmt"
)

func main() {
	fmt.Println("Parsing JSON ...")
	input := "{}"

	_, err := ParseJSON(input)
	if err != nil {
		fmt.Println(err)
	}
}
