# 🔐 Beat Identity Service

> A modern, secure, and scalable authentication microservice built with Go and Fiber.

## 📋 Project Overview

The Beat Identity Service is a production-ready authentication and authorization microservice designed for the Beat ecosystem. It provides comprehensive identity management including user registration, OAuth2-style token flows, JWT-based authentication, password recovery, and role-based access control.

Built with Clean Architecture principles, this service integrates seamlessly with event-driven systems using Kafka for asynchronous operations, Redis for distributed caching, and PostgreSQL for persistent storage. The service generates RS256-signed JWT tokens with JWKS support, making it ideal for microservice architectures requiring secure, stateless authentication.

### ✨ Key Features
- 🎫 OAuth2-style authentication with access and refresh tokens
- 🔑 RS256 JWT signing with Public Key Infrastructure (PKI)
- 👥 Multi-profile support with role-based permissions
- 📡 Event-driven architecture (Kafka integration)
- 🔄 Password recovery flows with secure code generation
- 📚 Interactive Swagger API documentation
- ⚡ Distributed token management with Redis
- 🗃️ Database migrations with version control

### 🛠️ Technologies
- **Go Version**: 1.23.2
- **Framework**: Fiber (v2.52.5) - Fast, Express-inspired web framework
- **Database**: PostgreSQL with GORM
- **Cache**: Redis for token and session management
- **Message Broker**: Kafka for event streaming
- **Security**: RS256 JWT, JWKS, bcrypt password hashing
- **API Documentation**: Swagger/OpenAPI with interactive UI
- **Migrations**: Goose for database version control
- **Docker**: Fully containerized for easy deployment

---

## 🚀 Setup for Development

### 📦 Prerequisites

Before you can start developing with this project, ensure you have the following tools installed:

- **[asdf](https://asdf-vm.com/)**: Version manager to manage Go versions. If you don’t have it, feel free to use **asdf** to install Go.
- **Go**: Version 1.23.2
- **Docker**: For running the service within a Docker container

### 💻 Installing Dependencies

1. **Clone the repository:**

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

## 🔧 Project Helper

This project includes a **Makefile** for managing various tasks. You can use the `make` command to interact with the project:

```sh
make
```

Here is a list of available commands:
### 🛠️ Tools:

- **setup**: Download and install all necessary tools
- **install-goose**: Install the Goose migration tool
- **install-swag**: Install the Swag tool for API documentation generation

### ⚡ Actions:

- **test**: Run tests
- **coverage**: Generate a coverage report
- **swagger**: Generate Swagger configuration files for API documentation

### 🗄️ Database Operations:

- **rollback**: Rollback the last migration
- **rebuild**: Rebuild the migrations
- **reset**: Rollback all migrations
- **status**: Show the current migration status
- **migrate**: Apply all pending migrations
- **create-migration**: Create a new migration with a user-provided name

## ⚙️ Configuration

The project relies on a `.env` file for environment-specific configuration. Here's an example `.env` file:

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

## 🗃️ Database Migration

To set up the application, you need to create a PostgreSQL database instance and run the **migrations**.

### 🚀 Run Database Migrations

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

## 🐳 Docker Setup

To run the service in a Docker container, follow these steps:

### 🏗️ Build the Docker Image

To build the Docker image, run:

```sh
docker build --build-arg GO_VERSION=1.23.2 --build-arg UID=10001 -t identity-service -f ./docker/api/Dockerfile .
```

This command will build the Docker image for the service with Go version `1.23.2` and a custom user ID (`10001`).

### ▶️ Run the Docker Container

Once the image is built, you can run the micro service in a Docker container:

```sh
docker run --env-file .env -p 3000:3000 identity-service
```

## 🏗️ Project Architecture

This service follows **Clean Architecture** principles, ensuring separation of concerns, testability, and maintainability:

### 📊 Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                     HTTP Layer (Fiber)                      │
│                  Swagger UI | Endpoints                     │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                    Controllers Layer                        │
│            Route handlers & Request validation              │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                   Use Cases Layer                           │
│        Business Logic | SignUp | Login | Refresh            │
│        ForgotPassword | ResetPassword | Profiles            │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│              Repositories & Services                        │
│   Auth Repo | Profile Repo | Token Service | Email Service │
└────┬─────────────────┬─────────────────┬────────────────────┘
     │                 │                 │
┌────▼─────┐    ┌──────▼──────┐    ┌────▼──────┐
│PostgreSQL│    │    Redis    │    │   Kafka   │
│  (GORM)  │    │  (Caching)  │    │ (Events)  │
└──────────┘    └─────────────┘    └───────────┘
```

### 📁 Directory Structure

- **`cmd/identity-service/`** 🚀 - Application entry point and initialization
  - PKI/JWKS generation
  - Server lifecycle management

- **`internal/`** 🔒 - Core business logic (private to this service)
  - **`usecases/`** 💼 - Business use cases (SignUp, Login, RefreshTokens, etc.)
  - **`repositories/`** 🗄️ - Data access layer (Auth, Profile, MemberChat)
  - **`adapters/`** 🔌 - External integrations (HTTP server, Kafka, Redis, Database)
  - **`middlewares/`** 🛡️ - HTTP middlewares (Authorization, JWT validation)
  - **`domain/`** 🎯 - Domain models and events
    - **`events/`** 📨 - Event definitions (UserCreated, GroupCreated, etc.)
    - **`handlers/`** 🎬 - Event handlers for Kafka consumers

- **`pkg/`** 📦 - Shared utilities and contracts (reusable across services)
  - **`services/`** ⚙️ - Core services (JWT, PKI, Password, Email, AES)
  - **`contracts/`** 📄 - Request/Response DTOs
  - **`shared/`** 🤝 - Validators, error handlers, base controllers

- **`config/`** ⚙️ - Configuration management and migrations
- **`migrations/`** 🗃️ - Goose database migration files
- **`docs/`** 📚 - Auto-generated Swagger documentation

### 🔑 Key Components

**🔐 Authentication Flow:**
1. User signs up → Account created → JWT tokens issued
2. User logs in → Credentials validated → Access + Refresh tokens returned
3. Access token expires → Client uses refresh token → New access token issued

**📡 Event-Driven Integration:**
- **Produces:** `user_created`, `email_queue` events via Kafka
- **Consumes:** `group_created`, `invite_accepted` events to update permissions

**🛡️ Security:**
- 🔐 RS256 asymmetric JWT signing
- 🔑 JWKS endpoint (`/.well-known/jwks.json`) for public key distribution
- 🔒 Bcrypt password hashing
- 👮 Scoped permissions for group-based access control

## 📚 API Documentation

All available API endpoints are documented with **Swagger UI**:

**http://localhost:3000/swagger/index.html**

The interactive documentation provides:
- 📖 Complete endpoint reference with request/response schemas
- 🔐 Authentication requirements and examples
- 🧪 Try-it-out functionality to test endpoints directly
- ⚠️ Error response details and status codes

### 🎯 Core Endpoints

| Endpoint | Description |
|----------|-------------|
| `/api/v1/auth/sign-up` | Register a new user account |
| `/api/v1/auth/token` | OAuth2-style token endpoint (login or refresh) |
| `/api/v1/auth/forgot-password` | Request password reset code via email |
| `/api/v1/auth/reset-password` | Reset password with verification code |
| `/api/v1/auth/profiles` | Attach profile to authenticated account |
| `/api/v1/auth/profiles/me` | Get current user profile and JWT claims |
| `/api/v1/auth/availability/check-field` | Check email availability |
| `/api/v1/auth/groups/permissions` | Fetch user permissions for a group |
| `/.well-known/jwks.json` | Public keys for JWT verification |

> For detailed request/response examples and payload structures, visit the Swagger documentation.