package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("expect filename as argument")
	}
	filename := os.Args[1]

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	root, err := Parse(bufio.NewReader(f))
	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}

	fmt.Printf("root:\n%s\n", root.String())
}
