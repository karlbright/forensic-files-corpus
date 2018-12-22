package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	forensicfilescorpus "github.com/karlbright/forensic-files-corpus"
)

func generate() {
	if len(os.Args) < 3 {
		fmt.Println("USAGE: ffcorpus generate sentences.txt [min] [max]")
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())
	var err error
	min := 140
	max := 280

	if len(os.Args) == 4 {
		min, err = strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
	}

	if len(os.Args) == 5 {
		min, err = strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}

		max, err = strconv.Atoi(os.Args[4])
		if err != nil {
			log.Fatal(err)
		}
	}

	path := os.Args[2]
	paragraph, err := forensicfilescorpus.GenerateFromFile(path, min, max)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(paragraph)
	os.Exit(0)
}
