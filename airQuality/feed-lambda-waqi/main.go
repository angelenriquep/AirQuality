package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	dynaClient dynamodbiface.DynamoDBAPI
)

// City - Stores a city.
type City struct {
	CityName  string
	CreatedAt string
	ID        int
	Pollution int
}

// Response - Weather API shape.
type Response struct {
	Data struct {
		Aqi  int `json:"aqi"`
		Idx  int `json:"idx"`
		City struct {
			Name string `json:"name"`
		} `json:"city"`
		Time struct {
			Iso string `json:"iso"`
		} `json:"time"`
	} `json:"data"`
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, e events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var tableName = os.Getenv("DYNAMODB_TABLE_NAME")
	var wakiToken = os.Getenv("WAKI_TOKEN")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	resp, err := http.Get(fmt.Sprintf("https://api.waqi.info/feed/Madrid/?token=" + wakiToken))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return generateErrorResponse("Got error marshalling a city item: %s", err), nil
	}

	city := City{
		CityName:  result.Data.City.Name,
		CreatedAt: result.Data.Time.Iso,
		ID:        result.Data.Idx,
		Pollution: result.Data.Aqi,
	}

	av, err := dynamodbattribute.MarshalMap(city)
	if err != nil {
		return generateErrorResponse("Got error marshalling a city item: %s", err), nil
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 201,
		Body:       string("Successfully added: " + city.CityName),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

// Helper gunction to generate error responses.
func generateErrorResponse(error string, err error) events.APIGatewayV2HTTPResponse {
	log.Fatalf(error, err)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 500,
		Body:       string(error),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}
