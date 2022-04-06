# TODO:
# Check if node is installed in machibe
# Check for env vars
# Check if go is installed to build
# npm i?

install-aws-cdk:
	npm install -g aws-cdk

build:
	cd ./airQuality/feed-lambda-waqi && go build
	cd ..
	cd ./airQuality/feed-lambda-processor && go build

deploy:
	cd ./airQuality/aws-cdk && cdk deploy

destroy:
	cd ./airQuality/aws-cdk && cdk destroy

synth:
	cd ./airQuality/aws-cdk && cdk synth

generate-yaml-cloudformation:
	cd ./airQuality/aws-cdk && cdk synth > cloudFormation.yaml
