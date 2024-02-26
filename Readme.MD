# Golang Microservices

This is a Golang Microservices playground where we can find the following:

## Broker Service
The broker service is a main service that receives external requests and forwards them to the specific microservices.
There are 4 types of communication.
- REST API. The classic way applications communicate by sending/receiving JSON.
- RabbitMQ. The broker service publishes some events that each service consume.
- RPC. We communicate via RPC for quicker results.
- gRPC. With the help of a proto file we communicate via Google's gRPC.

## Authentication Service
The authentication service has a PostgresDB where our users are stored.

## Listener Service
The listener service is a service that is here to consume our RabbitMQ queues.

## Logger Service
The logger service keeps logs of our system in a MongoDB database.

## Mail Service
The mail service is sending mails using mailhog. We can use plain text or templates.

## Extras & Comments
The project folder contains files needed to run/build our project.
After downloading the code, running `make up_build`, everything should run fine.