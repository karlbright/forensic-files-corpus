package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	case "strip":
		strip()
	case "pick":
		pick()
	case "generate":
		generate()
	default:
		usage()
	}

	os.Exit(1)
}

func usage() {
	fmt.Println("USAGE: ffcorpus generate sentences.txt [min] [max]")
	fmt.Println("USAGE: ffcorpus pick sentences.txt [min] [max]")
	fmt.Println("USAGE: ffcorpus strip *.srt sentences.txt")
	fmt.Println("USAGE: ffcorpus tweet sentences.txt")
	os.Exit(0)
}
