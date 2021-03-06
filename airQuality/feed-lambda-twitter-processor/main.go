package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

const (
	aqiIndexGood                  string = "Good"
	aqiIndexModerate              string = "Moderate"
	aqiIndexUnhealthyForSensitive string = "Unhealthy for Sensitive Groups"
	aqiIndexUnhealthy             string = "Unhealthy"
	aqiIndexVeryUnhealthy         string = "Very Unhealthy"
	aqiIndexHazardous             string = "Hazardous"
)

// Credentials stores all of our access/consumer tokens and secret keys needed
// for authentication against the twitter REST API.
type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// City - Stores a city.
type City struct {
	CityName  string
	CreatedAt string
	ID        int
	Pollution int
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, snsEvent events.SNSEvent) {
	creds := Credentials{
		AccessToken:       os.Getenv("TWITTER_ACCESS_TOKEN_KEY"),
		AccessTokenSecret: os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("TWITTER_CONSUMER_SECRET"),
	}

	client, err := getClient(&creds)
	if err != nil {
		log.Println("Error getting Twitter Client")
		log.Println(err)
	}

	cityStatuses := make([]string, 0, len(snsEvent.Records))

	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		fmt.Printf("[%s %s] Message = %s \n", record.EventSource, snsRecord.Timestamp, snsRecord.Message)

		var city City
		if err := json.Unmarshal([]byte(snsRecord.Message), &city); err != nil {
			log.Println(err)
			return
		}

		var strPollution string

		switch level := city.Pollution; {
		case level >= 0 && level < 51:
			strPollution = aqiIndexGood
		case level >= 51 && level <= 100:
			strPollution = aqiIndexModerate
		case level >= 101 && level <= 150:
			strPollution = aqiIndexUnhealthyForSensitive
		case level >= 151 && level <= 200:
			strPollution = aqiIndexUnhealthy
		case level >= 201 && level <= 300:
			strPollution = aqiIndexVeryUnhealthy
		case level >= 301:
			strPollution = aqiIndexHazardous
		default:
			strPollution = "Not Defined"
		}

		message := fmt.Sprintf(
			"City: %s\nAir Quality: %s \nPollution level Idx: %d \nSource: https://waqi.info/",
			city.CityName,
			strPollution,
			city.Pollution,
		)

		cityStatuses = append(cityStatuses, message)
	}

	strCities := strings.Join(cityStatuses, "\n")

	_, resp, err := client.Statuses.Update(strCities, nil)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v\n", resp)
}

// getClient helper function that will return a twitter client that we can
// subsequently use to send tweets, or to stream new tweets
func getClient(creds *Credentials) (*twitter.Client, error) {
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	log.Printf("User's ACCOUNT:\n%+v\n", user)
	return client, nil
}
