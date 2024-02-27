# Application Overview
## Description
This application, developed using Golang v1.21, serves as a transmitter service that interacts with Artemis and MongoDB.

## Features
- Sends WhatsApp Message from queue on Artemis message broker.
- Stores Outbound message and Delivery Report in MongoDB database and Artemis.

# Deployment
## Installation
1. Clone this repository to your machine.
2. Navigate to the project directory.
3. Change configuration file on bin/config.yaml.

## Build Docker images
```bash
cd bin/
docker build -t transmitter-artemis .
```

## Run Docker Compose
```bash
docker compose up -d --build
```
This will start the application container using the Docker Compose file provided.

# Manual
## Step 1: Set Up Meta Developer Account
1. Make sure you already have a Meta developer account.
2. Ensure that you have set up the WhatsApp API, including the phone-number-id, and have a valid recipient number and token..

## Step 2: Set Up MongoDB
1. Ensure MongoDB is installed and running on your system.
2. Open a terminal and connect to your MongoDB instance.
3. Create a new database named transmitterdb:
```bash
use transmitterdb
```
4. Create a collection named **client-info** and insert a sample document:
```bash
db.client-info.insertOne({
  "token": "your-token",
  "client_name": "xxxxx",
  "phone_number_id": "your-phone-id",
  "wa_host": "https://graph.facebook.com/"
})
```
## Step 3: Run Application
1. Navigate to the directory containing the transmitter-artemis binary and Dockerfile.
2. Build the Docker image
3. Run the Docker container
4. Check the application logs to ensure successful connections to MongoDB and Artemis
5. Verify that the application has subscribed to the appropriate queue (xxxxx-msg-queue) for the client with name *xxxxx*.

## Step 4: Send Test Message
1. Send a test message to the **xxxxx-msg-queue queue**. Example message format:
```plain
message_id=1&to=6283872750005&type=text&text[preview_url]=false&text[body]=Hello world !!
```
2. Check success log on application.
3. Receive WhatsApp message from Meta.