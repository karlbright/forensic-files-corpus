package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	forensicfilescorpus "github.com/karlbright/forensic-files-corpus"
)

func Handler() (string, error) {
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))
	client := twitter.NewClient(config.Client(oauth1.NoContext, token))

	rand.Seed(time.Now().UnixNano())

	pick, err := forensicfilescorpus.PickFromFile("sentences.txt", 0, 280)
	if err != nil {
		return "", err
	}

	tweet, _, err := client.Statuses.Update(pick, nil)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://twitter.com/%v/status/%v", tweet.User.ScreenName, tweet.IDStr)
	return url, nil
}

func main() {
	lambda.Start(Handler)
}
