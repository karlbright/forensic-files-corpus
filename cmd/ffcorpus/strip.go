package main

import (
	"fmt"
	"log"
	"os"

	forensicfilescorpus "github.com/karlbright/forensic-files-corpus"
)

func strip() {
	if len(os.Args) < 4 {
		fmt.Println("USAGE: ffcorpus strip *.srt sentences.txt")
		os.Exit(1)
	}

	paths := os.Args[2 : len(os.Args)-2]
	output := os.Args[len(os.Args)-1]

	if err := forensicfilescorpus.StripAllToFile(paths, output); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
