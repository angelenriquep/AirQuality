package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sns"
)

// DynamoDBStreamRecord - Shape for dynamo strem record.
type DynamoDBStreamRecord struct {
	ApproximateCreationDateTime events.SecondsEpochTime             `json:"ApproximateCreationDateTime,omitempty"`
	Keys                        map[string]*dynamodb.AttributeValue `json:"Keys,omitempty"`
	NewImage                    map[string]*dynamodb.AttributeValue `json:"NewImage,omitempty"`
	OldImage                    map[string]*dynamodb.AttributeValue `json:"OldImage,omitempty"`
	SequenceNumber              string                              `json:"SequenceNumber"`
	SizeBytes                   int64                               `json:"SizeBytes"`
	StreamViewType              string                              `json:"StreamViewType"`
}

// DynamoDBEventRecord - Shape for dynamo event record.
type DynamoDBEventRecord struct {
	AWSRegion      string                       `json:"awsRegion"`
	Change         DynamoDBStreamRecord         `json:"dynamodb"`
	EventID        string                       `json:"eventID"`
	EventName      string                       `json:"eventName"`
	EventSource    string                       `json:"eventSource"`
	EventVersion   string                       `json:"eventVersion"`
	EventSourceArn string                       `json:"eventSourceARN"`
	UserIdentity   *events.DynamoDBUserIdentity `json:"userIdentity,omitempty"`
}

// DynamoDBEvent - A list of Dynamodb records.
type DynamoDBEvent struct {
	Records []DynamoDBEventRecord `json:"Records"`
}

// City - Struct to store a city pollution data.
type City struct {
	CityName  string `json:"cityName"`
	CreatedAt string `json:"createdAt"`
	ID        int    `json:"id"`
	Pollution int    `json:"pollution"`
}

// Message - SNS message shape.
type Message struct {
	Default string `json:"default"`
}

// Handles a Dynamodb event from a DynamoDB stream and sends the data to a SNS
// topic.
func lambdaHandler(event DynamoDBEvent) error {
	var snsArn = os.Getenv("SNS_TOPIC_ARN")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sns.New(sess)

	for _, record := range event.Records {

		change := record.Change
		newImage := change.NewImage

		var item City
		// Converts from dynamodb attributes to Go types
		err := dynamodbattribute.UnmarshalMap(newImage, &item)
		if err != nil {
			log.Println(err.Error())
			return err
		}

		itemStr, _ := json.Marshal(item)

		message := Message{
			Default: string(itemStr),
		}

		messageBytes, _ := json.Marshal(message)
		messageStr := string(messageBytes)

		result, err := svc.Publish(&sns.PublishInput{
			TopicArn:         aws.String(snsArn),
			Message:          aws.String(messageStr),
			MessageStructure: aws.String("json"),
		})

		if err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}

		log.Println(*result)
	}

	return nil
}

func main() {
	lambda.Start(lambdaHandler)
}
