# Go Microservices Architecture

This project is a microservices-based application built with Go, demonstrating a distributed system architecture with multiple services, databases, and message queues.

## Architecture Overview

```mermaid
graph LR
    style Client fill:#64B5F6,stroke:#1976D2,stroke-width:2px,color:#333,font-size:14px,rx:10,ry:10
    style Frontend fill:#42A5F5,stroke:#1565C0,stroke-width:2px,color:#333,rx:10,ry:10
    style Broker fill:#81C784,stroke:#388E3C,stroke-width:2px,color:#333,rx:10,ry:10
    style Auth fill:#FFB74D,stroke:#F57C00,stroke-width:2px,color:#333,rx:10,ry:10
    style Logger fill:#BA68C8,stroke:#8E24AA,stroke-width:2px,color:#333,rx:10,ry:10
    style Mail fill:#FFCCBC,stroke:#E64A19,stroke-width:2px,color:#333,rx:10,ry:10
    style Listener fill:#FFF59D,stroke:#FBC02D,stroke-width:2px,color:#333,rx:10,ry:10
    style Postgres fill:#A5D6A7,stroke:#2E7D32,stroke-width:2px,color:#333,rx:10,ry:10
    style Mongo fill:#90CAF9,stroke:#1E88E5,stroke-width:2px,color:#333,rx:10,ry:10
    style RabbitMQ fill:#FFE082,stroke:#FFB300,stroke-width:2px,color:#333,rx:10,ry:10
    style MailHog fill:#C5E1A5,stroke:#558B2F,stroke-width:2px,color:#333,rx:10,ry:10

    subgraph "Client Layer"
        Client[Client/Browser]
    end

    subgraph "Application Services"
        Frontend[Frontend Service<br/>Port: 80]
        Broker[Broker Service<br/>Port: 8080<br/>API Gateway]
        Auth[Auth Service<br/>Port: 8081]
        Logger[Logger Service<br/>Port: 80<br/>RPC: 5001]
        Mail[Mail Service<br/>Port: 80]
        Listener[Listener Service<br/>Event Consumer]
    end

    subgraph "Data and Infrastructure"
        Postgres[PostgreSQL<br/>Port: 5433<br/>Users DB]
        Mongo[MongoDB<br/>Port: 27017<br/>Logs DB]
        RabbitMQ[RabbitMQ<br/>Port: 5672<br/>Message Queue]
        MailHog[MailHog<br/>Ports: 1025, 8025<br/>Email Testing]
    end

    Client -->|HTTP| Frontend
    Frontend -->|HTTP :8080| Broker

    Broker -->|HTTP /authenticate| Auth
    Broker -->|RPC :5001| Logger
    Broker -->|HTTP /send| Mail
    Broker -->|AMQP| RabbitMQ

    Auth -->|SQL| Postgres
    Auth -->|HTTP /log| Logger

    Logger -->|MongoDB| Mongo

    Mail -->|SMTP| MailHog

    Listener -->|AMQP Consumer| RabbitMQ
    Listener -->|logs_topic| RabbitMQ

    linkStyle 0 stroke:#1976D2,stroke-width:2px
    linkStyle 1 stroke:#1565C0,stroke-width:2px
    linkStyle 2 stroke:#388E3C,stroke-width:2px
    linkStyle 3 stroke:#8E24AA,stroke-width:2px
    linkStyle 4 stroke:#E64A19,stroke-width:2px
    linkStyle 5 stroke:#FFB300,stroke-width:2px
    linkStyle 6 stroke:#2E7D32,stroke-width:2px
    linkStyle 7 stroke:#1565C0,stroke-width:2px
    linkStyle 8 stroke:#1E88E5,stroke-width:2px
    linkStyle 9 stroke:#558B2F,stroke-width:2px
    linkStyle 10 stroke:#FFB300,stroke-width:2px
    linkStyle 11 stroke:#FFB300,stroke-width:2px
```

## Services

### 1. Frontend Service
- **Port**: 80
- **Description**: Web frontend that provides a UI for testing microservices
- **Technology**: Go with HTML templates
- **Endpoints**: 
  - `GET /` - Test page with service interaction buttons

### 2. Broker Service (API Gateway)
- **Port**: 8080 (mapped from internal port 80)
- **Description**: Acts as the main API gateway/router for all service requests
- **Technology**: Go with Chi router
- **Endpoints**:
  - `POST /` - Health check endpoint
  - `POST /handle` - Routes requests to appropriate services based on action type
- **Responsibilities**:
  - Routes authentication requests to Auth Service
  - Routes logging requests to Logger Service (via RPC)
  - Routes mail requests to Mail Service
  - Can publish events to RabbitMQ

### 3. Auth Service
- **Port**: 8081 (mapped from internal port 80)
- **Description**: Handles user authentication
- **Technology**: Go with PostgreSQL
- **Database**: PostgreSQL (users database)
- **Endpoints**:
  - `POST /authenticate` - Authenticates users
- **Features**:
  - User authentication with email/password
  - Logs authentication events to Logger Service

### 4. Logger Service
- **Port**: 80 (HTTP), 5001 (RPC)
- **Description**: Centralized logging service
- **Technology**: Go with MongoDB, RPC support
- **Database**: MongoDB (logs database)
- **Endpoints**:
  - `POST /log` - HTTP endpoint for logging
  - RPC endpoint on port 5001 for direct logging
- **Features**:
  - Stores logs in MongoDB
  - Supports both HTTP and RPC communication

### 5. Mail Service
- **Port**: 80
- **Description**: Handles email sending
- **Technology**: Go with SMTP
- **Endpoints**:
  - `POST /send` - Sends emails
- **Features**:
  - Sends emails via SMTP
  - Uses MailHog for email testing in development

### 6. Listener Service
- **Description**: Event consumer that listens to RabbitMQ messages
- **Technology**: Go with RabbitMQ AMQP client
- **Features**:
  - Consumes messages from RabbitMQ `logs_topic` exchange
  - Listens for log events (log.INFO, log.WARNING, log.ERROR)

## Infrastructure

### Databases

#### PostgreSQL
- **Port**: 5433 (mapped from internal port 5432)
- **Database**: users
- **Used by**: Auth Service
- **Purpose**: Stores user authentication data

#### MongoDB
- **Port**: 27017
- **Database**: logs
- **Used by**: Logger Service
- **Purpose**: Stores application logs

### Message Queue

#### RabbitMQ
- **Port**: 5672
- **Used by**: Broker Service, Listener Service
- **Purpose**: Asynchronous message processing and event-driven communication
- **Exchange**: `logs_topic` (topic exchange)

### Email Testing

#### MailHog
- **Ports**: 1025 (SMTP), 8025 (Web UI)
- **Used by**: Mail Service
- **Purpose**: Email testing and development tool

## Communication Patterns

1. **HTTP/REST**: Primary communication method between services
   - Frontend → Broker Service
   - Broker Service → Auth Service
   - Broker Service → Mail Service
   - Auth Service → Logger Service

2. **RPC**: Direct procedure calls for high-performance logging
   - Broker Service → Logger Service (port 5001)

3. **AMQP/RabbitMQ**: Asynchronous event-driven communication
   - Broker Service → RabbitMQ (publishes events)
   - Listener Service → RabbitMQ (consumes events)

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for local development)

### Running the Project

1. Navigate to the project directory:
```bash
cd project
```

2. Start all services using Docker Compose:
```bash
docker-compose up --build
```

3. Access the services:
   - Frontend: http://localhost:80 (if exposed)
   - Broker Service: http://localhost:8080
   - Auth Service: http://localhost:8081
   - MailHog UI: http://localhost:8025

### Testing Services

Use the frontend test page to interact with services:
- Test Broker: Tests the broker service health
- Test Auth: Tests authentication flow
- Test Log: Tests logging via RPC
- Test Mail: Tests email sending

## Project Structure

```
go-microservices/
├── auth-service/          # Authentication service
├── broker-service/        # API gateway/broker
├── front-end/             # Web frontend
├── listener-service/      # RabbitMQ event consumer
├── logger-service/        # Logging service
├── mail-service/          # Email service
└── project/              # Docker Compose configuration
    ├── docker-compose.yml
    └── db-data/          # Database volumes
```

## Environment Variables

### Auth Service
- `DSN`: PostgreSQL connection string

### Mail Service
- `MAIL_DOMAIN`: Mail domain
- `MAIL_HOST`: SMTP host (mailhog)
- `MAIL_PORT`: SMTP port (1025)
- `MAIL_FROM_ADDRESS`: Sender email address
- `MAIL_FROM_NAME`: Sender name

## Development

Each service is a standalone Go application that can be run independently or together via Docker Compose. Services communicate through well-defined HTTP endpoints, RPC calls, or message queues.

## License

[Add your license here]

