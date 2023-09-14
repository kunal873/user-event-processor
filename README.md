# User Event Processor

The User Event Processor is a Go application that allows you to publish user events to a Redis stream and process those events asynchronously. It provides an HTTP API endpoint for publishing events and a background process for reading and processing events from the Redis stream.

## Features

- Publish user events with a JSON payload.
- Asynchronously process user events with retry logic.
- Redis integration for event storage and communication.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Contributions](#contributions)

## Prerequisites

Before you start using this application, ensure you have the following installed:

- [Go](https://golang.org/dl/) (v1.21.0 or higher)
- [Docker](https://www.docker.com/get-started) (for running Redis in a container)

## Getting Started

1. **Clone the repository** to your local machine:

   ```shell
   git clone https://github.com/kunal873/user-event-processor.git
   cd user-event-processor
    ```
2. **Compose Redis Container** in a Docker container:

   ```shell
   docker-compose up redis
   ```
3. **Compose the publisher** in a new terminal window:

   ```shell
    docker-compose up pub
   ```
4. **Compose the subscriber** in a new terminal window:

   ```shell
    docker-compose up sub
   ```
 
## Configuration

You can configure the application by modifying the .env file in the project directory. Make sure to restart the application after making changes to the environment variables.

## Contributions
Contributions to this project are welcome. Feel free to open issues or create pull requests.