package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// @see https://github.com/aws/aws-lambda-go/blob/main/events/README_DynamoDB.md
func handleRequest(ctx context.Context, e events.DynamoDBEvent) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sns.New(sess)

	for _, record := range e.Records {
		fmt.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)

		// Print new values for attributes of type String
		for name, value := range record.Change.NewImage {
			if value.DataType() == events.DataTypeString {
			}

			fmt.Printf("Attribute name: %s, value: %s\n", name, value.String())

			// result, err := svc.Publish(&sns.PublishInput{
			// 	Message:  msgPtr,
			// 	TopicArn: topicPtr,
			// })
			// if err != nil {
			// 	fmt.Println(err.Error())
			// 	os.Exit(1)
			// }

			// fmt.Println(*result.MessageId)
		}
	}
}

func main() {
	lambda.Start(handleRequest)
}
