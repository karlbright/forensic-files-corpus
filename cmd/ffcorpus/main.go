package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	forensicfilescorpus "github.com/karlbright/forensic-files-corpus"
)

func main() {
	strip := flag.Bool("strip", false, "Strip sentences from a given set of subtitle files")
	pick := flag.Bool("pick", false, "Pick a random sentence from the given sentences sources")
	flag.Parse()

	if len(os.Args) == 1 {
		fmt.Println("USAGE: ffcorpus --strip *.srt sentences.txt")
		fmt.Println("USAGE: ffcorpus --pick sentences.txt [min?] [max?]")
		fmt.Println("USAGE: ffcorpus sentences.txt [min?] [max?]")
		os.Exit(1)
	}

	if *strip {
		if len(os.Args) < 4 {
			fmt.Println("USAGE: ffcorpus --strip *.srt sentences.txt")
			os.Exit(1)
		}

		dest, err := os.OpenFile(os.Args[len(os.Args)-1], os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			if err := dest.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		sentences := forensicfilescorpus.StripAll(os.Args[2 : len(os.Args)-2])
		for _, sentence := range sentences {
			dest.WriteString(sentence)
			dest.WriteString("\n")
		}

		os.Exit(0)
	}

	if *pick {
		if len(os.Args) < 3 {
			fmt.Println("USAGE: ffcorpus --pick sentences.txt [min?] [max?]")
			os.Exit(1)
		}

		src, err := os.Open(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			if err := src.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		var sentences []string
		scanner := bufio.NewScanner(src)
		for scanner.Scan() {
			sentences = append(sentences, scanner.Text())
		}

		min := -1
		max := -1
		rand.Seed(time.Now().UnixNano())

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

		sentence, err := forensicfilescorpus.Pick(sentences, min, max)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(sentence)
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Println("USAGE: ffcorpus sentences.txt [min?] [max?]")
	}

	src, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := src.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var sentences []string
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		sentences = append(sentences, scanner.Text())
	}

	min := 140
	max := 280
	rand.Seed(time.Now().UnixNano())

	if len(os.Args) == 3 {
		min, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
	}

	if len(os.Args) == 4 {
		min, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}

		max, err = strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
	}

	paragraph, err := forensicfilescorpus.Generate(sentences, min, max)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(paragraph)
}
