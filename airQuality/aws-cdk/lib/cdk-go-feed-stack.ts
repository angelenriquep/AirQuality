import { Table, AttributeType, BillingMode, StreamViewType, } from 'aws-cdk-lib/aws-dynamodb';
import { AssetCode, Function, Runtime, StartingPosition } from "aws-cdk-lib/aws-lambda";
import { App, Duration, RemovalPolicy, Stack, StackProps } from 'aws-cdk-lib';
import { DynamoEventSource } from 'aws-cdk-lib/aws-lambda-event-sources';
import { LambdaSubscription } from 'aws-cdk-lib/aws-sns-subscriptions'
import { LambdaFunction } from 'aws-cdk-lib/aws-events-targets';
import { Schedule, Rule } from 'aws-cdk-lib/aws-events';
import { Topic } from 'aws-cdk-lib/aws-sns';
import { Construct } from 'constructs';

const wakiToken = process.env.ENV_WAKI_TOKEN_KEY || ''

const DYNAMODB_TABLE_NAME = process.env.ENV_DYNAMODB_TABLE_NAME || 'airQualityCities'
const LAMBDA_WAQI_FUNCTION = 'airQualityWaqiFeed'
const LAMBDA_FEED_PROCESSOR_FUNCTION = 'airQualityFeedProcessor'
const LAMBDA_FEED_TWITTER_PUBLISHER = 'airQualityTwitterPublisher'
const SNS_TOPIC = 'airQualitySNSTopic'
const LAMBDA_CRON_RULE_EACH_5_MINUTES = 'ruleCronEach5Minutes'

export class FeedApp extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    // DynamoDB
    const table = new Table(this, DYNAMODB_TABLE_NAME, {
      tableName: DYNAMODB_TABLE_NAME,
      partitionKey: {
        name: 'ID',
        type: AttributeType.NUMBER
      },
      sortKey: {
        name: 'CityName',
        type: AttributeType.STRING
      },
      billingMode: BillingMode.PAY_PER_REQUEST,
      removalPolicy: RemovalPolicy.DESTROY,
      stream: StreamViewType.NEW_AND_OLD_IMAGES,
    });

    // WAQI Feed Lambda
    const lambdaFeedFunctionWaqi = new Function(this, LAMBDA_WAQI_FUNCTION, {
      code: new AssetCode('../feed-lambda-waqi'),
      handler: 'main',
      runtime: Runtime.GO_1_X,
      timeout: Duration.seconds(300),
      memorySize: 256,
      environment: {
        "DYNAMODB_TABLE_NAME": DYNAMODB_TABLE_NAME,
        "WAKI_TOKEN": wakiToken,
      }
    });

    table.grantReadWriteData(lambdaFeedFunctionWaqi)

    const rule = new Rule(this, LAMBDA_CRON_RULE_EACH_5_MINUTES, {
      schedule: Schedule.expression('cron(0/5 * ? * * *)')
    });

    rule.addTarget(new LambdaFunction(lambdaFeedFunctionWaqi));

    // Lambda Feed processor
    const lambdaFunctionFeedProcessor = new Function(this, LAMBDA_FEED_PROCESSOR_FUNCTION, {
      code: new AssetCode('../feed-lambda-processor'),
      handler: 'main',
      runtime: Runtime.GO_1_X,
      timeout: Duration.seconds(300),
      memorySize: 256,
      environment: {
        "DYNAMODB_TABLE_NAME": DYNAMODB_TABLE_NAME,
        "SNS_TOPIC": SNS_TOPIC
      }
    });

    // Add an SQS Event Source from the DynamoDB Table to the Lambda Function
    lambdaFunctionFeedProcessor.addEventSource(new DynamoEventSource(table, {
      startingPosition: StartingPosition.LATEST,
    }));

    table.grantStreamRead(lambdaFunctionFeedProcessor);

    // Create sns topic
    const topic = new Topic(this, SNS_TOPIC, {
      displayName: 'New Pollution Data',
    });

    // Lambda twitter publisher
    const twitterLambda = new Function(this, LAMBDA_FEED_TWITTER_PUBLISHER, {
      code: new AssetCode('../feed-lambda-processor'),
      handler: 'main',
      runtime: Runtime.GO_1_X,
      timeout: Duration.seconds(300),
      memorySize: 256,
      environment: {}
    });

    topic.addSubscription(new LambdaSubscription(twitterLambda));


  }
}
const app = new App();
new FeedApp(app, `feed-app`);
app.synth();