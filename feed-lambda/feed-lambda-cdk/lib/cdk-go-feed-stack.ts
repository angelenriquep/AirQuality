import { Table, AttributeType, BillingMode, StreamViewType } from 'aws-cdk-lib/aws-dynamodb';
import { App, Duration, RemovalPolicy, Stack, StackProps } from 'aws-cdk-lib';
import { AssetCode, Function, Runtime } from "aws-cdk-lib/aws-lambda";
import targets = require('aws-cdk-lib/aws-events-targets');
import events = require('aws-cdk-lib/aws-events');
import { Construct } from 'constructs';

// import config from './config.json';
const tableName = 'AirQualityTable'
const wakiToken = process.env.ENV_WAKI_TOKEN_KEY || ''

export class FeedApp extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const table = new Table(this, tableName, {
      tableName: tableName,
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


    const lambdaFeedFunction = new Function(this, 'AirQualityFeed', {
      code: new AssetCode('./lambda'),
      handler: 'main',
      runtime: Runtime.GO_1_X,
      timeout: Duration.seconds(300),
      memorySize: 256,
      environment: {
        "DYNAMODB_TABLE_NAME": tableName,
        "WAKI_TOKEN": wakiToken,
      }
    });

    table.grantReadWriteData(lambdaFeedFunction)

    const rule = new events.Rule(this, 'LambdaFeedEach5Minutes', {
      schedule: events.Schedule.expression('cron(0/5 * ? * * *)')
    });

    rule.addTarget(new targets.LambdaFunction(lambdaFeedFunction));
  }
}
const app = new App();
new FeedApp(app, `feed-app`);
app.synth();