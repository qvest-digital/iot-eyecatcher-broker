# iot-eyecatcher-broker
The websocket message broker for the IoT eye catcher project.

Simple message broker without persistence written in golang.

## Build and run

To build the docker image call:
    
    docker build -t iot-eyecatcher-broker:latest .
    
To run the docker image call:

    docker run -e USERNAME={username] -e PASSWORD={password} iot-eyecatcher-broker:latest

## HTTP API

### GET /{topicName}
Get the last message of the given topic. 404 if topic does not exist.

### POST /{topicName}
Send an new message of the given topic. Creates topic if it does not exist.
You need to authenticate with HTTP Basic authentication.

## Websocket API 

### GET /

Once connected, send the following message to...

...subscribe to a topic: 

    {"operation":"subscribe", "topic":"insert-topic-name-here"}

...unsubscribe from a topic:

    {"operation":"unsubscribe", "topic":"insert-topic-name-here"}


