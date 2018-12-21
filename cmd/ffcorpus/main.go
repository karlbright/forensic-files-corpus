package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/Shopify/ejson"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	forensicfilescorpus "github.com/karlbright/forensic-files-corpus"
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
	case "tweet":
		tweet()
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

func strip() {
	if len(os.Args) < 4 {
		fmt.Println("USAGE: ffcorpus strip *.srt sentences.txt")
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

func pick() {
	if len(os.Args) < 3 {
		fmt.Println("USAGE: ffcorpus pick sentences.txt [min] [max]")
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

func generate() {
	if len(os.Args) < 3 {
		fmt.Println("USAGE: ffcorpus generate sentences.txt [min] [max]")
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

	min := 140
	max := 280
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

	paragraph, err := forensicfilescorpus.Generate(sentences, min, max)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(paragraph)
	os.Exit(0)
}

func tweet() {
	if len(os.Args) < 3 {
		fmt.Println("USAGE: ffcorpus tweet sentences.txt")
	}

	var secrets map[string]string
	decrypted, err := ejson.DecryptFile("secrets.ejson", "", os.Getenv("SECRETS_PRIVATE_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(decrypted, &secrets); err != nil {
		log.Fatal(err)
	}

	config := oauth1.NewConfig(secrets["consumer_key"], secrets["consumer_secret"])
	token := oauth1.NewToken(secrets["access_token"], secrets["access_secret"])
	client := twitter.NewClient(config.Client(oauth1.NoContext, token))

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
	rand.Seed(time.Now().UnixNano())

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		sentences = append(sentences, scanner.Text())
	}

	pick, err := forensicfilescorpus.Pick(sentences, 0, 280)
	if err != nil {
		log.Fatal(err)
	}

	tweet, _, err := client.Statuses.Update(pick, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IDStr)

	os.Exit(1)
}
