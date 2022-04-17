# Change this line acording to your env variables file.
ENV := .env.dev

# Environment variables for project
include $(ENV)

export

install-aws-cdk:
	npm install -g aws-cdk

build:
	cd ./airQuality/feed-lambda-waqi && go build
	cd ..
	cd ./airQuality/feed-lambda-processor && go build
	cd ..
	cd ./airQuality/feed-lambda-twitter-processor && go build

deploy:
	cd ./airQuality/aws-cdk && cdk deploy

destroy:
	cd ./airQuality/aws-cdk && cdk destroy

synth:
	cd ./airQuality/aws-cdk && cdk synth

generate-yaml-cloudformation:
	cd ./airQuality/aws-cdk && cdk synth > cloudFormation.yaml
