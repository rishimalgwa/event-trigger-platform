# Event Trigger Platform

This is an event trigger platform built in Go that leverages GORM for ORM, Redis for task management, PostgreSQL for data storage, and Kafka for message queuing. The platform allows you to manage scheduled triggers, including one-time and recurring triggers, execute actions at scheduled times, and maintain event logs for fast retrieval.

## Features

- **Scheduled Triggers**: Support for one-time and recurring triggers that execute tasks based on a defined schedule.
- **Web Interface**: 
  - Fully functional web-based UI for managing triggers
  - Intuitive CRUD (Create, Read, Update, Delete) operations for event triggers
  - Real-time event log monitoring and management
- **Event Log Management**: 
  - Event logs are archived after 2 hours and deleted after 46 hours 
  - All event logs stored in Redis for fast retrieval
  - Separate views for active and archived logs
- **Comprehensive API Endpoints**:
  - Trigger Management:
    - ` /api/trigger`: Create a new trigger
    - ` /api/trigger/`: Retrieve all triggers
    - ` /api/trigger/:id`: Edit an existing trigger
    - ` /api/trigger/:id`: Delete a specific trigger
  - Event Log Management:
    - ` /api/eventlog/active/`: Retrieve active event logs
    - ` /api/eventlog/archived/`: Retrieve archived event logs
  - Testing:
    - ` /trigger/test/`: Create a test trigger
- **Kafka Integration**: Kafka is used for message queuing and task management. The producer is passed to the trigger and event log services.
- **Docker Support**: The project includes Docker configurations to set up services easily.

## Tech Stack

- **Go**: Programming language for backend development.
- **GORM**: ORM for interacting with PostgreSQL.
- **Redis**: Queue management for event logging and trigger task management.
- **PostgreSQL**: Database for storing event data and triggers.
- **Kafka**: Message queuing system used for trigger task management.
- **Docker**: Containerization tool for deploying services.

## Setup Instructions

### Prerequisites

- **Go** (version 1.20+)
- **Docker** (for containerization)
- **Docker Compose** (for setting up services)
- **PostgreSQL** (database setup via Docker)

### Installation Steps

1. **Clone the repository**:
   ```bash
   git clone https://github.com/rishimalgwa/event-trigger-platform.git
   cd event-trigger-platform
   ```

2. **Create a .env file to store environment variables**:
   ```bash 
   cp .env.example .env
   ```

3. **Set up PostgreSQL and Redis using Docker Compose**:
   ```bash
   docker-compose up --build
   ```
   This will start the backend API and services.

4. **Access the API**: 
   - Backend API: http://localhost:80/
   - API Endpoints:
     - Create Trigger: `POST http://localhost:80/api/trigger`
     - Get all triggers: `GET http://localhost:80/api/trigger/`
     - Edit Trigger: `PUT http://localhost:80/api/trigger/:id`
     - Delete Trigger: `DELETE http://localhost:80/api/trigger/:id`
     - Get Active Logs: `GET http://localhost:80/api/eventlog/active/`
     - Get Archived Logs: `GET http://localhost:80/api/eventlog/archived/`
     - Test Trigger: `POST http://localhost:80/trigger/test/`

## Accessing the Web Interface
Once the application is running, simply open `http://localhost:80/index.html` in your web browser to access the web interface.

## Testing the Application

- By default, the platform will hit the `/api/trigger` endpoint to create triggers.
- You can toggle the test button to hit `/api/trigger/test` for test purposes.

## Troubleshooting

- Ensure all environment variables are correctly set in the `.env` file
- Check that all Docker services are running with `docker-compose ps`

## Credits

- ChatGPT: Assisted with HTML UI, Kafka integration, and README writing
