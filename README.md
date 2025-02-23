# Metrics-Monitor

## Overview
Metrics-Monitor is a lightweight Go-based monitoring tool designed to track system metrics such as CPU usage, memory consumption, and other performance indicators. It provides an API to fetch real-time system metrics and logs data for analysis.

## Features
- **Real-time Metrics Collection**: Captures system metrics at configurable intervals.
- **REST API**: Exposes endpoints for fetching system statistics.
- **Gin Framework**: Uses the fast and minimalistic Gin web framework.
- **Graceful Shutdown**: Ensures smooth shutdown of services and database connections.
- **Logging**: Integrated with Uber Zap for structured and high-performance logging.
- **Dockerized**: Easily deployable using Docker and Docker Compose.

## Advantages
- **Lightweight & Fast**: Built in Go, ensuring minimal resource usage.
- **Modular Architecture**: Well-structured for scalability and maintainability.
- **Graceful Error Handling**: Ensures resilience and reliability.
- **Cloud-Ready**: Can be deployed on cloud platforms easily.

## Project Structure
```
├── config                # Configuration files
├── database              # Database connection and initialization
├── docs                  # API documentation (Swagger)
├── handler               # API handlers
├── logger                # Logging setup using Uber Zap
├── models                # Data models
├── router                # Routes and API definitions
├── service               # Core business logic and metric collection
├── utils                 # Utility functions
├── main.go               # Application entry point
├── Dockerfile            # Docker build file
├── docker-compose.yml    # Docker Compose configuration
└── README.md             # Project documentation
```

## Getting Started

### Prerequisites
- **Go 1.22 or later**
- **Docker & Docker Compose**
- **PostgreSQL 15**

### Running Locally (Without Docker)
```sh
git clone https://github.com/ROHITHSAKTHIVEL/Metrics-Monitor.git
cd Metrics-Monitor

# Set environment variables
export DB_HOST=localhost
export DB_USER=postgres
export DB_PASS=postgres
export DB_NAME=metrics_db
export DB_PORT=5432

# Run the application
go run main.go
```

### Running with Docker

#### **Build and Run the Application**
```sh
docker-compose up --build
```

#### **Stopping Containers**
```sh
docker-compose down
```

#### **Check Running Containers**
```sh
docker ps
```

#### **View Logs**
```sh
docker logs -f metrics-monitor
```

#### **Access PostgreSQL in the Container**
```sh
docker exec -it metrics-db psql -U postgres -d metrics_db
```

## API Endpoints
| Method | Endpoint                                             | Description            |
|--------|------------------------------------------------------|------------------------|
| GET    | `/metrics`                                           | Fetch current metrics with pagination Default Pagesizw =10 and default page = 1 |
| GET    | `/metrics?start=<timestamp>&end=<timestamp>`         | Filter metrics by time range. |
| GET    | `/metrics/average?start=<timestamp>&end=<timestamp>` | Return average CPU and memory usage over the specified period.|
## Troubleshooting
### **Common Issues**
1. **Container fails to start**
   - Run `docker logs metrics-monitor` to check logs.
   - Ensure PostgreSQL is running properly.
2. **Database Connection Issues**
   - Verify that `pg_isready` command is available inside the container.
   - Check if `DB_HOST` is set correctly.
3. **Server Not Stopping Gracefully**
   - Ensure the shutdown logic is properly handling signals.
   - Run `docker-compose down` and wait for containers to terminate.


