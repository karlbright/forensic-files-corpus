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

func pick() {
	if len(os.Args) < 3 {
		fmt.Println("USAGE: ffcorpus pick sentences.txt [min] [max]")
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())
	var err error
	min := -1
	max := -1

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
	sentence, err := forensicfilescorpus.PickFromFile(path, min, max)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(sentence)
	os.Exit(0)
}
