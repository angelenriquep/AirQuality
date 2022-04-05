import * as cdk from '@aws-cdk/core';
import * as dynamodb from '@aws-cdk/aws-dynamodb';
import * as lambda from '@aws-cdk/aws-lambda';
import { DynamoEventSource } from '@aws-cdk/aws-lambda-event-sources';

export class LambdaCronStack extends cdk.Stack {
  constructor(app: cdk.App, id: string) {
    super(app, id);

    const table = new dynamodb.Table(this, 'Table', {
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      partitionKey: { name: 'id', type: dynamodb.AttributeType.STRING },
      stream: dynamodb.StreamViewType.NEW_AND_OLD_IMAGES,
      tableName: 'Table',
    });

    const lambdaFunction = new lambda.Function(this, 'Function', {
      code: lambda.Code.fromAsset('src'),
      handler: 'index.handler',
      functionName: 'TableStreamHandler',
      runtime: lambda.Runtime.NODEJS_12_X,
    });

    lambdaFunction.addEventSource(new DynamoEventSource(table, {
      startingPosition: lambda.StartingPosition.LATEST,
    }));
  }
}

const app = new cdk.App();
new LambdaCronStack(app, 'LambdaCronExample');
app.synth();
