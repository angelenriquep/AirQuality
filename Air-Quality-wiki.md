# AirQualitySP high level system design

## Table of Contents

1. [AirQualitySP high level system design](#airQualitySP-high-level-system-design)

* 1.1 [Definitions and Acronyms](#definitions-and-acronyms)
* 1.2 [Abstract](#abstract)
* 1.3 [Goal/Objectives](#goal/objectives)
* 1.4 [Stakeholders](#stakeholders)
* 1.5 [Assumptions](#assumptions)
* 1.6 [Limitations & Unknowns](#limitations-&-unknowns)
* 1.7 [Supported use-cases](#supported-use-cases)
* 1.8 [Out of scope](#out-of-scope)

2. [Proposal](#out-of-scope)

* 2.1 [Architecture](#out-of-scope)
* 2.2 [Phase 1 description (covered in this document)](#out-of-scope)
* 2.3 [Phase 2(next steps)](#out-of-scope)
* 2.4 [Brief description of the problem](#out-of-scope)
* 2.5 [Cloudformation has grown complex](#out-of-scope)

3. [Costs](#costs)

* 3.1 [Monthly Cost (30 day) (Phase 1)](#monthly-cost-(30 day)-(Phase 1))

4. [References](#references)

## Definitions and Acronyms

Note: many of these terms require knowledge on AWS Services and Products, for
more information refer to <https://aws.amazon.com/>

* AQI: Air Quality Index
* AWS: Amazon Web Services
* API: application program interface
* CityFeed: A city feed represents the collection of airquality data reported by
  stations within a city (or in certain cases a Spanish county). For example, if
  Madrid City has "N" number of stations, Madrid's CityFeed would
  represent Madrid as a single city, so in other words, it represents all the
  Air quality data from these "N" stations.
* DTO: Data transfer object
* ETA: estimated time of arrival(used to describe how long something might take
  until it's done)
* Lambda: Aws solution to runcode in a serverless fashion(no servers required).
* SQS: Simple Queue Service (from AWS)
* SNS: Simple notification service(from AWS)
* Dynamo DB: A No-SQL database product from AWS
* Dynamo DB Streams: it's a notification mechanism that allows to send a
  DynamoDB item to any listening service in AWS, triggered when Dynamo DB
  detects changes on such item.
* MIT License: Permissive License by Massachusetts Institute of Technology used
  in open-source software.
* WAQI: The World Air Quality Index project
* City feed: represents the data feed from a supported city, includes one or
  more stations.

## Abstract

This project aims to use public data of air quality stations located in Spain,
and use this data to send notifications to a given user in a near-realtime
fashion. The system will be Cloud based and notifications may be published to
the public while providing corresponding attribution(s) to the data source(s).

## Goal/Objectives

Provide a suitable system architecture which would freely provide near-realtime
notifications regarding air quality in Spain. This will allow users to get
notified if air quality changes allowing them to knowingly take actions to
safeguard their health, for instance deciding not to go out if air quality is
not good, or doing outdoors activities when air quality is good. Provided AQI
will use the AQI US-EPA 2016 standard. This high level design does not inted to
dive into the specifics of implementation of the presented components.

Additionally, this project targets the open-source software development
community, to share the code of this project, as long as
guidance/practices/methodologies used in its development.

## Stakeholders

* Project owner: Angel Enrique.
* Users: anyone interested in using or actively using this system with the
  intention of receiving air quality notifications in a supported area (Spain).
* Open-source community: anyone insterested in using this software for learning
  purposes or those allowed by MIT License.

## Assumptions

* This is a __high level__ design, is not ment to hold specifics of implementation of
the presented components on this system.
* Each "feed" or "city feed" represents all the stations within a Spanish city
  or in certain cases a Spanish county.
* Adding support to a new city (becoming a "supported city feed") is treated as
  separate from this design, assume this is done manually.

* Data obtained from the World Air Quality Project programmatic API is
  considered the __source of truth__.
* Data retrieval frequency is by default __every 5 minutes__, it is subject to
  change.
* Notifications sent (ex. in the form of a Tweet) will be validated, won't send
  the same notification twice(only if there are changes).
* Data communication between components will be by using DTOs.
* For Rate limits at the "publisher" only Twitter's limits will be considered
  seeTwitters Limits. As of 4/21/2022, the Tweet limit is 900
  requests/15-minutes. More info check
  [Twitter-rate-limits](https://developer.twitter.com/en/docs/twitter-api/rate-limits)
* As per WAQI's terms ["The data can not be redistributed as cached or archived
  data"](https://aqicn.org/api/tos/), however the way system stores data is to
  detect changes, and does not redistribute cached or archived data, only new
  information is published.
* As per previously mentioned WAQI's terms, all messages/notifications that
  fail, will be discarded.
* Next steps (Phase 2) section attempts to reflect feedback, is not to be
  considered final, it may be updated continuously.

## Limitations & Unknowns

* WAQI's data seems presented as "real-time", however there's no confirmed
  update frequency for specific cities.
* There's no known push method to listen for data changes in WAQI's API, so we
  "pull" every 5 minutes via GET HTTP request.
* Twitter needs to approve a development account. ETA is not known for your use
  case, in our case it was inmediately.
* How does Twitter detects Spam. Should there be any considerations?

## Supported use-cases

* System is able to access AQI data from WAQI's programmatic API.
* Users can be notified when any station located in a determined city show
  changes in respect to air quality and will use AQI scale.
* Specific rules that determine wether or not a notification should be
  published, should prevent spam (Not publishing same notification more than
  once or avoid publishing too frequently).
* Users are notified via twitter (users should follow a determined Twitter
  account to access messages). [Project Account](https://twitter.com/ngelEnr27558455).
* Notified users can access a text message containing the AQI, Air Pollution
  Level (Good, Moderate, Unhealthy for sensitive groups, Unhealthy, Very
  Unhealthy and Hazardous) according to the [AQI scale](https://www.airnow.gov/aqi/aqi-basics/).
* System is able to publish text messages as tweets on a determined twitter
  account with AQI.
* Text messages delivered by the System must include attribution to WAQI and any
  other data source(such as, but not limited to organizations, institutions,
  public APIs, person/people, etc).
* In case of failure, message notifications must be discarded.

## Out of scope

* Any form of user data collection: no data will be collected in this project.
* Notifications messages specific to an individual "station" within a city.
* Notifications on areas that are not located in Spain.
* Search/Lookup functionality
* Email notifications are out of scope
* Automation of setup for adding support of cities.
* Error handling
* Retry mechanisms
* As system will be divided in multiple development phases, Cost calculation
  won't be updated from initial calculation, which represents a Worst Case
  Scenario.

# Proposal

## Architecture

The proposed system will use a __Serverless__ approach using AWS services.
Architecture will be developed in multiple phases. This document covers Phase 1,
although details for an upcomming are considered out of scope, notes for next
steps will be considered as Phase 2.

## Phase 1 description (covered in this document)

Relies on Cloudwatch events to trigger periodic calls to the data source, and
uses AWS Lambdas, SNS topics and DynamoDB streams to deliver the data to its
final destination. It also performs validations on it's first Lambda component,
the __feed-lambda-waqi__ which only stores information on DynamoDB if it meets the
criteria of specific set of rules. Such rules may change within the system's
implementation itself, the objective is to only notify on relevant, new
information(relevant changes in data), and avoid spam.

## Phase 2(next steps)

The following are features considered for upcoming development.

* Error handling flow: consider use of aws step functions or equivalent, to
  properly handle failures of the components within the system.
* SMS notification on failures: should failures occur, automatically notify the
  system's owner by using SMS.
* Retry strategies: upon failures, system may retry within a determined period
  of time (ex. 3 minutes), and may re-attempt to continue the flow of the
  system, if it is meets validation criteria(ex. Relevant new information, not
  spam).
* feed-lambda-waqi will write to an SQS queue instead of writing directly into
  DynamoDB. This will open the possibility of having a new Lambda which can
  (within 1 execution), process batches of messages available in the SQS queue.
  Also opens the possibility of parallel execution, having multiple lambdas
  picking up messages from the same SQS queue.
* __To be determined__: moving architecture to an Azure Cloud.

## Brief description of the problem

### Cloudformation has grown complex

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

## Data structures

Sample WAQI API /feed/:city response

```json
{
   "status":"ok",
   "data":{
      "aqi":42,
      "idx":6732,
      "attributions":[
         {
            "url":"http://www.euskadi.eus/gobierno-vasco/medio-ambiente/",
            "name":"Departamento de Medio Ambiente, Planificación Territorial y Vivienda · Gobierno Vasco",
            "logo":"Spain-GobiernoVasco.png"
         },
         {
            "url":"http://www.eea.europa.eu/themes/air/",
            "name":"European Environment Agency",
            "logo":"Europe-EEA.png"
         },
         {
            "url":"https://waqi.info/",
            "name":"World Air Quality Index Project"
         }
      ],
      "city":{
         "geo":[
            43.267505511797,
            -2.9351881103382
         ],
         "name":"Mazarredo, Bilbao, País Vasco, Spain",
         "url":"https://aqicn.org/city/spain/pais-vasco/bilbao/mazarredo",
         "location":""
      },
      "dominentpol":"pm25",
      "iaqi":{
         "co":{
            "v":0.1
         },
         "dew":{
            "v":9
         },
         "h":{
            "v":81
         },
         "no2":{
            "v":15.6
         },
         "o3":{
            "v":16.7
         },
         "p":{
            "v":1006
         },
         "pm10":{
            "v":19
         },
         "pm25":{
            "v":42
         },
         "so2":{
            "v":4.6
         },
         "t":{
            "v":12
         },
         "w":{
            "v":3
         },
         "wg":{
            "v":11.8
         }
      },
      "time":{
         "s":"2022-04-21 17:00:00",
         "tz":"+02:00",
         "v":1650560400,
         "iso":"2022-04-21T17:00:00+02:00"
      },
      "forecast":{
         
      },
      "debug":{
         "sync":"2022-04-22T03:21:26+09:00"
      }
   }
```

# Costs

## Monthly Cost (30 day) (Phase 1)

UNDER STIMATION

## References

* [poweruser](https://poweruser.blog/aws-cdk-with-go-part1-4075eeeceaad)
* [sns-example-publish](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/sns-example-publish.html)
* [waqi-api](https://waqi.info/)
* [official-aws-repo](https://github.com/aws/aws-lambda-go)
