# Special Status Check Service

## Description

This service implements OS signal handling and graceful shutdown, working with HTTP server, HTTP middleware, OpenAPI (Swagger) specifications, net/http, mux, gin libraries, Prometheus metrics collection, and logging. It also implements the Health Check API pattern.

## Features

The service provides three HTTP endpoints:

1. `GET /success`
   - Always returns 200 OK.

2. `GET /internal`
   - Always returns 500 Internal Server Error.

3. `POST /random`
   - 20% chance of panic, returning 500 internal server error.
   - 30% chance of returning a 500 status error.
   - Otherwise, returns 200 OK.
   - All request bodies and server responses are logged.

4. `GET /metrics`
   - get Prometheus metrics on port 8081.

5. `GET /healthz`
   - Health check endpoint on port 8081.

6. `GET /readyz`
   - Readiness check endpoint on port 8081.

All endpoints are covered with basic Prometheus metrics.

## Installation and Run

### Step 1: Clone the repository

```sh
git clone https://github.com/popeskul/special-status-check.git
cd special-status-check
```

### Step 2: Create the configuration file

Create the `config/config.yaml` file with the following content:

```yaml
server:
  port: 8080
  timeouts:
    write: "15s"
    read: "15s"
    idle: "60s"
  health_check_port: 8081
```

### Step 3: Build and run with Docker Compose

Use Docker Compose to build and run the service.

```sh
make docker-compose-up
```

The service will be available on port 8081(changed to avoid conflicts with other services).

### Step 4: Stop the service

To stop the service, run:

```sh
make docker-compose-down
```

### Endpoints

### `GET /health`

GET /success

Example request:

```sh
curl -i -X GET http://localhost:8080/success
```

### `GET /internal`

Always returns 500 Internal Server Error.

Example request:

```sh
curl -i -X GET http://localhost:8080/internal
```

### `POST /random`

- 20% chance of panic, returning 500 internal server error.
- 30% chance of returning a 500 status error.
- Otherwise, returns 200 OK.

Example request:

```sh
curl -i -X POST http://localhost:8080/random
```

### Prometheus Metrics

The service provides Prometheus metrics on the `/metrics` endpoint.

Example request:

```sh
curl -i -X GET http://localhost:8081/metrics
```

### Health Check

The service provides a health check endpoint on the `/healthz` endpoint.

Example request:

```sh
curl -i -X GET http://localhost:8081/healthz
```

### Readiness Check

The service provides a readiness check endpoint on the `/readyz` endpoint.

Example request:

```sh
curl -i -X GET http://localhost:8081/readyz
```

### Continuous Integration

This project uses GitHub Actions for continuous integration. The CI workflow is defined in `.github/workflows/ci.yml` and includes steps for:

- Checking out the code
- Setting up Go
- Installing dependencies
- Running lint checks
- Building the project
- Running tests
- Building the Docker image
- Running Docker Compose
- Running integration tests
- Tearing down Docker Compose

### License

This project is licensed under the MIT License.
