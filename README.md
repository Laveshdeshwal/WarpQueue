# WarpQueue

WarpQueue is a Go-based job queue prototype with:

- an HTTP API for submitting and inspecting jobs
- an in-memory queue for ordered job processing
- an in-memory job store for lookup by job ID
- worker lifecycle tracking
- automatic retries for failed jobs
- centralized tests

## Current Features

- Submit jobs with `POST /jobs`
- Inspect a single job with `GET /jobs/:id`
- View queue and worker stats with `GET /stats`
- Health check with `GET /health`
- Job lifecycle:
  - `pending`
  - `running`
  - `retrying`
  - `completed`
  - `failed`
- Automatic retry with fixed delay
- Config via environment variables with defaults

## Project Structure

```text
WarpQueue/
├── cmd/
│   ├── server/
│   │   └── main.go
│   └── worker/
│       └── main.go
├── internal/
│   ├── api/
│   │   └── server.go
│   ├── config/
│   │   └── config.go
│   ├── job/
│   │   ├── model.go
│   │   └── store.go
│   ├── logger/
│   │   └── logger.go
│   ├── queue/
│   │   ├── interface.go
│   │   └── memory.go
│   └── worker/
│       ├── handler.go
│       └── pool.go
├── tests/
│   ├── api_test.go
│   └── worker_test.go
├── Makefile
└── go.mod
```

## Architecture

WarpQueue currently uses two in-memory layers:

1. Queue order
   - Maintains FIFO order of job IDs.

2. Job store
   - Maintains a `map[string]*Job` for lookup, updates, stats, and filtering by status.

Flow:

1. Client submits a job through `POST /jobs`
2. Job is saved to the in-memory store
3. Job ID is added to the in-memory queue
4. Worker dequeues the job
5. Job status changes to `running`
6. If handler succeeds, job becomes `completed`
7. If handler fails and retries remain, job becomes `retrying`, then returns to the queue
8. If retries are exhausted, job becomes `failed`

## Configuration

WarpQueue reads config from environment variables and falls back to defaults.

| Variable | Default | Description |
|---|---|---|
| `SERVER_PORT` | `8080` | HTTP server port |
| `WORKER_COUNT` | `3` | Number of worker goroutines |
| `LOG_LEVEL` | `info` | Logger level |
| `QUEUE_TYPE` | `memory` | Queue backend type |
| `SHUTDOWN_TIMEOUT` | `10s` | Graceful shutdown timeout |

Example:

```bash
export SERVER_PORT=8080
export WORKER_COUNT=3
export LOG_LEVEL=debug
export QUEUE_TYPE=memory
export SHUTDOWN_TIMEOUT=10s
```

## API Endpoints

### `POST /jobs`

Create a job.

Request:

```json
{
  "type": "send_email",
  "payload": {
    "to": "user@example.com"
  },
  "priority": 1,
  "max_retries": 2
}
```

Response:

```json
{
  "id": "abc123"
}
```

### `GET /jobs/:id`

Inspect a submitted job.

Response:

```json
{
  "id": "abc123",
  "type": "send_email",
  "status": "completed",
  "attempts": 1
}
```

### `GET /stats`

Operational visibility for job states.

Response:

```json
{
  "total": 10,
  "pending": 3,
  "running": 2,
  "retrying": 1,
  "completed": 3,
  "failed": 1,
  "workers": 3
}
```

### `GET /health`

Health endpoint.

Response:

```json
{
  "Server": "go-warp-queue-server",
  "Status": "running"
}
```

## Running The Project

### Run server

```bash
make run-server
```

### Run worker

```bash
make run-worker
```

### Build binaries

```bash
make build
```

This creates:

- `warpqueue-server`
- `warpqueue-worker`

## Testing

Run all tests:

```bash
make test
```

Run verbose tests:

```bash
make test-verbose
```

The centralized test suite lives in:

- `tests/api_test.go`
- `tests/worker_test.go`

Current tests cover:

- job creation and lookup by ID
- stats endpoint behavior
- retry flow for failed jobs

## Manual API Testing

Start the server:

```bash
make run-server
```

Create a job:

```bash
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{"type":"send_email","payload":{"to":"user@example.com"},"priority":1,"max_retries":2}'
```

Inspect the job:

```bash
curl http://localhost:8080/jobs/<JOB_ID>
```

Check stats:

```bash
curl http://localhost:8080/stats
```

Check health:

```bash
curl http://localhost:8080/health
```



