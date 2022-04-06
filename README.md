# AirQualitySP project: air quality notifications for Spain

__UNDER CONSTRUCTION__

## What is AirQualitySP

## Project status

## Coding guidelines

This will provide a brief look at deploying cloud infrastructure with AWS Cloud
Development Kit (CDK) with Golang.

### The problem: Cloudformation has grown complex

At least a decade ago modern config management tools became a thing and were
teaching Sysadmins to abandon shell scripts and CLIs to manage infrastructure in
an imperative way in favor of declaring the desired state of the infrastructure
in a yaml (or so) document and letting the config management tool to figure out
how to get there. This way of defining system config and infrastructure
declaratively was praised as less error prone (guaranteed repeatable, etc) and
got adopted by then rising cloud providers as the standard way to get
infrastructure up, running and updated in the cloud. In combination with version
control the term “infrastructure as code” was coined. On AWS the platform’s
built-in IaC service and declaration syntax is called Cloudformation (Cfn).

AWS CDK (AWS Cloud Development Kit) is a polyglot framework and toolkit for
generating and deploying apps with one or more Cloudformation stacks from a
number of programming languages. Currently there are bindings for TypeScript,
Python, Java, .NET, and Go.

The support for Go is currently in “Developer Preview”, which means there can be
breaking API changes.

Typescript is CDK’s native language and plays a special role. CDK itself is
written in Typescript. Bindings for other languages get generated from
Typescript using an open source framework called jsii. (jsii has been
specifically developed by Amazon for CDK).

## Instructions

### Deploy cloud infrastructure on AWS

Create a new stack from scratch:

```bash
  cd feed-lambda && cdk init --language=go
```

Check the cloudformation template generated:

```bash
  cdk synth
```

(If it’s the first time that we deploy via CDK to an AWS account and region, we
need to run cdk bootstrap before we can run cdk deploy. This creates some
preparation in the AWS account:

```bash
  cdk bootstrap && cdk deploy
```

## Contributions

Would you like to provide any feedback?, please open up an Issue. I appreciate
feedback and comments, although please keep in mind the project is incomplete,
and I'm doing my best to keep it up to date.

## References

[poweruser](https://poweruser.blog/aws-cdk-with-go-part1-4075eeeceaad)

```go
package main

import (
 "github.com/aws/aws-cdk-go/awscdk/v2"
 "github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
 "github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
 "github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
 "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
 "github.com/aws/aws-sdk-go/aws"
 "github.com/aws/constructs-go/constructs/v10"
 "github.com/aws/jsii-runtime-go"
)

// StackProps -
type StackProps struct {
 awscdk.StackProps
}

func main() {
 app := awscdk.NewApp(nil)

 NewStack(app, "AirQualitySpStack", &StackProps{
  awscdk.StackProps{
   Env: env(),
  },
 })

 app.Synth(nil)
}

var (
 dynamoTableName string = "AirQualityTable"
)

// NewStack - Defines new Cloudformation Stack
func NewStack(scope constructs.Construct, id string, props *StackProps) awscdk.Stack {
 var sprops awscdk.StackProps
 if props != nil {
  sprops = props.StackProps
 }
 stack := awscdk.NewStack(scope, &id, &sprops)

 lambdaFn := awslambda.NewFunction(stack, jsii.String("Singleton"),
  &awslambda.FunctionProps{
   FunctionName: jsii.String("FeedFunction"),
   Code:         awslambda.NewAssetCode(jsii.String("lambda"), nil),
   Description:  jsii.String("A lambda function that retreived the air quality from cities"),
   Handler:      jsii.String("main"),
   Timeout:      awscdk.Duration_Seconds(jsii.Number(300)),
   Runtime:      awslambda.Runtime_GO_1_X(),
   Environment: &map[string]*string{
    "DYNAMODB_TABLE_NAME": jsii.String(dynamoTableName),
   },
  },
 )

 // See https://docs.aws.amazon.com/lambda/latest/dg/tutorial-scheduled-events-schedule-expressions.html
 rule := awsevents.NewRule(stack, jsii.String("Rule"), &awsevents.RuleProps{
  Schedule: awsevents.Schedule_Expression(jsii.String("cron(0/1 * * * ? *)")),
 })

 rule.AddTarget(awseventstargets.NewLambdaFunction(lambdaFn, nil))

 awsdynamodb.NewTable(stack, jsii.String("AirQualityTable"),
  &awsdynamodb.TableProps{
   TableName: aws.String(dynamoTableName),
   PartitionKey: &awsdynamodb.Attribute{
    Name: jsii.String("id"),
    Type: awsdynamodb.AttributeType_STRING,
   },
   BillingMode: &awsdynamodb.BillingMode_PAY_PER_REQUEST,
  },
 )

 return stack
}

func env() *awscdk.Environment {
 return nil
}
```
