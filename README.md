# AirQualitySP project: air quality and pollution notifications for Spain

## Coding guidelines

This will provide a brief look at how to deploy cloud infrastructure with AWS Cloud
Development Kit (CDK) with Golang.

## Instructions

A __makefile__ is provided in order to provide better deployment experience.

Ensure you have node.js and the AWS cdk installed, also make sure you have go
lang installed. As a reference the app has been developed using node 17 and go 1.18.

This Repo is not intended to teach about the AWS resources or how to create
resources using the AWS CDK, its only a prof of concept.

In order to make this app up and running you will need to preconfigure a Twitter
dev app with Auth1 and populate the environmental variables acordingly.

An example `.env` file is provided as example but __you will need to create a
`.env.dev` in order to geet this app up and running in the root directory.__

Env variables should be populated:

- ENV_WAQI_TOKEN_KEY (You will need to create an account in WAQI to get it)
- ENV_CITY_LIST (Alist of comma separated City Names)
- ENV_DYNAMODB_TABLE_NAME (The table name)
- ENV_TWITTER_CONSUMER_KEY (Auth1 twitter provided)
- ENV_TWITTER_CONSUMER_SECRET (Auth1 twitter provided)
- ENV_TWITTER_ACCESS_TOKEN_KEY (Auth1 twitter provided)
- ENV_TWITTER_ACCESS_TOKEN_SECRET (Auth1 twitter provided)

Check the available locations on <https://waqi.info>

The App is preconfigured to be executed every 5 minutes.

### Deploy cloud infrastructure on AWS

Generate the required Go binaries:

```bash
  make build
```

Deploy the App into AWS:

```bash
  make deploy
```

## Contributions

Would you like to provide any feedback?, please open up an Issue. I appreciate
feedback and comments, although please keep in mind the project is incomplete,
and I'm doing my best to keep it up to date.

## Demo

[DEMO](https://twitter.com/ngelEnr27558455)
