# Identity Service

## Overview

This repository contains a **microservice** built using **Go** with the **Fiber** framework. It serves as the **identity service** for the **Beat Back End**, responsible for managing account creation, access token generation, and user profile management.

### Technologies
- **Go Version**: 1.23.2
- **Framework**: **Fiber** (Version: 2.52.5)
- **Docker**: Containerized microservice
- **.env**: Environment-specific configurations for the project

---

## Setup for Development

### Prerequisites

Before you can start developing with this project, ensure you have the following tools installed:

- **[asdf](https://asdf-vm.com/)**: Version manager to manage Go versions. If you don’t have it, feel free to use **asdf** to install Go.
- **Go**: Version 1.23.2
- **Docker**: For running the service within a Docker container

### Installing Dependencies

1. Clone the repository:

    ```sh
    git clone https://github.com/BeatEcoprove/beat_identity_server.git identity
    cd identity
    ```

2. Install the required version of Go (1.23.2) using **asdf**:

    ```sh
    asdf install golang 1.23.2
    ```

3. Install project dependencies:

    ```sh
    go mod tidy
    ```

---

## Project Helper

This project includes a **Makefile** for managing various tasks. You can use the `make` command to interact with the project:

```sh
make
```

Here is a list of available commands:
### Tools:

- **setup**: Download and install all necessary tools
- **install-goose**: Install the Goose migration tool
- **install-swag**: Install the Swag tool for API documentation generation

### Actions:

- **test**: Run tests
- **coverage**: Generate a coverage report
- **swagger**: Generate Swagger configuration files for API documentation

### Database Operations:

- **rollback**: Rollback the last migration
- **rebuild**: Rebuild the migrations
- **reset**: Rollback all migrations
- **status**: Show the current migration status
- **migrate**: Apply all pending migrations
- **create-migration**: Create a new migration with a user-provided name

## Configuration

The project relies on a `.env` file for environment-specific configuration. Here’s an example `.env` file:

```env
BEAT_IDENTITY_SERVER=3000

POSTGRES_DB=identity
POSTGRES_USER=beat
POSTGRES_PASSWORD=beat
POSTGRES_PORT=5432
POSTGRES_HOST=ecoprove

JWT_AUDIENCE=Beat
JWT_ISSUER=Beat
JWT_ACCESS_EXPIRED=10
JWT_REFRESH_EXPIRED=4
JWT_SECRET=ed395d0b3852a9917aedf1ec651bf92bf46ed418017982a312984f704395bcff

REDIS_HOST=redis
REDIS_UI_PORT=8000
REDIS_PORT=6379
REDIS_DB=0

RABBIT_MQ_HOST=broker
RABBIT_MQ_PORT=5672
RABBITMQ_DEFAULT_USER=beat
RABBITMQ_DEFAULT_PASS=beat
RABBITMQ_DEFAULT_VHOST=beat
RABBIT_MQ_EXCHANGE=email_exchange
RABBIT_MQ_QUEUE_MAIL=mails
RABBIT_MQ_ROUTING_KEY=send_email
RABBIT_MQ_UI_PORT=15672
```

> Make sure to set up these environment variables correctly before running the application.

## Database Migration

To set up the application, you need to create a PostgreSQL database instance and run the **migrations**.

### Run Database Migrations

Use the following `make` command to apply the migrations:

```sh
make migrate
```

Example output:

```sh
Applying migrations...
2024/12/02 01:46:10 OK   20241113171834_create_auth_table.sql (12.5ms)
2024/12/02 01:46:10 OK   20241113171952_create_profile_table.sql (5.78ms)
2024/12/02 01:46:10 goose: successfully migrated database to version: 20241113171952
```

This command will apply all pending migrations to your database.

## Docker Setup

To run the service in a Docker container, follow these steps:

### Build the Docker Image

To build the Docker image, run:

```sh
docker build --build-arg GO_VERSION=1.23.2 --build-arg UID=10001 -t identity-service -f ./docker/api/Dockerfile .
```

This command will build the Docker image for the service with Go version `1.23.2` and a custom user ID (`10001`).

### Run the Docker Container

Once the image is built, you can run the micro service in a Docker container:

```sh
docker run --env-file .env -p 3000:3000 identity-service
```

## Additional Information


#### HTTP Endpoints Documentation

All available API endpoints are documented in `swagger` located in http://localhost:3000/swagger/index.html.

Documentation:

- [`http://localhost:3000/api/v2/auth/sign-up`](https://github.com/BeatEcoprove/beat_identity_server/wiki/Sign-Up-Endpoint)
- [`http://localhost:3000/api/v2/auth/profile`](https://github.com/BeatEcoprove/beat_identity_server/wiki/Attach-Profile-Endpoint)
- [`http://localhost:3000//api/v2/check-field?email=<email>`](https://github.com/BeatEcoprove/beat_identity_server/wiki/Check-Filed-Endpoint)
- [`http://localhost:3000/api/v2/auth/forgot-password`](https://github.com/BeatEcoprove/beat_identity_server/wiki/Forgot-Password-Endpoint)
- [`http://localhost:3000/api/v2/auth/login`](https://github.com/BeatEcoprove/beat_identity_server/wiki/Login-Endpoint)
- [`http://localhost:3000/api/v2/auth/refresh-token?param=<profile-id>`](https://github.com/BeatEcoprove/beat_identity_server/wiki/Refresh-Tokens-Endpoint)
- [`http://localhost:3000/api/v2/auth/reset-password`](https://github.com/BeatEcoprove/beat_identity_server/wiki/Reset-Password-Endpoint)
- [`http://localhost:3000/api/v2/auth/token`](https://github.com/BeatEcoprove/beat_identity_server/wiki/Token-Endpoint)

Plus, follow this instructions to view the full documentation of each endpoint.